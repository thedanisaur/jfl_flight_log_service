package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/thedanisaur/jfl_platform/types"
	"github.com/thedanisaur/jfl_platform/util"

	"github.com/google/uuid"
)

func GetAirCrews(txid uuid.UUID, flight_log_id uuid.UUID) ([]types.AircrewDTO, error) {
	log.Printf("%s | %s\n", util.GetFunctionName(GetAirCrews), txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return nil, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) AS id
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
		FROM aircrews
		WHERE flight_log_id = UUID_TO_BIN(?)
	`
	rows, err := database.Query(query, flight_log_id)
	if err != nil {
		log.Printf("Failed to retrieve aircrew members for flight log: %s \n%s\n", flight_log_id, err.Error())
		return nil, fmt.Errorf("failed to retrieve aircrew members for flight log: %s", flight_log_id)
	}
	defer rows.Close()

	aircrews := make([]types.AircrewDTO, 0)
	for rows.Next() {
		var aircrew types.AircrewDTO
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

func GetFlightLogComments(txid uuid.UUID, flight_log_id uuid.UUID) ([]types.CommentDTO, error) {
	log.Printf("%s | %s\n", util.GetFunctionName(GetFlightLogComments), txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return nil, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) AS id
			, flight_log_id
			, user_id
			, role_id
			, comment
		FROM flight_log_comments
		WHERE flight_log_id = UUID_TO_BIN(?)
	`
	rows, err := database.Query(query, flight_log_id)
	if err != nil {
		log.Printf("Failed to retrieve comments for flight log: %s \n%s\n", flight_log_id, err.Error())
		return nil, fmt.Errorf("failed to retrieve comments for flight log: %s", flight_log_id)
	}
	defer rows.Close()

	comments := make([]types.CommentDTO, 0)
	for rows.Next() {
		var comment types.CommentDTO
		err := rows.Scan(
			&comment.ID,
			&comment.FlightLogID,
			&comment.UserID,
			&comment.RoleID,
			&comment.Comment,
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
	log.Printf("%s | %s\n", util.GetFunctionName(GetFlightlog), txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return types.FlightLogDTO{}, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) AS id
			, user_id
			, unit_id
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
		&flight_log_dto.UnitID,
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

func GetMissions(txid uuid.UUID, flight_log_id uuid.UUID) ([]types.MissionDTO, error) {
	log.Printf("%s | %s\n", util.GetFunctionName(GetMissions), txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return nil, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) AS id
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
		FROM missions
		WHERE flight_log_id = UUID_TO_BIN(?)
	`
	rows, err := database.Query(query, flight_log_id)
	if err != nil {
		log.Printf("Failed to retrieve missions for flight log: %s \n%s\n", flight_log_id, err.Error())
		return nil, fmt.Errorf("failed to retrieve missions for flight log: %s", flight_log_id)
	}
	defer rows.Close()

	missions := make([]types.MissionDTO, 0)
	for rows.Next() {
		var mission types.MissionDTO
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

func InsertAircrew(txid uuid.UUID, flight_log types.FlightLogDTO) ([]uuid.UUID, error) {
	log.Printf("%s | %s\n", util.GetFunctionName(InsertAircrew), txid.String())
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
				UUID_TO_BIN(?), -- id
				UUID_TO_BIN(?), -- flight_log_id
				UUID_TO_BIN(?), -- user_id
				? , -- flying_origin
				? , -- flight_auth_code
				? , -- time_primary
				? , -- time_secondary
				? , -- time_instructor
				? , -- time_evaluator
				? , -- time_other
				? , -- total_aircrew_duration_decimal
				? , -- total_aircrew_sorties
				? , -- cond_night_time
				? , -- cond_instrument_time
				? , -- cond_sim_instrument_time
				? , -- cond_nvg_time
				? , -- cond_combat_time
				? , -- cond_combat_sortie
				? , -- cond_combat_support_time
				? , -- cond_combat_support_sortie
				? -- aircrew_role_type
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
	log.Printf("%s | %s\n", util.GetFunctionName(InsertFlightLog), txid.String())
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
			, unit_id
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
			UUID_TO_BIN(?), -- unit_id
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
		flight_log.UnitID,
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

func InsertFlightLogComments(txid uuid.UUID, request_user_id uuid.UUID, flight_log types.FlightLogDTO) ([]uuid.UUID, error) {
	log.Printf("%s | %s\n", util.GetFunctionName(InsertFlightLogComments), txid.String())
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
				, role_id
				, comment
			)
			VALUES
			(
				UUID_TO_BIN(?), -- id
				UUID_TO_BIN(?), -- flight_log_id
				UUID_TO_BIN(?), -- user_id
				UUID_TO_BIN(?), -- role_id
				? -- comment
			)
		`
		id := uuid.New()
		_, err = database.Exec(
			query,
			id,
			flight_log.ID,
			request_user_id,
			comment.RoleID,
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
	log.Printf("%s | %s\n", util.GetFunctionName(InsertMissions), txid.String())
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
