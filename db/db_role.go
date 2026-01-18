package db

import (
	"errors"
	"log"

	"github.com/thedanisaur/jfl_platform/types"
	"github.com/thedanisaur/jfl_platform/util"

	"github.com/google/uuid"
)

func GetRoleById(id uuid.UUID) (types.RoleDbo, error) {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetRoleById), txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return types.RoleDbo{}, errors.New("failed to connect to database")
	}
	query := `
		SELECT id
			, name
			, displayname
		FROM roles
		WHERE id = UUID_TO_BIN(?)
	`
	row := database.QueryRow(query, id)
	var role types.RoleDbo
	err = row.Scan(&role.Id, &role.Name, &role.DisplayName)
	if err != nil {
		log.Printf("invalid role id: %s\n", err.Error())
		return types.RoleDbo{}, errors.New("failed to connect to database")
	}
	return role, nil
}

func GetRoleByName(name string) (types.RoleDbo, error) {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetRoleByName), txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return types.RoleDbo{}, errors.New("failed to connect to database")
	}
	query := `
		SELECT id
			, name
			, displayname
		FROM roles
		WHERE name = LOWER(?)
	`
	row := database.QueryRow(query, name)
	var role types.RoleDbo
	err = row.Scan(&role.Id, &role.Name, &role.DisplayName)
	if err != nil {
		log.Printf("invalid role name: %s\n", err.Error())
		return types.RoleDbo{}, errors.New("failed to connect to database")
	}
	return role, nil
}

func GetRolesAsMap() (map[string]types.RoleDbo, error) {
	txid := uuid.New()
	log.Printf("%s | %s\n", util.GetFunctionName(GetRolesAsMap), txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return map[string]types.RoleDbo{}, errors.New("failed to connect to database")
	}
	query := `
		SELECT BIN_TO_UUID(id)
			, name
			, displayname
		FROM roles
	`
	rows, err := database.Query(query)
	if err != nil {
		log.Printf("failed to query database:\n%s\n", err.Error())
		return map[string]types.RoleDbo{}, errors.New("failed to connect to database")
	}

	roles := make(map[string]types.RoleDbo)
	for rows.Next() {
		var role types.RoleDbo
		err := rows.Scan(
			&role.Id,
			&role.Name,
			&role.DisplayName,
		)
		if err != nil {
			log.Printf("failed to scan row\n%s\n", err.Error())
			continue
		}
		roles[role.Id.String()] = role
	}

	err = rows.Err()
	if err != nil {
		log.Println("error scanning rows")
		return map[string]types.RoleDbo{}, errors.New("failed to connect to database")
	}

	return roles, nil
}
