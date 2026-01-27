package auth

import (
	"errors"
	"flight_log_service/db"
	"fmt"
	"log"
	"strings"

	"github.com/thedanisaur/jfl_platform/types"
	"github.com/thedanisaur/jfl_platform/util"

	"github.com/google/cel-go/cel"
	"github.com/google/uuid"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func compileCelToSQL(expr string, table_name string, request_user map[string]interface{}) (string, []interface{}, error) {
	// Parse CEL expression into AST
	env, err := cel.NewEnv(
		cel.Variable("record", cel.MapType(cel.StringType, cel.DynType)),
		cel.Variable("request_user", cel.MapType(cel.StringType, cel.DynType)),
	)
	if err != nil {
		return "", nil, err
	}
	ast, iss := env.Parse(expr)
	if iss.Err() != nil {
		return "", nil, iss.Err()
	}
	checked, iss := env.Check(ast)
	if iss.Err() != nil {
		return "", nil, iss.Err()
	}
	checked_expression, err := cel.AstToCheckedExpr(checked)
	if err != nil {
		return "", nil, err
	}

	// Convert AST to SQL
	return compileExpression(checked_expression.GetExpr(), table_name, request_user)
}

func compileExpression(expression *exprpb.Expr, table_name string, request_user map[string]interface{}) (string, []interface{}, error) {
	switch expression_kind := expression.ExprKind.(type) {

	// boolean constants
	case *exprpb.Expr_ConstExpr:
		switch v := expression_kind.ConstExpr.ConstantKind.(type) {

		case *exprpb.Constant_BoolValue:
			if v.BoolValue {
				return "1=1", nil, nil
			} else {
				return "1=0", nil, nil
			}

		case *exprpb.Constant_Int64Value:
			return "?", []interface{}{v.Int64Value}, nil

		case *exprpb.Constant_DoubleValue:
			return "?", []interface{}{v.DoubleValue}, nil

		case *exprpb.Constant_StringValue:
			return "?", []interface{}{v.StringValue}, nil

		case *exprpb.Constant_NullValue:
			return "NULL", nil, nil
		}

		return "", nil, errors.New("unsupported constant type")

	// identifiers: true / false
	case *exprpb.Expr_IdentExpr:
		switch expression_kind.IdentExpr.Name {
		case "true":
			return "1=1", nil, nil
		case "false":
			return "1=0", nil, nil
		default:
			return "", nil, fmt.Errorf("unsupported identifier: %s", expression_kind.IdentExpr.Name)
		}

	// record.field OR request_user.field
	case *exprpb.Expr_SelectExpr:
		// operand must be an identifier
		ident, ok := expression_kind.SelectExpr.Operand.ExprKind.(*exprpb.Expr_IdentExpr)
		if !ok {
			return "", nil, errors.New("unsupported select operand")
		}

		switch ident.IdentExpr.Name {
		case "record":
			return fmt.Sprintf("%s.%s", table_name, expression_kind.SelectExpr.Field), nil, nil

		case "request_user":
			val, ok := request_user[expression_kind.SelectExpr.Field]
			if !ok {
				return "", nil, fmt.Errorf("request_user.%s not provided", expression_kind.SelectExpr.Field)
			}

			// If it's a slice, return as-is for `in`
			if s, ok := val.([]interface{}); ok {
				return "?", s, nil
			}

			return "?", []interface{}{val}, nil

		default:
			return "", nil, fmt.Errorf("unsupported select base: %s", ident.IdentExpr.Name)
		}

	// binary operators
	case *exprpb.Expr_CallExpr:
		if len(expression_kind.CallExpr.Args) != 2 {
			return "", nil, errors.New("only binary operators are supported")
		}

		left_sql, left_args, err := compileExpression(expression_kind.CallExpr.Args[0], table_name, request_user)
		if err != nil {
			return "", nil, err
		}

		right_sql, right_args, err := compileExpression(expression_kind.CallExpr.Args[1], table_name, request_user)
		if err != nil {
			return "", nil, err
		}

		sql_operation := map[string]string{
			"_==_":     "=",
			"_!=_":     "!=",
			"_<_":      "<",
			"_<=_":     "<=",
			"_>_":      ">",
			"_>=_":     ">=",
			"_&&_":     "AND",
			"_||_":     "OR",
			"_in_":     "IN",
			"_not_in_": "NOT IN",
		}[expression_kind.CallExpr.Function]

		if sql_operation == "" {
			return "", nil, fmt.Errorf("unsupported operator: %s", expression_kind.CallExpr.Function)
		}

		// Support `in`/`not in`. Only this order will work: 'single in list'
		// i.e.: record.unit_id in request_user.unit_ids
		if sql_operation == "IN" || sql_operation == "NOT IN" {
			// right_args must be a slice
			if len(right_args) == 0 {
				return "", nil, errors.New("right side of `in`/`not in` must be a list")
			}

			// build (?, ?, ?)
			placeholders := make([]string, len(right_args))
			for i := range placeholders {
				placeholders[i] = "?"
			}

			sql := fmt.Sprintf("(%s %s (%s))", left_sql, sql_operation, strings.Join(placeholders, ", "))
			return sql, right_args, nil
		}

		return fmt.Sprintf("(%s %s %s)", left_sql, sql_operation, right_sql), append(left_args, right_args...), nil
	}

	return "", nil, errors.New("unsupported expression type")
}

func EvaluateRead(txid uuid.UUID, role_name string, resource string, operation string, table_name string, request_user map[string]interface{}) (string, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(EvaluateRead))

	policies, err := loadPermissions(txid, role_name, resource, operation)
	if err != nil {
		return "", err
	}
	if len(policies) <= 0 {
		return "", errors.New("not authorized")
	}

	for _, policy := range policies {
		// Implicit deny overrides any allow
		if policy.Effect != "allow" {
			return "", errors.New("not authorized")
		}
	}

	// If we made it here we are authorized to read, let's build the mysql filters
	var filters []string
	var args []interface{}

	for _, policy := range policies {
		// compile CEL -> SQL filter (we'll implement this next)
		sql, p, err := compileCelToSQL(policy.ConditionExpression, table_name, request_user)
		if err != nil {
			return "", err
		}
		filters = append(filters, sql)
		args = append(args, p...)
	}

	// Combine filters using OR
	filter_string := strings.Join(filters, " OR ")
	// This shouldn't happen
	if filter_string == "" {
		// Implicit deny overrides any allow
		filter_string = "1=0"
	}
	return filter_string, nil
}

