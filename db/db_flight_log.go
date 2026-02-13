package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/thedanisaur/jfl_platform/types"
	"github.com/thedanisaur/jfl_platform/util"

	"github.com/google/uuid"
)

func DeleteFlightlog(txid uuid.UUID, user_id uuid.UUID, flight_log_id uuid.UUID) (uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(DeleteFlightlog))
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return uuid.Nil, errors.New("failed to connect to DB")
	}
	transaction, err := database.BeginTx(context.Background(), nil)
	if err != nil {
		log.Printf("Failed to initiate transaction\n%s\n", err.Error())
		return uuid.Nil, errors.New("failed to connect to DB")
	}
	defer transaction.Rollback()

	// Delete flight log's comment records
	comments_query := `DELETE FROM flight_log_comments WHERE flight_log_id = UUID_TO_BIN(?)`
	comments_result, err := database.Exec(comments_query, flight_log_id)
	if err != nil {
		log.Printf("Failed to delete flight log comments: %s for user: %s\n%s\n", flight_log_id, user_id, err.Error())
		return uuid.Nil, errors.New("failed to delete flight log comments")
	}
	_, err = comments_result.RowsAffected()
	if err != nil {
		return uuid.Nil, errors.New("failed to delete flight log comments")
	}

	// Delete flight log's aircrew records
	aircrews_query := `DELETE FROM aircrews WHERE flight_log_id = UUID_TO_BIN(?)`
	aircrews_result, err := database.Exec(aircrews_query, flight_log_id)
	if err != nil {
		log.Printf("Failed to delete aircrews: %s for user: %s\n%s\n", flight_log_id, user_id, err.Error())
		return uuid.Nil, errors.New("failed to delete aircrews")
	}
	_, err = aircrews_result.RowsAffected()
	if err != nil {
		return uuid.Nil, errors.New("failed to delete aircrews")
	}

	// Delete flight log's mission records
	missions_query := `DELETE FROM missions WHERE flight_log_id = UUID_TO_BIN(?)`
	missions_result, err := database.Exec(missions_query, flight_log_id)
	if err != nil {
		log.Printf("Failed to delete missions: %s for user: %s\n%s\n", flight_log_id, user_id, err.Error())
		return uuid.Nil, errors.New("failed to delete missions")
	}
	_, err = missions_result.RowsAffected()
	if err != nil {
		return uuid.Nil, errors.New("failed to delete missions")
	}

	// Delete flight log
	flight_log_query := `DELETE FROM flight_logs WHERE id = UUID_TO_BIN(?)`
	flight_log_result, err := database.Exec(flight_log_query, flight_log_id)
	if err != nil {
		log.Printf("Failed to delete flight log: %s for user: %s\n%s\n", flight_log_id, user_id, err.Error())
		return uuid.Nil, errors.New("failed to delete flight log")
	}
	_, err = flight_log_result.RowsAffected()
	if err != nil {
		return uuid.Nil, errors.New("failed to delete flight log")
	}

	return flight_log_id, nil
}

