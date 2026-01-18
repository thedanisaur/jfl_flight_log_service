package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/thedanisaur/jfl_platform/types"

	_ "github.com/go-sql-driver/mysql"
)

var lock = &sync.Mutex{}

var database *sql.DB

func AddNullableBool(col string, field types.NullableBool, set_clauses []string, arguments []interface{}) ([]string, []interface{}) {
	if field.Set {
		if field.Value == nil {
			set_clauses = append(set_clauses, fmt.Sprintf("%s = NULL", col))
		} else {
			set_clauses = append(set_clauses, fmt.Sprintf("%s = ?", col))
			arguments = append(arguments, *field.Value)
		}
	}
	return set_clauses, arguments
}

func AddNullableBytes(col string, field types.NullableBytes, set_clauses []string, arguments []interface{}) ([]string, []interface{}) {
	if field.Set {
		if field.Value == nil {
			set_clauses = append(set_clauses, fmt.Sprintf("%s = NULL", col))
		} else {
			set_clauses = append(set_clauses, fmt.Sprintf("%s = ?", col))
			arguments = append(arguments, *field.Value)
		}
	}
	return set_clauses, arguments
}

func AddNullableString(col string, field types.NullableString, set_clauses []string, arguments []interface{}) ([]string, []interface{}) {
	if field.Set {
		if field.Value == nil {
			set_clauses = append(set_clauses, fmt.Sprintf("%s = NULL", col))
		} else {
			set_clauses = append(set_clauses, fmt.Sprintf("%s = ?", col))
			arguments = append(arguments, *field.Value)
		}
	}
	return set_clauses, arguments
}

func AddNullableTime(col string, field types.NullableTime, set_clauses []string, arguments []interface{}) ([]string, []interface{}) {
	if field.Set {
		if field.Value == nil {
			set_clauses = append(set_clauses, fmt.Sprintf("%s = NULL", col))
		} else {
			set_clauses = append(set_clauses, fmt.Sprintf("%s = ?", col))
			arguments = append(arguments, *field.Value)
		}
	}
	return set_clauses, arguments
}

func AddNullableUUID(col string, field types.NullableString, set_clauses []string, arguments []interface{}) ([]string, []interface{}) {
	if field.Set {
		if field.Value == nil {
			set_clauses = append(set_clauses, fmt.Sprintf("%s = NULL", col))
		} else {
			set_clauses = append(set_clauses, fmt.Sprintf("%s = UUID_TO_BIN(?)", col))
			arguments = append(arguments, *field.Value)
		}
	}
	return set_clauses, arguments
}

// TODO move all db logic to db package
func GetInstance() (*sql.DB, error) {
	if database == nil {
		lock.Lock()
		defer lock.Unlock()
		if database == nil {
			env, err := readDatabaseEnv()
			if err != nil {
				return nil, fmt.Errorf("connection error: %s", err.Error())
			}
			conn_str := fmt.Sprintf("%s:%s@/%s?parseTime=true", env.Username, env.Password, env.Name)
			db, err := sql.Open(env.Driver, conn_str)
			if err != nil {
				return nil, fmt.Errorf("connection error: %s", err.Error())
			}
			// defer db.Close()
			// See "Important settings" section.
			db.SetConnMaxLifetime(time.Minute * 3)
			db.SetMaxOpenConns(10)
			db.SetMaxIdleConns(10)

			// Open doesn't open a connection. Validate DSN data:
			err = db.Ping()
			if err != nil {
				return nil, fmt.Errorf("connection error: %s", err.Error())
			} else {
				log.Printf("Connected to: %s", env.Name)
			}
			database = db
		}
	}
	return database, nil
}

func readDatabaseEnv() (*types.DbConfig, error) {
	username, username_set := os.LookupEnv("MSDBUSERNAME")
	password, password_set := os.LookupEnv("MSDBPASSWORD")
	db_name, db_name_set := os.LookupEnv("MSDBNAME")
	db_driver, db_driver_set := os.LookupEnv("MSDBDRIVER")
	var db_conn types.DbConfig
	if !username_set || !password_set || !db_name_set || !db_driver_set {
		json_file, err := os.Open("./secrets/db.env")
		if err != nil {
			log.Printf("Reading env file error: %s", err.Error())
			return nil, errors.New("could not read db env file aborting")
		}
		defer json_file.Close()
		bytes, _ := io.ReadAll(json_file)
		json.Unmarshal(bytes, &db_conn)
	} else {
		db_conn = types.DbConfig{
			Username: username,
			Password: password,
			Name:     db_name,
			Driver:   db_driver,
		}
	}
	return &db_conn, nil
}
