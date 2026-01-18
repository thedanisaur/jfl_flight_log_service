package db

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/thedanisaur/jfl_platform/types"
	"github.com/thedanisaur/jfl_platform/types/UserStatus"
	"github.com/thedanisaur/jfl_platform/util"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetPassword(username string) (string, error) {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetPassword), txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return "", errors.New("failed to connect to DB")
	}
	query := `
		SELECT password_hash
		FROM users
		WHERE email = LOWER(?)
	`
	row := database.QueryRow(query, username)
	var password_hash string
	err = row.Scan(&password_hash)
	if err != nil {
		log.Printf("Invalid username: %s\n", err.Error())
		return "", errors.New("failed to connect to DB")
	}
	return password_hash, nil
}

func GetUserById(user_id uuid.UUID) (types.UserResponse, error) {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetUserById), txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return types.UserResponse{}, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) id
			, email
			, password_hash
			, first_name
			, last_name
			, call_sign
			, primary_mds
			, secondary_mds
			, ssn_last_4
			, flight_auth_code
			, issuing_unit
			, unit_charged
			, harm_location
			, status
			, is_instructor
			, is_evaluator
			, role_id
			, role_requested_id
			, created_on
			, updated_on
			, last_logged_in
		FROM users
		WHERE id = UUID_TO_BIN(?)
	`
	row := database.QueryRow(query, user_id)
	var user_dbo types.UserDbo
	err = row.Scan(
		&user_dbo.ID,
		&user_dbo.Email,
		&user_dbo.PasswordHash,
		&user_dbo.FirstName,
		&user_dbo.LastName,
		&user_dbo.CallSign,
		&user_dbo.PrimaryMDS,
		&user_dbo.SecondaryMDS,
		&user_dbo.SSNLast4,
		&user_dbo.FlightAuthCode,
		&user_dbo.IssuingUnit,
		&user_dbo.UnitCharged,
		&user_dbo.HarmLocation,
		&user_dbo.Status,
		&user_dbo.IsInstructor,
		&user_dbo.IsEvaluator,
		&user_dbo.RoleId,
		&user_dbo.RoleRequestedId,
		&user_dbo.CreatedOn,
		&user_dbo.UpdatedOn,
		&user_dbo.LastLoggedIn,
	)
	if err != nil {
		log.Printf("Failed to retrieve user:\n%s\n", err.Error())
		return types.UserResponse{}, errors.New("failed to retrieve user")
	}
	user, err := userDboToResponse(user_dbo)
	if err != nil {
		return types.UserResponse{}, err
	}
	return user, nil
}

func GetUserByEmail(email string) (types.UserResponse, error) {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetUserByEmail), txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return types.UserResponse{}, errors.New("failed to connect to database")
	}
	query := `
		SELECT BIN_TO_UUID(id) id
			, email
			, password_hash
			, first_name
			, last_name
			, call_sign
			, primary_mds
			, secondary_mds
			, ssn_last_4
			, flight_auth_code
			, issuing_unit
			, unit_charged
			, harm_location
			, status
			, is_instructor
			, is_evaluator
			, role_id
			, role_requested_id
			, created_on
			, updated_on
			, last_logged_in
		FROM users
		WHERE email = ?
	`
	row := database.QueryRow(query, email)
	var user_dbo types.UserDbo
	err = row.Scan(
		&user_dbo.ID,
		&user_dbo.Email,
		&user_dbo.PasswordHash,
		&user_dbo.FirstName,
		&user_dbo.LastName,
		&user_dbo.CallSign,
		&user_dbo.PrimaryMDS,
		&user_dbo.SecondaryMDS,
		&user_dbo.SSNLast4,
		&user_dbo.FlightAuthCode,
		&user_dbo.IssuingUnit,
		&user_dbo.UnitCharged,
		&user_dbo.HarmLocation,
		&user_dbo.Status,
		&user_dbo.IsInstructor,
		&user_dbo.IsEvaluator,
		&user_dbo.RoleId,
		&user_dbo.RoleRequestedId,
		&user_dbo.CreatedOn,
		&user_dbo.UpdatedOn,
		&user_dbo.LastLoggedIn,
	)
	if err != nil {
		log.Printf("Failed to retrieve user:\n%s\n", err.Error())
		return types.UserResponse{}, errors.New("failed to retrieve user")
	}
	user, err := userDboToResponse(user_dbo)
	if err != nil {
		return types.UserResponse{}, err
	}
	return user, nil
}

func getUserDboById(user_id uuid.UUID) (types.UserDbo, error) {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(getUserDboById), txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return types.UserDbo{}, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) id
			, email
			, password_hash
			, first_name
			, last_name
			, call_sign
			, primary_mds
			, secondary_mds
			, ssn_last_4
			, flight_auth_code
			, issuing_unit
			, unit_charged
			, harm_location
			, status
			, is_instructor
			, is_evaluator
			, role_id
			, role_requested_id
			, created_on
			, updated_on
			, last_logged_in
		FROM users
		WHERE id = UUID_TO_BIN(?)
	`
	row := database.QueryRow(query, user_id)
	var user types.UserDbo
	err = row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.CallSign,
		&user.PrimaryMDS,
		&user.SecondaryMDS,
		&user.SSNLast4,
		&user.FlightAuthCode,
		&user.IssuingUnit,
		&user.UnitCharged,
		&user.HarmLocation,
		&user.Status,
		&user.IsInstructor,
		&user.IsEvaluator,
		&user.RoleId,
		&user.RoleRequestedId,
		&user.CreatedOn,
		&user.UpdatedOn,
		&user.LastLoggedIn,
	)
	if err != nil {
		log.Printf("Failed to retrieve user:\n%s\n", err.Error())
		return types.UserDbo{}, errors.New("failed to retrieve user")
	}
	return user, nil
}