func GetAirCrews(txid uuid.UUID, flight_log_id uuid.UUID) ([]types.FlightLogAircrewDTO, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetAirCrews))
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return nil, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) AS id
			, BIN_TO_UUID(flight_log_id) AS flight_log_id
			, user_id
			, flying_origin
			, flight_auth_code
			, time_primary
			, time_secondary
			, time_instructor
			, time_evaluator
			, time_other
			, total_aircrew_duration_decimal
			, total_aircrew_sorties
			, cond_night_time
			, cond_instrument_time
			, cond_sim_instrument_time
			, cond_nvg_time
			, cond_combat_time
			, cond_combat_sortie
			, cond_combat_support_time
			, cond_combat_support_sortie
			, aircrew_role_type
		FROM aircrews
		WHERE flight_log_id = UUID_TO_BIN(?)
	`
	rows, err := database.Query(query, flight_log_id)
	if err != nil {
		log.Printf("Failed to retrieve aircrew members for flight log: %s \n%s\n", flight_log_id, err.Error())
		return nil, fmt.Errorf("failed to retrieve aircrew members for flight log: %s", flight_log_id)
	}
	defer rows.Close()

	aircrews := make([]types.FlightLogAircrewDTO, 0)
	for rows.Next() {
		var aircrew types.FlightLogAircrewDTO
		err := rows.Scan(
			&aircrew.ID,
			&aircrew.FlightLogID,
			&aircrew.UserID,
			&aircrew.FlyingOrigin,
			&aircrew.FlightAuthCode,
			&aircrew.TimePrimary,
			&aircrew.TimeSecondary,
			&aircrew.TimeInstructor,
			&aircrew.TimeEvaluator,
			&aircrew.TimeOther,
			&aircrew.TotalAircrewDurationDecimal,
			&aircrew.TotalAircrewSorties,
			&aircrew.CondNightTime,
			&aircrew.CondInstrumentTime,
			&aircrew.CondSimInstrumentTime,
			&aircrew.CondNvgTime,
			&aircrew.CondCombatTime,
			&aircrew.CondCombatSortie,
			&aircrew.CondCombatSupportTime,
			&aircrew.CondCombatSupportSortie,
			&aircrew.AircrewRoleType,
		)
		if err != nil {
			log.Printf("Failed to parse aircrew member for flight log: %s \n%s\n", flight_log_id, err.Error())
			return nil, fmt.Errorf("failed to parse aircrew member for flight log: %s", flight_log_id)
		}
		aircrews = append(aircrews, aircrew)
	}
	return aircrews, nil
}

func GetFlightLogComments(txid uuid.UUID, flight_log_id uuid.UUID) ([]types.FlightLogCommentDTO, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetFlightLogComments))
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return nil, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) AS id
			, BIN_TO_UUID(flight_log_id) AS flight_log_id
			, user_id
			, role_name
			, comment
			, created_on
			, updated_on
		FROM flight_log_comments
		WHERE flight_log_id = UUID_TO_BIN(?)
	`
	rows, err := database.Query(query, flight_log_id)
	if err != nil {
		log.Printf("Failed to retrieve comments for flight log: %s \n%s\n", flight_log_id, err.Error())
		return nil, fmt.Errorf("failed to retrieve comments for flight log: %s", flight_log_id)
	}
	defer rows.Close()

	comments := make([]types.FlightLogCommentDTO, 0)
	for rows.Next() {
		var comment types.FlightLogCommentDTO
		err := rows.Scan(
			&comment.ID,
			&comment.FlightLogID,
			&comment.UserID,
			&comment.RoleName,
			&comment.Comment,
			&comment.CreatedOn,
			&comment.UpdatedOn,
		)
		if err != nil {
			log.Printf("Failed to parse comment for flight log: %s \n%s\n", flight_log_id, err.Error())
			return nil, fmt.Errorf("failed to parse comment for flight log: %s", flight_log_id)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func GetFlightlog(txid uuid.UUID, user_id uuid.UUID, flight_log_id uuid.UUID) (types.FlightLogDTO, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetFlightlog))
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return types.FlightLogDTO{}, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) AS id
			, BIN_TO_UUID(user_id) AS user_id
			, mds
			, flight_log_date
			, serial_number
			, unit_charged
			, harm_location
			, flight_authorization
			, issuing_unit
			, is_training_flight
			, is_training_only
			, total_flight_decimal_time
			, scheduler_signature_id
			, sarm_signature_id
			, instructor_signature_id
			, student_signature_id
			, training_officer_signature_id
			, type
			, remarks
		FROM flight_logs
		WHERE id = UUID_TO_BIN(?) AND user_id = UUID_TO_BIN(?)
	`
	row := database.QueryRow(query, flight_log_id, user_id)
	var flight_log_dto types.FlightLogDTO
	err = row.Scan(
		&flight_log_dto.ID,
		&flight_log_dto.UserID,
		&flight_log_dto.MDS,
		&flight_log_dto.FlightLogDate,
		&flight_log_dto.SerialNumber,
		&flight_log_dto.UnitCharged,
		&flight_log_dto.HarmLocation,
		&flight_log_dto.FlightAuthorization,
		&flight_log_dto.IssuingUnit,
		&flight_log_dto.IsTrainingFlight,
		&flight_log_dto.IsTrainingOnly,
		&flight_log_dto.TotalFlightDecimalTime,
		&flight_log_dto.SchedulerSignatureID,
		&flight_log_dto.SarmSignatureID,
		&flight_log_dto.InstructorSignatureID,
		&flight_log_dto.StudentSignatureID,
		&flight_log_dto.TrainingOfficerSignatureID,
		&flight_log_dto.Type,
		&flight_log_dto.Remarks,
	)
	if err != nil {
		log.Printf("Failed to retrieve flight log: %s for user: %s\n%s\n", flight_log_id, user_id, err.Error())
		return types.FlightLogDTO{}, errors.New("failed to retrieve flight log")
	}
	return flight_log_dto, nil
}

func GetFlightlogs(txid uuid.UUID, user_id uuid.UUID) ([]types.FlightLogDTO, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetFlightlogs))
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return nil, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) AS id
			, BIN_TO_UUID(user_id) AS user_id
			, mds
			, flight_log_date
			, serial_number
			, unit_charged
			, harm_location
			, flight_authorization
			, issuing_unit
			, is_training_flight
			, is_training_only
			, total_flight_decimal_time
			, scheduler_signature_id
			, sarm_signature_id
			, instructor_signature_id
			, student_signature_id
			, training_officer_signature_id
			, type
			, remarks
		FROM flight_logs
		WHERE user_id = UUID_TO_BIN(?)
	`
	rows, err := database.Query(query, user_id)
	if err != nil {
		log.Printf("Failed to retrieve flight logs for user: %s\n%s\n", user_id, err.Error())
		return nil, errors.New("failed to retrieve flight logs")
	}
	defer rows.Close()

	flight_logs := make([]types.FlightLogDTO, 0)
	for rows.Next() {
		var flight_log types.FlightLogDTO
		err := rows.Scan(
			&flight_log.ID,
			&flight_log.UserID,
			&flight_log.MDS,
			&flight_log.FlightLogDate,
			&flight_log.SerialNumber,
			&flight_log.UnitCharged,
			&flight_log.HarmLocation,
			&flight_log.FlightAuthorization,
			&flight_log.IssuingUnit,
			&flight_log.IsTrainingFlight,
			&flight_log.IsTrainingOnly,
			&flight_log.TotalFlightDecimalTime,
			&flight_log.SchedulerSignatureID,
			&flight_log.SarmSignatureID,
			&flight_log.InstructorSignatureID,
			&flight_log.StudentSignatureID,
			&flight_log.TrainingOfficerSignatureID,
			&flight_log.Type,
			&flight_log.Remarks,
		)
		if err != nil {
			log.Printf("Failed to parse a flight log for user: %s \n%s\n", user_id, err.Error())
			return nil, fmt.Errorf("failed to parse a flight log for user: %s", user_id)
		}
		flight_logs = append(flight_logs, flight_log)
	}
	return flight_logs, nil
}