// func Evaluate(txid uuid.UUID, role_name string, resource string, operation string, context map[string]interface{}) (bool, error) {
// 	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(Evaluate))

// 	policies, err := loadPermissions(txid, role_name, resource, operation)
// 	if err != nil {
// 		return false, err
// 	}

// 	// Default deny policy
// 	allowed := false
// 	for _, policy := range policies {
// 		ok, err := evaluateCEL(policy.ConditionExpression, context)
// 		if err != nil {
// 			return false, err
// 		}

// 		if ok && policy.Effect == "allow" {
// 			allowed = true
// 		}
// 		if ok && policy.Effect == "deny" {
// 			return false, nil // deny overrides allow
// 		}
// 	}
// 	return allowed, nil
// }

// func evaluateCEL(condition_expression string, context map[string]interface{}) (bool, error) {
// 	env, _ := cel.NewEnv(
// 		cel.Declarations(
// 			decls.NewVar("user", decls.NewMapType(decls.String, decls.Dyn)),
// 			decls.NewVar("resource", decls.NewMapType(decls.String, decls.Dyn)),
// 			decls.NewVar("extra", decls.NewMapType(decls.String, decls.Dyn)),
// 		),
// 	)

// 	ast, _ := env.Parse(condition_expression)
// 	checked, _ := env.Check(ast)
// 	prg, _ := env.Program(checked)

// 	out, _, err := prg.Eval(context)
// 	if err != nil {
// 		return false, err
// 	}

// 	result, ok := out.Value().(bool)
// 	if !ok {
// 		return false, errors.New("policy must return bool")
// 	}
// 	return result, nil
// }

func loadPermissions(txid uuid.UUID, role_name string, resource string, operation string) ([]types.PermissionDTO, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(loadPermissions))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())
	database, err := db.GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return nil, errors.New(err_string)
	}

	permissions_query := `
		SELECT
			p.resource
			, p.operation
			, p.effect
			, p.cond_type
			, p.cond_expr
		FROM permissions p
		WHERE p.role_name = ?
		  AND p.resource = ?
		  AND p.operation = ?
	`
	rows, err := database.Query(permissions_query, role_name, resource, operation)
	if err != nil {
		log.Printf("Failed to retrieve policies for role: %s\n%s\n", role_name, err.Error())
		return nil, errors.New("failed to retrieve policies")
	}
	defer rows.Close()

	permissions := make([]types.PermissionDTO, 0)
	for rows.Next() {
		var permission types.PermissionDTO
		err := rows.Scan(
			&permission.Resource,
			&permission.Operation,
			&permission.Effect,
			&permission.ConditionType,
			&permission.ConditionExpression,
		)
		if err != nil {
			log.Printf("Failed to parse a permission for role: %s \n%s\n", role_name, err.Error())
			return nil, fmt.Errorf("failed to parse a permission for role: %s", role_name)
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}