func GetUsers() ([]types.UserResponse, error) {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetUsers), txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return nil, errors.New("failed to connect to database")
	}
	query := fmt.Sprintf(`
		SELECT BIN_TO_UUID(id) id
			, email
			, password_hash
			, first_name
			, last_name
			, call_sign
			, primary_mds
			, secondary_mds
			, ssn_last_4
			, flight_auth_code
			, issuing_unit
			, unit_charged
			, harm_location
			, status
			, is_instructor
			, is_evaluator
			, role_id
			, role_requested_id
			, created_on
			, updated_on
			, last_logged_in
		FROM users
		WHERE status != '%s'
	`, UserStatus.Deactivated)
	rows, err := database.Query(query)
	if err != nil {
		log.Printf("failed to query database:\n%s\n", err.Error())
		return nil, errors.New("failed to connect to database")
	}

	roles, err := GetRolesAsMap()
	if err != nil {
		return nil, err
	}
	var users []types.UserResponse
	for rows.Next() {
		var user_dbo types.UserDbo
		err = rows.Scan(
			&user_dbo.ID,
			&user_dbo.Email,
			&user_dbo.PasswordHash,
			&user_dbo.FirstName,
			&user_dbo.LastName,
			&user_dbo.CallSign,
			&user_dbo.PrimaryMDS,
			&user_dbo.SecondaryMDS,
			&user_dbo.SSNLast4,
			&user_dbo.FlightAuthCode,
			&user_dbo.IssuingUnit,
			&user_dbo.UnitCharged,
			&user_dbo.HarmLocation,
			&user_dbo.Status,
			&user_dbo.IsInstructor,
			&user_dbo.IsEvaluator,
			&user_dbo.RoleId,
			&user_dbo.RoleRequestedId,
			&user_dbo.CreatedOn,
			&user_dbo.UpdatedOn,
			&user_dbo.LastLoggedIn,
		)
		if err != nil {
			log.Printf("failed to scan row\n%s\n", err.Error())
			continue
		}
		role := roles[user_dbo.RoleId.String()]
		requested_role := roles[user_dbo.RoleRequestedId.String()]
		var user types.UserResponse = types.UserResponse{
			ID:             user_dbo.ID,
			Email:          user_dbo.Email,
			PasswordHash:   user_dbo.PasswordHash,
			FirstName:      user_dbo.FirstName,
			LastName:       user_dbo.LastName,
			CallSign:       user_dbo.CallSign,
			PrimaryMDS:     user_dbo.PrimaryMDS,
			SecondaryMDS:   user_dbo.SecondaryMDS,
			SSNLast4:       user_dbo.SSNLast4,
			FlightAuthCode: user_dbo.FlightAuthCode,
			IssuingUnit:    user_dbo.IssuingUnit,
			UnitCharged:    user_dbo.UnitCharged,
			HarmLocation:   user_dbo.HarmLocation,
			Status:         user_dbo.Status,
			IsInstructor:   user_dbo.IsInstructor,
			IsEvaluator:    user_dbo.IsEvaluator,
			Role:           role.Name,
			RoleRequested:  &requested_role.Name,
			CreatedOn:      user_dbo.CreatedOn,
			UpdatedOn:      user_dbo.UpdatedOn,
			LastLoggedIn:   user_dbo.LastLoggedIn,
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		log.Println("error scanning rows")
		return nil, errors.New("failed to connect to database")
	}
	return users, nil
}

