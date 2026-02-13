package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/thedanisaur/jfl_platform/types"
	"github.com/thedanisaur/jfl_platform/util"
)

func LoadPermissions(txid uuid.UUID, role_name string, resource string, operation string) ([]types.PermissionDTO, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(LoadPermissions))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())
	database, err := GetInstance()
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
