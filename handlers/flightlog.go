package handlers

import (
	"fmt"
	"log"

	"flight_log_service/db"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/thedanisaur/jfl_platform/auth"
	"github.com/thedanisaur/jfl_platform/types"
	"github.com/thedanisaur/jfl_platform/util"
)

func CreateFlightlog(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(CreateFlightlog))

		var flight_log types.FlightLogDTO
		err := c.BodyParser(&flight_log)
		if err != nil {
			log.Printf("Failed to parse flight log data\n%s\n", err.Error())
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Failed to parse flight log data: %s\n", txid.String()))
		}
		/* Get the requesting user */
		request_user := c.Locals("user_claims").(types.UserClaims)

		/* Now start inserting the flight log */
		flight_log.ID, err = db.InsertFlightLog(txid, request_user.UserID, flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		mission_ids, err := db.InsertMissions(txid, flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		aircrew_ids, err := db.InsertAircrews(txid, flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		comment_ids, err := db.InsertFlightLogComments(txid, request_user.UserID, flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		response := fiber.Map{
			"txid":          txid.String(),
			"flight_log_id": flight_log.ID.String(),
			"mission_ids":   mission_ids,
			"aircrew_ids":   aircrew_ids,
			"comment_ids":   comment_ids,
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func GetFlightlog(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetFlightlog))

		user_id, err := uuid.Parse(c.Params("user_id"))
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString("invalid user")
		}
		flight_log_id, err := uuid.Parse(c.Params("flight_log_id"))
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString("invalid flight log")
		}

		flight_log, err := db.GetFlightlog(txid, user_id, flight_log_id)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}

		flight_log.Missions, err = db.GetMissions(txid, flight_log_id)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}

		flight_log.Aircrew, err = db.GetAirCrews(txid, flight_log_id)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}

		flight_log.Comments, err = db.GetFlightLogComments(txid, flight_log_id)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}

		// response := fiber.Map{
		// 	"txid": txid.String(),
		// }

		return c.Status(fiber.StatusOK).JSON(flight_log)
	}
}

func GetFlightlogs(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetFlightlogs))

		user_id, err := uuid.Parse(c.Params("user_id"))
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString("invalid user")
		}

		flight_logs, err := db.GetFlightlogs(txid, user_id)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}

		for index := range flight_logs {
			flight_logs[index].Missions, err = db.GetMissions(txid, flight_logs[index].ID)
			if err != nil {
				return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
			}

			flight_logs[index].Aircrew, err = db.GetAirCrews(txid, flight_logs[index].ID)
			if err != nil {
				return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
			}

			flight_logs[index].Comments, err = db.GetFlightLogComments(txid, flight_logs[index].ID)
			if err != nil {
				return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
			}
		}

		// response := fiber.Map{
		// 	"txid": txid.String(),
		// }

		return c.Status(fiber.StatusOK).JSON(flight_logs)
	}
}

func GetFlightlogsAll(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetFlightlogsAll))

		/* Get the requesting user's info */
		request_user := c.Locals("user_claims").(types.UserClaims)
		resource := "flight-logs"
		operation := "read"

		/* Load Permissions */
		policies, err := db.LoadPermissions(txid, request_user.RoleName, resource, operation)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}
		if len(policies) <= 0 {
			return c.Status(fiber.StatusUnauthorized).SendString("not authorized")
		}

		/* Authorize */
		where_clause, arguments, err := auth.EvaluateRead(txid, resource, operation, "flight_logs_vw", request_user, policies)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}

		/* Get the flight logs */
		flight_logs, err := db.GetFlightlogsAll(txid, request_user.UserID, where_clause, arguments)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}

		for index := range flight_logs {
			flight_logs[index].Missions, err = db.GetMissions(txid, flight_logs[index].ID)
			if err != nil {
				return c.Status(fiber.StatusNotFound).SendString(err.Error())
			}

			flight_logs[index].Aircrew, err = db.GetAirCrews(txid, flight_logs[index].ID)
			if err != nil {
				return c.Status(fiber.StatusNotFound).SendString(err.Error())
			}

			flight_logs[index].Comments, err = db.GetFlightLogComments(txid, flight_logs[index].ID)
			if err != nil {
				return c.Status(fiber.StatusNotFound).SendString(err.Error())
			}
		}

		// TODO [drd] include txid in the response
		// response := fiber.Map{
		// 	"txid": txid.String(),
		// }

		return c.Status(fiber.StatusOK).JSON(flight_logs)
	}
}

func UpdateFlightlog(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(UpdateFlightlog))

		var flight_log types.FlightLogDTO
		err := c.BodyParser(&flight_log)
		if err != nil {
			log.Printf("Failed to parse flight log data\n%s\n", err.Error())
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Failed to parse flight log data: %s\n", txid.String()))
		}
		// TODO [drd] validate that this action is allowed.
		/* Get the requesting user's id */
		// request_user_id := c.Locals("user_id").(uuid.UUID)

		/* Now update the flight log */
		flight_log.ID, err = db.UpdateFlightLog(txid, flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		response := fiber.Map{
			"txid":          txid.String(),
			"flight_log_id": flight_log.ID.String(),
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}