func InsertUser(hashed_password string, user types.UserRequest) (int64, error) {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(InsertUser), txid.String())
	err_string := fmt.Sprintf("database error: %s\n", txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return -1, errors.New(err_string)
	}
	role, err := GetRoleByName(*user.Role.Value)
	if err != nil {
		return -1, err
	}
	requested_role, err := GetRoleByName(*user.RoleRequested.Value)
	if err != nil {
		return -1, err
	}
	query := `
		INSERT INTO users
		(
			email
			, password_hash
			, first_name
			, last_name
			, call_sign
			, primary_mds
			, secondary_mds
			, ssn_last_4
			, flight_auth_code
			, issuing_unit
			, unit_charged
			, harm_location
			, status
			, is_instructor
			, is_evaluator
			, role_id
			, role_requested_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, UUID_TO_BIN(?), UUID_TO_BIN(?))
	`
	result, err := database.Exec(
		query,
		user.Email,
		hashed_password,
		user.FirstName,
		user.LastName,
		user.CallSign,
		user.PrimaryMDS,
		user.SecondaryMDS,
		user.SSNLast4,
		user.FlightAuthCode,
		user.IssuingUnit,
		user.UnitCharged,
		user.HarmLocation,
		UserStatus.Pending,
		user.IsInstructor,
		user.IsEvaluator,
		role.Id,
		requested_role.Id,
	)
	if err != nil {
		log.Printf("failed user insert\n%s\n", err.Error())
		return -1, errors.New(err_string)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("failed retrieve inserted id\n%s\n", err.Error())
		return -1, errors.New(err_string)
	}
	return id, nil
}