func GetFlightlogsAll(txid uuid.UUID, user_id uuid.UUID, where_clause string, where_args []interface{}) ([]types.FlightLogDTO, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetFlightlogsAll))
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return nil, errors.New("failed to connect to DB")
	}
	flight_log_query := `
		SELECT BIN_TO_UUID(flight_logs.id) AS "flight_log_id"
			, BIN_TO_UUID(flight_logs.user_id) AS "user_id"
			, flight_logs.mds
			, flight_logs.flight_log_date
			, flight_logs.serial_number
			, flight_logs.unit_charged
			, flight_logs.harm_location
			, flight_logs.flight_authorization
			, flight_logs.issuing_unit
			, flight_logs.is_training_flight
			, flight_logs.is_training_only
			, flight_logs.total_flight_decimal_time
			, flight_logs.scheduler_signature_id
			, flight_logs.sarm_signature_id
			, flight_logs.instructor_signature_id
			, flight_logs.student_signature_id
			, flight_logs.training_officer_signature_id
			, flight_logs.type
			, flight_logs.remarks
		FROM flight_logs
	`
	// where_clause = strings.ReplaceAll(where_clause, "?", "UUID_TO_BIN(?)")
	flight_log_query_str := strings.Join([]string{flight_log_query, "WHERE", where_clause}, " ")

	/* TODO [drd] remove this logging */
	log.Printf("Query string: %s\n", flight_log_query_str)
	b, _ := json.Marshal(where_args)
	log.Printf("Arguments: %s\n", b)

	rows, err := database.Query(flight_log_query_str, where_args...)
	if err != nil {
		log.Printf("Failed to retrieve flight logs for user: %s\n%s\n", user_id, err.Error())
		return nil, errors.New("failed to retrieve flight logs")
	}
	defer rows.Close()

	flight_logs := make([]types.FlightLogDTO, 0)
	for rows.Next() {
		var flight_log types.FlightLogDTO
		err := rows.Scan(
			&flight_log.ID,
			&flight_log.UserID,
			&flight_log.MDS,
			&flight_log.FlightLogDate,
			&flight_log.SerialNumber,
			&flight_log.UnitCharged,
			&flight_log.HarmLocation,
			&flight_log.FlightAuthorization,
			&flight_log.IssuingUnit,
			&flight_log.IsTrainingFlight,
			&flight_log.IsTrainingOnly,
			&flight_log.TotalFlightDecimalTime,
			&flight_log.SchedulerSignatureID,
			&flight_log.SarmSignatureID,
			&flight_log.InstructorSignatureID,
			&flight_log.StudentSignatureID,
			&flight_log.TrainingOfficerSignatureID,
			&flight_log.Type,
			&flight_log.Remarks,
		)
		if err != nil {
			log.Printf("Failed to parse a flight log for user: %s \n%s\n", user_id, err.Error())
			return nil, fmt.Errorf("failed to parse a flight log for user: %s", user_id)
		}
		flight_logs = append(flight_logs, flight_log)
	}
	return flight_logs, nil
}

