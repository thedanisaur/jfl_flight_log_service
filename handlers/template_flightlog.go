package handlers

import (
	"fmt"
	"log"

	"flight_log_service/db"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/thedanisaur/jfl_platform/types"
	"github.com/thedanisaur/jfl_platform/util"
)

func CreateTemplateFlightlog(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(CreateTemplateFlightlog))

		var flight_log types.FlightLogDTO
		err := c.BodyParser(&flight_log)
		if err != nil {
			log.Printf("Failed to parse template flight log data\n%s\n", err.Error())
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Failed to parse template flight log data: %s\n", txid.String()))
		}
		/* Get the requesting user */
		request_user := c.Locals("user_claims").(types.UserClaims)

		/* Now start inserting the flight log */
		flight_log.ID, err = db.InsertTemplateFlightLog(txid, request_user.UserID, flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		mission_ids, err := db.InsertTemplateMissions(txid, flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		aircrew_ids, err := db.InsertTemplateAircrews(txid, flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		response := fiber.Map{
			"txid":          txid.String(),
			"flight_log_id": flight_log.ID.String(),
			"mission_ids":   mission_ids,
			"aircrew_ids":   aircrew_ids,
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func GetTemplateFlightlog(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetTemplateFlightlog))

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

func GetTemplateFlightlogs(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(GetTemplateFlightlogs))

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

func UpdateTemplateFlightlog(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(UpdateTemplateFlightlog))

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
