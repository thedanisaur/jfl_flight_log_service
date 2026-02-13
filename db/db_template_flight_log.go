package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/thedanisaur/jfl_platform/types"
	"github.com/thedanisaur/jfl_platform/util"

	"github.com/google/uuid"
)

func DeleteTemplateFlightlog(txid uuid.UUID, user_id uuid.UUID, template_id uuid.UUID) (uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(DeleteTemplateFlightlog))
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

	// Delete template flight log's aircrew records
	aircrews_query := `DELETE FROM template_aircrews WHERE flight_log_id = UUID_TO_BIN(?)`
	aircrews_result, err := database.Exec(aircrews_query, template_id)
	if err != nil {
		log.Printf("Failed to delete template aircrews: %s for user: %s\n%s\n", template_id, user_id, err.Error())
		return uuid.Nil, errors.New("failed to delete template aircrews")
	}
	_, err = aircrews_result.RowsAffected()
	if err != nil {
		return uuid.Nil, errors.New("failed to delete template aircrews")
	}

	// Delete template flight log's mission records
	missions_query := `DELETE FROM template_missions WHERE flight_log_id = UUID_TO_BIN(?)`
	missions_result, err := database.Exec(missions_query, template_id)
	if err != nil {
		log.Printf("Failed to delete template missions: %s for user: %s\n%s\n", template_id, user_id, err.Error())
		return uuid.Nil, errors.New("failed to delete template missions")
	}
	_, err = missions_result.RowsAffected()
	if err != nil {
		return uuid.Nil, errors.New("failed to delete template missions")
	}

	// Delete template flight log
	template_query := `DELETE FROM template_flight_logs WHERE id = UUID_TO_BIN(?)`
	template_result, err := database.Exec(template_query, template_id)
	if err != nil {
		log.Printf("Failed to delete template flight log: %s for user: %s\n%s\n", template_id, user_id, err.Error())
		return uuid.Nil, errors.New("failed to delete template flight log")
	}
	_, err = template_result.RowsAffected()
	if err != nil {
		return uuid.Nil, errors.New("failed to delete template flight log")
	}

	return template_id, nil
}

func GetTemplateAirCrews(txid uuid.UUID, template_id uuid.UUID) ([]types.FlightLogAircrewDTO, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetTemplateAirCrews))
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
		FROM template_aircrews
		WHERE flight_log_id = UUID_TO_BIN(?)
	`
	rows, err := database.Query(query, template_id)
	if err != nil {
		log.Printf("Failed to retrieve template aircrew members for template flight log: %s \n%s\n", template_id, err.Error())
		return nil, fmt.Errorf("failed to retrieve template aircrew members for template flight log: %s", template_id)
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
			log.Printf("Failed to parse template aircrew member for template flight log: %s \n%s\n", template_id, err.Error())
			return nil, fmt.Errorf("failed to parse template aircrew member for template flight log: %s", template_id)
		}
		aircrews = append(aircrews, aircrew)
	}
	return aircrews, nil
}

func GetTemplateFlightlog(txid uuid.UUID, user_id uuid.UUID, template_id uuid.UUID) (types.TemplateFlightLogDTO, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetTemplateFlightlog))
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return types.TemplateFlightLogDTO{}, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) AS id
			, name
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
		FROM template_flight_logs
		WHERE id = UUID_TO_BIN(?) AND user_id = UUID_TO_BIN(?)
	`
	row := database.QueryRow(query, template_id, user_id)
	var template_flight_log_dto types.TemplateFlightLogDTO
	err = row.Scan(
		&template_flight_log_dto.ID,
		&template_flight_log_dto.Name,
		&template_flight_log_dto.UserID,
		&template_flight_log_dto.MDS,
		&template_flight_log_dto.FlightLogDate,
		&template_flight_log_dto.SerialNumber,
		&template_flight_log_dto.UnitCharged,
		&template_flight_log_dto.HarmLocation,
		&template_flight_log_dto.FlightAuthorization,
		&template_flight_log_dto.IssuingUnit,
		&template_flight_log_dto.IsTrainingFlight,
		&template_flight_log_dto.IsTrainingOnly,
		&template_flight_log_dto.TotalFlightDecimalTime,
		&template_flight_log_dto.SchedulerSignatureID,
		&template_flight_log_dto.SarmSignatureID,
		&template_flight_log_dto.InstructorSignatureID,
		&template_flight_log_dto.StudentSignatureID,
		&template_flight_log_dto.TrainingOfficerSignatureID,
		&template_flight_log_dto.Type,
		&template_flight_log_dto.Remarks,
	)
	if err != nil {
		log.Printf("Failed to retrieve template flight log: %s for user: %s\n%s\n", template_id, user_id, err.Error())
		return types.TemplateFlightLogDTO{}, errors.New("failed to retrieve template flight log")
	}
	return template_flight_log_dto, nil
}