func GetMissions(txid uuid.UUID, flight_log_id uuid.UUID) ([]types.FlightLogMissionDTO, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetMissions))
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return nil, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) AS id
			, BIN_TO_UUID(flight_log_id) AS flight_log_id
			, mission_number
			, mission_symbol
			, mission_from
			, mission_to
			, takeoff_time
			, land_time
			, total_time_decimal
			, total_time_display
			, touch_and_gos
			, full_stops
			, total_landings
			, sorties
		FROM missions
		WHERE flight_log_id = UUID_TO_BIN(?)
	`
	rows, err := database.Query(query, flight_log_id)
	if err != nil {
		log.Printf("Failed to retrieve missions for flight log: %s \n%s\n", flight_log_id, err.Error())
		return nil, fmt.Errorf("failed to retrieve missions for flight log: %s", flight_log_id)
	}
	defer rows.Close()

	missions := make([]types.FlightLogMissionDTO, 0)
	for rows.Next() {
		var mission types.FlightLogMissionDTO
		err := rows.Scan(
			&mission.ID,
			&mission.FlightLogID,
			&mission.MissionNumber,
			&mission.MissionSymbol,
			&mission.MissionFrom,
			&mission.MissionTo,
			&mission.TakeoffTime,
			&mission.LandTime,
			&mission.TotalTimeDecimal,
			&mission.TotalTimeDisplay,
			&mission.TouchAndGos,
			&mission.FullStops,
			&mission.TotalLandings,
			&mission.Sorties,
		)
		if err != nil {
			log.Printf("Failed to parse a mission leg for flight log: %s \n%s\n", flight_log_id, err.Error())
			return nil, fmt.Errorf("failed to parse a mission leg for flight log: %s", flight_log_id)
		}
		missions = append(missions, mission)
	}
	return missions, nil
}

func InsertAircrews(txid uuid.UUID, flight_log types.FlightLogDTO) ([]uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(InsertAircrews))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return nil, errors.New(err_string)
	}
	ids := []uuid.UUID{}
	for _, aircrew := range flight_log.Aircrew {
		query := `
			INSERT INTO aircrews
			(
				id
				, flight_log_id
				, user_id
				, flying_origin
				, flight_auth_code
				, time_primary
				, time_secondary
				, time_instructor
				, time_evaluator
				, time_other
				, total_aircrew_duration_decimal
				, total_aircrew_sorties
				, cond_night_time
				, cond_instrument_time
				, cond_sim_instrument_time
				, cond_nvg_time
				, cond_combat_time
				, cond_combat_sortie
				, cond_combat_support_time
				, cond_combat_support_sortie
				, aircrew_role_type
			)
			VALUES
			(
				UUID_TO_BIN(?) -- id
				, UUID_TO_BIN(?) -- flight_log_id
				, UUID_TO_BIN(?) -- user_id
				, ? -- flying_origin
				, ? -- flight_auth_code
				, ? -- time_primary
				, ? -- time_secondary
				, ? -- time_instructor
				, ? -- time_evaluator
				, ? -- time_other
				, ? -- total_aircrew_duration_decimal
				, ? -- total_aircrew_sorties
				, ? -- cond_night_time
				, ? -- cond_instrument_time
				, ? -- cond_sim_instrument_time
				, ? -- cond_nvg_time
				, ? -- cond_combat_time
				, ? -- cond_combat_sortie
				, ? -- cond_combat_support_time
				, ? -- cond_combat_support_sortie
				, ? -- aircrew_role_type
			)
		`
		id := uuid.New()
		_, err = database.Exec(
			query,
			id,
			flight_log.ID,
			aircrew.UserID,
			aircrew.FlyingOrigin,
			aircrew.FlightAuthCode,
			aircrew.TimePrimary,
			aircrew.TimeSecondary,
			aircrew.TimeInstructor,
			aircrew.TimeEvaluator,
			aircrew.TimeOther,
			aircrew.TotalAircrewDurationDecimal,
			aircrew.TotalAircrewSorties,
			aircrew.CondNightTime,
			aircrew.CondInstrumentTime,
			aircrew.CondSimInstrumentTime,
			aircrew.CondNvgTime,
			aircrew.CondCombatTime,
			aircrew.CondCombatSortie,
			aircrew.CondCombatSupportTime,
			aircrew.CondCombatSupportSortie,
			aircrew.AircrewRoleType,
		)
		if err != nil {
			log.Printf("failed aircrew insert\n%s\n", err.Error())
			return nil, errors.New(err_string)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func InsertFlightLog(txid uuid.UUID, request_user_id uuid.UUID, flight_log types.FlightLogDTO) (uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(InsertFlightLog))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return uuid.Nil, errors.New(err_string)
	}
	query := `
		INSERT INTO flight_logs
		(
			id
			, user_id
			, mds
			, flight_log_date
			, serial_number
			, unit_charged
			, harm_location
			, flight_authorization
			, issuing_unit
			, is_training_flight
			, is_training_only
			, total_flight_decimal_time
			, scheduler_signature_id
			, sarm_signature_id
			, instructor_signature_id
			, student_signature_id
			, training_officer_signature_id
			, type
			, remarks
		)
		VALUES
		(
			UUID_TO_BIN(?), -- id
			UUID_TO_BIN(?), -- user_id
			?, -- mds
			?, -- flight_log_date
			?, -- serial_number
			?, -- unit_charged
			?, -- harm_location
			?, -- flight_authorization
			?, -- issuing_unit
			?, -- is_training_flight
			?, -- is_training_only
			?, -- total_flight_decimal_time
			UUID_TO_BIN(?), -- scheduler_signature_id
			UUID_TO_BIN(?), -- sarm_signature_id
			UUID_TO_BIN(?), -- instructor_signature_id
			UUID_TO_BIN(?), -- student_signature_id
			UUID_TO_BIN(?), -- training_officer_signature_id
			?, -- type
			? -- remarks
		)
	`
	id := uuid.New()
	_, err = database.Exec(
		query,
		id,
		request_user_id,
		flight_log.MDS,
		flight_log.FlightLogDate,
		flight_log.SerialNumber,
		flight_log.UnitCharged,
		flight_log.HarmLocation,
		flight_log.FlightAuthorization,
		flight_log.IssuingUnit,
		flight_log.IsTrainingFlight,
		flight_log.IsTrainingOnly,
		flight_log.TotalFlightDecimalTime,
		flight_log.SchedulerSignatureID,
		flight_log.SarmSignatureID,
		flight_log.InstructorSignatureID,
		flight_log.StudentSignatureID,
		flight_log.TrainingOfficerSignatureID,
		flight_log.Type,
		flight_log.Remarks,
	)
	if err != nil {
		log.Printf("failed flight log insert\n%s\n", err.Error())
		return uuid.Nil, errors.New(err_string)
	}
	return id, nil
}

func InsertFlightLogComment(txid uuid.UUID, request_user_id uuid.UUID, flight_log_comment types.FlightLogCommentDTO) (uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(InsertFlightLogComment))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return uuid.Nil, errors.New(err_string)
	}
	query := `
		INSERT INTO flight_log_comments
		(
			id
			, flight_log_id
			, user_id
			, role_name
			, comment
		)
		VALUES
		(
			UUID_TO_BIN(?), -- id
			UUID_TO_BIN(?), -- flight_log_id
			UUID_TO_BIN(?), -- user_id
			?, -- role_name
			? -- comment
		)
	`
	id := uuid.New()
	_, err = database.Exec(
		query,
		id,
		flight_log_comment.FlightLogID,
		request_user_id,
		flight_log_comment.RoleName,
		flight_log_comment.Comment,
	)
	if err != nil {
		log.Printf("failed flight log comment insert\n%s\n", err.Error())
		return uuid.Nil, errors.New(err_string)
	}
	return id, nil
}

func InsertFlightLogComments(txid uuid.UUID, request_user_id uuid.UUID, flight_log types.FlightLogDTO) ([]uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(InsertFlightLogComments))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return nil, errors.New(err_string)
	}
	ids := []uuid.UUID{}
	for _, comment := range flight_log.Comments {
		query := `
			INSERT INTO flight_log_comments
			(
				id
				, flight_log_id
				, user_id
				, role_name
				, comment
			)
			VALUES
			(
				UUID_TO_BIN(?), -- id
				UUID_TO_BIN(?), -- flight_log_id
				UUID_TO_BIN(?), -- user_id
				?, -- role_name
				? -- comment
			)
		`
		id := uuid.New()
		_, err = database.Exec(
			query,
			id,
			flight_log.ID,
			request_user_id,
			comment.RoleName,
			comment.Comment,
		)
		if err != nil {
			log.Printf("failed flight log comment insert\n%s\n", err.Error())
			return nil, errors.New(err_string)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func InsertMissions(txid uuid.UUID, flight_log types.FlightLogDTO) ([]uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(InsertMissions))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return nil, errors.New(err_string)
	}
	ids := []uuid.UUID{}
	for _, mission := range flight_log.Missions {
		query := `
			INSERT INTO missions
			(
				id
				, flight_log_id
				, mission_number
				, mission_symbol
				, mission_from
				, mission_to
				, takeoff_time
				, land_time
				, total_time_decimal
				, total_time_display
				, touch_and_gos
				, full_stops
				, total_landings
				, sorties
			)
			VALUES
			(
				UUID_TO_BIN(?), -- id
				UUID_TO_BIN(?), -- flight_log_id
				?, -- mission_number
				?, -- mission_symbol
				?, -- mission_from
				?, -- mission_to
				?, -- takeoff_time
				?, -- land_time
				?, -- total_time_decimal
				?, -- total_time_display
				?, -- touch_and_gos
				?, -- full_stops
				?, -- total_landings
				? -- sorties
			)
		`
		id := uuid.New()
		_, err = database.Exec(
			query,
			id,
			flight_log.ID,
			mission.MissionNumber,
			mission.MissionSymbol,
			mission.MissionFrom,
			mission.MissionTo,
			mission.TakeoffTime,
			mission.LandTime,
			mission.TotalTimeDecimal,
			mission.TotalTimeDisplay,
			mission.TouchAndGos,
			mission.FullStops,
			mission.TotalLandings,
			mission.Sorties,
		)
		if err != nil {
			log.Printf("failed mission insert\n%s\n", err.Error())
			return nil, errors.New(err_string)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func UpdateAircrews(txid uuid.UUID, flight_log types.FlightLogDTO) ([]uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(UpdateAircrews))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())

	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return nil, errors.New(err_string)
	}

	ids := []uuid.UUID{}
	for _, aircrew := range flight_log.Aircrew {
		query := `
			UPDATE aircrews
			SET
				flight_log_id = UUID_TO_BIN(?)
				, user_id = UUID_TO_BIN(?)
				, flying_origin = ?
				, flight_auth_code = ?
				, time_primary = ?
				, time_secondary = ?
				, time_instructor = ?
				, time_evaluator = ?
				, time_other = ?
				, total_aircrew_duration_decimal = ?
				, total_aircrew_sorties = ?
				, cond_night_time = ?
				, cond_instrument_time = ?
				, cond_sim_instrument_time = ?
				, cond_nvg_time = ?
				, cond_combat_time = ?
				, cond_combat_sortie = ?
				, cond_combat_support_time = ?
				, cond_combat_support_sortie = ?
				, aircrew_role_type = ?
			WHERE id = UUID_TO_BIN(?)
		`
		_, err := database.Exec(
			query,
			flight_log.ID,
			aircrew.UserID,
			aircrew.FlyingOrigin,
			aircrew.FlightAuthCode,
			aircrew.TimePrimary,
			aircrew.TimeSecondary,
			aircrew.TimeInstructor,
			aircrew.TimeEvaluator,
			aircrew.TimeOther,
			aircrew.TotalAircrewDurationDecimal,
			aircrew.TotalAircrewSorties,
			aircrew.CondNightTime,
			aircrew.CondInstrumentTime,
			aircrew.CondSimInstrumentTime,
			aircrew.CondNvgTime,
			aircrew.CondCombatTime,
			aircrew.CondCombatSortie,
			aircrew.CondCombatSupportTime,
			aircrew.CondCombatSupportSortie,
			aircrew.AircrewRoleType,
			// WHERE clause
			aircrew.ID,
		)
		if err != nil {
			log.Printf("failed aircrew update\n%s\n", err.Error())
			return nil, errors.New(err_string)
		}
		ids = append(ids, aircrew.ID)
	}
	return ids, nil
}

func UpdateFlightLog(txid uuid.UUID, flight_log types.FlightLogDTO) (uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(UpdateFlightLog))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())

	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return uuid.Nil, errors.New(err_string)
	}
	query := `
		UPDATE flight_logs
		SET
			mds = ?
			, flight_log_date = ?
			, serial_number = ?
			, unit_charged = ?
			, harm_location = ?
			, flight_authorization = ?
			, issuing_unit = ?
			, is_training_flight = ?
			, is_training_only = ?
			, total_flight_decimal_time = ?
			, scheduler_signature_id = UUID_TO_BIN(?)
			, sarm_signature_id = UUID_TO_BIN(?)
			, instructor_signature_id = UUID_TO_BIN(?)
			, student_signature_id = UUID_TO_BIN(?)
			, training_officer_signature_id = UUID_TO_BIN(?)
			, type = ?
			, remarks = ?
		WHERE id = UUID_TO_BIN(?)
	`
	_, err = database.Exec(
		query,
		flight_log.MDS,
		flight_log.FlightLogDate,
		flight_log.SerialNumber,
		flight_log.UnitCharged,
		flight_log.HarmLocation,
		flight_log.FlightAuthorization,
		flight_log.IssuingUnit,
		flight_log.IsTrainingFlight,
		flight_log.IsTrainingOnly,
		flight_log.TotalFlightDecimalTime,
		flight_log.SchedulerSignatureID,
		flight_log.SarmSignatureID,
		flight_log.InstructorSignatureID,
		flight_log.StudentSignatureID,
		flight_log.TrainingOfficerSignatureID,
		flight_log.Type,
		flight_log.Remarks,
		// WHERE clause
		flight_log.ID,
	)
	if err != nil {
		log.Printf("failed flight log update\n%s\n", err.Error())
		return uuid.Nil, errors.New(err_string)
	}
	return flight_log.ID, nil
}

func UpdateMissions(txid uuid.UUID, flight_log types.FlightLogDTO) ([]uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(UpdateMissions))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())

	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return nil, errors.New(err_string)
	}

	ids := []uuid.UUID{}
	for _, mission := range flight_log.Missions {
		query := `
			UPDATE missions
			SET
				flight_log_id = UUID_TO_BIN(?)
				, mission_number = ?
				, mission_symbol = ?
				, mission_from = ?
				, mission_to = ?
				, takeoff_time = ?
				, land_time = ?
				, total_time_decimal = ?
				, total_time_display = ?
				, touch_and_gos = ?
				, full_stops = ?
				, total_landings = ?
				, sorties = ?
			WHERE id = UUID_TO_BIN(?)
		`
		_, err := database.Exec(
			query,
			flight_log.ID,
			mission.MissionNumber,
			mission.MissionSymbol,
			mission.MissionFrom,
			mission.MissionTo,
			mission.TakeoffTime,
			mission.LandTime,
			mission.TotalTimeDecimal,
			mission.TotalTimeDisplay,
			mission.TouchAndGos,
			mission.FullStops,
			mission.TotalLandings,
			mission.Sorties,
			// WHERE clause
			mission.ID,
		)
		if err != nil {
			log.Printf("failed mission update\n%s\n", err.Error())
			return nil, errors.New(err_string)
		}
		ids = append(ids, mission.ID)
	}
	return ids, nil
}