func UpdateUser(request_user_id uuid.UUID, target_user types.UserRequest) (int64, error) {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(UpdateUser), txid.String())
	request_user, err := getUserDboById(request_user_id)
	if err != nil {
		log.Printf("Failed to get User\n%s\n", err.Error())
		return -1, errors.New("Could not fetch current user")
	}

	var set_clauses []string
	var arguments []interface{}

	set_clauses, arguments = AddNullableString("email", target_user.Email, set_clauses, arguments)
	if target_user.UpdatePassword.Set {
		// Compare the password sent to the password we expect.
		err := bcrypt.CompareHashAndPassword([]byte(request_user.PasswordHash), []byte(*target_user.Password.Value))
		if err != nil {
			return -1, errors.New("invalid password")
		}
		// Now hash the new password
		hashed_password, err := bcrypt.GenerateFromPassword([]byte(*target_user.UpdatePassword.Value), 12)
		if err != nil {
			return -1, errors.New("failed to hash password")
		}
		update_hashed_password := string(hashed_password)
		// Make sure to override the old password
		target_user.Password = types.NullableString{Set: true, Value: &update_hashed_password}
		set_clauses, arguments = AddNullableString("password_hash", target_user.Password, set_clauses, arguments)
	}
	set_clauses, arguments = AddNullableString("first_name", target_user.FirstName, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("last_name", target_user.LastName, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("call_sign", target_user.CallSign, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("primary_mds", target_user.PrimaryMDS, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("secondary_mds", target_user.SecondaryMDS, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("ssn_last_4", target_user.SSNLast4, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("flight_auth_code", target_user.FlightAuthCode, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("issuing_unit", target_user.IssuingUnit, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("unit_charged", target_user.UnitCharged, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("harm_location", target_user.HarmLocation, set_clauses, arguments)
	set_clauses, arguments = AddNullableString("status", target_user.Status, set_clauses, arguments)
	set_clauses, arguments = AddNullableBool("is_instructor", target_user.IsInstructor, set_clauses, arguments)
	set_clauses, arguments = AddNullableBool("is_evaluator", target_user.IsEvaluator, set_clauses, arguments)
	if target_user.Role.Set {
		role, err := GetRoleByName(*target_user.Role.Value)
		if err != nil {
			return -1, err
		}
		role_id_string := role.Id.String()
		role_id := types.NullableString{Set: true, Value: &role_id_string}
		set_clauses, arguments = AddNullableUUID("role_id", role_id, set_clauses, arguments)
	}
	if target_user.RoleRequested.Set {
		role, err := GetRoleByName(*target_user.RoleRequested.Value)
		if err != nil {
			return -1, err
		}
		role_id_string := role.Id.String()
		role_id := types.NullableString{Set: true, Value: &role_id_string}
		set_clauses, arguments = AddNullableUUID("role_requested_id", role_id, set_clauses, arguments)
	}
	set_clauses, arguments = AddNullableTime("created_on", target_user.CreatedOn, set_clauses, arguments)
	set_clauses, arguments = AddNullableTime("updated_on", target_user.UpdatedOn, set_clauses, arguments)
	set_clauses, arguments = AddNullableTime("last_logged_in", target_user.LastLoggedIn, set_clauses, arguments)

	if len(set_clauses) == 0 {
		return -1, fmt.Errorf("no fields to update")
	}

	query := bytes.Buffer{}
	query.WriteString("UPDATE users SET ")
	query.WriteString(set_clauses[0])
	for _, clause := range set_clauses[1:] {
		query.WriteString(", ")
		query.WriteString(clause)
	}
	query.WriteString(" WHERE id = UUID_TO_BIN(?)")

	arguments = append(arguments, target_user.ID)
	result, err := database.Exec(query.String(), arguments...)
	if err != nil {
		log.Printf("failed user update\n%s\n", err.Error())
		return -1, errors.New("database error")
	}
	count, err := result.RowsAffected()
	if err != nil {
		log.Printf("failed retrieve count\n%s\n", err.Error())
		return -1, errors.New("database error")
	}
	return count, nil
}

func userDboToResponse(user_dbo types.UserDbo) (types.UserResponse, error) {
	// Grab the roles and then create the return object
	role, err := GetRoleById(user_dbo.RoleId)
	if err != nil {
		return types.UserResponse{}, err
	}
	var requested_role types.RoleDbo
	if user_dbo.RoleRequestedId != nil {
		requested_role, err = GetRoleById(*user_dbo.RoleRequestedId)
		if err != nil {
			return types.UserResponse{}, err
		}
	}
	var user types.UserResponse = types.UserResponse{
		ID:             user_dbo.ID,
		Email:          user_dbo.Email,
		PasswordHash:   user_dbo.PasswordHash,
		FirstName:      user_dbo.FirstName,
		LastName:       user_dbo.LastName,
		CallSign:       user_dbo.CallSign,
		PrimaryMDS:     user_dbo.PrimaryMDS,
		SecondaryMDS:   user_dbo.SecondaryMDS,
		SSNLast4:       user_dbo.SSNLast4,
		FlightAuthCode: user_dbo.FlightAuthCode,
		IssuingUnit:    user_dbo.IssuingUnit,
		UnitCharged:    user_dbo.UnitCharged,
		HarmLocation:   user_dbo.HarmLocation,
		Status:         user_dbo.Status,
		IsInstructor:   user_dbo.IsInstructor,
		IsEvaluator:    user_dbo.IsEvaluator,
		Role:           role.Name,
		RoleRequested:  &requested_role.Name,
		CreatedOn:      user_dbo.CreatedOn,
		UpdatedOn:      user_dbo.UpdatedOn,
		LastLoggedIn:   user_dbo.LastLoggedIn,
	}
	return user, nil
}