func GetTemplateFlightlogs(txid uuid.UUID, user_id uuid.UUID) ([]types.TemplateFlightLogDTO, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetTemplateFlightlogs))
	database, err := GetInstance()
	if err != nil {
		log.Printf("Failed to connect to DB\n%s\n", err.Error())
		return nil, errors.New("failed to connect to DB")
	}
	query := `
		SELECT BIN_TO_UUID(id) AS id
			, name
			, BIN_TO_UUID(user_id) AS user_id
		FROM template_flight_logs
		WHERE user_id = UUID_TO_BIN(?)
	`
	rows, err := database.Query(query, user_id)
	if err != nil {
		log.Printf("Failed to retrieve template flight logs for user: %s\n%s\n", user_id, err.Error())
		return nil, errors.New("failed to retrieve template flight logs")
	}
	defer rows.Close()

	template_flight_logs := make([]types.TemplateFlightLogDTO, 0)
	for rows.Next() {
		var template_flight_log types.TemplateFlightLogDTO
		err := rows.Scan(
			&template_flight_log.ID,
			&template_flight_log.Name,
			&template_flight_log.UserID,
		)
		if err != nil {
			log.Printf("Failed to parse a template flight log for user: %s \n%s\n", user_id, err.Error())
			return nil, fmt.Errorf("failed to parse a template flight log for user: %s", user_id)
		}
		template_flight_logs = append(template_flight_logs, template_flight_log)
	}
	return template_flight_logs, nil
}

func GetTemplateMissions(txid uuid.UUID, template_id uuid.UUID) ([]types.FlightLogMissionDTO, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetTemplateMissions))
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
		FROM template_missions
		WHERE flight_log_id = UUID_TO_BIN(?)
	`
	rows, err := database.Query(query, template_id)
	if err != nil {
		log.Printf("Failed to retrieve template missions for template flight log: %s \n%s\n", template_id, err.Error())
		return nil, fmt.Errorf("failed to retrieve template missions for template flight log: %s", template_id)
	}
	defer rows.Close()

	template_missions := make([]types.FlightLogMissionDTO, 0)
	for rows.Next() {
		var template_mission types.FlightLogMissionDTO
		err := rows.Scan(
			&template_mission.ID,
			&template_mission.FlightLogID,
			&template_mission.MissionNumber,
			&template_mission.MissionSymbol,
			&template_mission.MissionFrom,
			&template_mission.MissionTo,
			&template_mission.TakeoffTime,
			&template_mission.LandTime,
			&template_mission.TotalTimeDecimal,
			&template_mission.TotalTimeDisplay,
			&template_mission.TouchAndGos,
			&template_mission.FullStops,
			&template_mission.TotalLandings,
			&template_mission.Sorties,
		)
		if err != nil {
			log.Printf("Failed to parse a template mission leg for template flight log: %s \n%s\n", template_id, err.Error())
			return nil, fmt.Errorf("failed to parse a template mission leg for template flight log: %s", template_id)
		}
		template_missions = append(template_missions, template_mission)
	}
	return template_missions, nil
}

func InsertTemplateAircrews(txid uuid.UUID, flight_log types.TemplateFlightLogDTO) ([]uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(InsertTemplateAircrews))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return nil, errors.New(err_string)
	}
	ids := []uuid.UUID{}
	for _, aircrew := range flight_log.Aircrew {
		query := `
			INSERT INTO template_aircrews
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
			log.Printf("failed template aircrew insert\n%s\n", err.Error())
			return nil, errors.New(err_string)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func InsertTemplateFlightLog(txid uuid.UUID, request_user_id uuid.UUID, template_flight_log types.TemplateFlightLogDTO) (uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(InsertTemplateFlightLog))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return uuid.Nil, errors.New(err_string)
	}
	query := `
		INSERT INTO template_flight_logs
		(
			id
			, name
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
			?, -- name
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
		template_flight_log.Name,
		request_user_id,
		template_flight_log.MDS,
		template_flight_log.FlightLogDate,
		template_flight_log.SerialNumber,
		template_flight_log.UnitCharged,
		template_flight_log.HarmLocation,
		template_flight_log.FlightAuthorization,
		template_flight_log.IssuingUnit,
		template_flight_log.IsTrainingFlight,
		template_flight_log.IsTrainingOnly,
		template_flight_log.TotalFlightDecimalTime,
		template_flight_log.SchedulerSignatureID,
		template_flight_log.SarmSignatureID,
		template_flight_log.InstructorSignatureID,
		template_flight_log.StudentSignatureID,
		template_flight_log.TrainingOfficerSignatureID,
		template_flight_log.Type,
		template_flight_log.Remarks,
	)
	if err != nil {
		log.Printf("failed template flight log insert\n%s\n", err.Error())
		return uuid.Nil, errors.New(err_string)
	}
	return id, nil
}

func InsertTemplateMissions(txid uuid.UUID, flight_log types.TemplateFlightLogDTO) ([]uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(InsertTemplateMissions))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())
	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return nil, errors.New(err_string)
	}
	ids := []uuid.UUID{}
	for _, mission := range flight_log.Missions {
		query := `
			INSERT INTO template_missions
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
			log.Printf("failed template mission insert\n%s\n", err.Error())
			return nil, errors.New(err_string)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func UpdateTemplateAircrews(txid uuid.UUID, flight_log types.TemplateFlightLogDTO) ([]uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(UpdateTemplateAircrews))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())

	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return nil, errors.New(err_string)
	}

	ids := []uuid.UUID{}
	for _, aircrew := range flight_log.Aircrew {
		query := `
			UPDATE template_aircrews
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
			log.Printf("failed template aircrew update\n%s\n", err.Error())
			return nil, errors.New(err_string)
		}
		ids = append(ids, aircrew.ID)
	}
	return ids, nil
}

func UpdateTemplateFlightLog(txid uuid.UUID, template_flight_log types.TemplateFlightLogDTO) (uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(UpdateTemplateFlightLog))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())

	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return uuid.Nil, errors.New(err_string)
	}
	query := `
		UPDATE template_flight_logs
		SET
			name = ?
			, mds = ?
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
		template_flight_log.Name,
		template_flight_log.MDS,
		template_flight_log.FlightLogDate,
		template_flight_log.SerialNumber,
		template_flight_log.UnitCharged,
		template_flight_log.HarmLocation,
		template_flight_log.FlightAuthorization,
		template_flight_log.IssuingUnit,
		template_flight_log.IsTrainingFlight,
		template_flight_log.IsTrainingOnly,
		template_flight_log.TotalFlightDecimalTime,
		template_flight_log.SchedulerSignatureID,
		template_flight_log.SarmSignatureID,
		template_flight_log.InstructorSignatureID,
		template_flight_log.StudentSignatureID,
		template_flight_log.TrainingOfficerSignatureID,
		template_flight_log.Type,
		template_flight_log.Remarks,
		// WHERE clause
		template_flight_log.ID,
	)
	if err != nil {
		log.Printf("failed template flight log update\n%s\n", err.Error())
		return uuid.Nil, errors.New(err_string)
	}
	return template_flight_log.ID, nil
}

func UpdateTemplateMissions(txid uuid.UUID, flight_log types.TemplateFlightLogDTO) ([]uuid.UUID, error) {
	log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(UpdateTemplateMissions))
	err_string := fmt.Sprintf("database error: %s\n", txid.String())

	database, err := GetInstance()
	if err != nil {
		log.Printf("failed to connect to database\n%s\n", err.Error())
		return nil, errors.New(err_string)
	}

	ids := []uuid.UUID{}
	for _, mission := range flight_log.Missions {
		query := `
			UPDATE template_missions
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
			log.Printf("failed template mission update\n%s\n", err.Error())
			return nil, errors.New(err_string)
		}
		ids = append(ids, mission.ID)
	}
	return ids, nil
}
