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

func CreateFlightlog(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", util.GetFunctionName(CreateFlightlog), txid.String())

		var flight_log types.FlightLogDTO
		err := c.BodyParser(&flight_log)
		if err != nil {
			log.Printf("Failed to parse flight log data\n%s\n", err.Error())
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Failed to parse flight log data: %s\n", txid.String()))
		}
		/* Get the requesting user's id */
		request_user_id := c.Locals("user_id").(uuid.UUID)

		/* Now start inserting the flight log */
		flight_log.ID, err = db.InsertFlightLog(txid, request_user_id, flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		mission_ids, err := db.InsertMissions(txid, flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		aircrew_ids, err := db.InsertAircrew(txid, flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		comment_ids, err := db.InsertFlightLogComments(txid, request_user_id, flight_log)
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

func GetAllFlightlogs(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", util.GetFunctionName(GetAllFlightlogs), txid.String())
		response := fiber.Map{
			"txid": txid.String(),
		}
		return c.Status(fiber.StatusMethodNotAllowed).JSON(response)
	}
}

func GetFlightlog(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", util.GetFunctionName(GetFlightlog), txid.String())

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
		log.Printf("%s | %s\n", util.GetFunctionName(GetFlightlogs), txid.String())
		response := fiber.Map{
			"txid": txid.String(),
		}
		return c.Status(fiber.StatusMethodNotAllowed).JSON(response)
	}
}

func UpdateFlightlog(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", util.GetFunctionName(UpdateFlightlog), txid.String())
		response := fiber.Map{
			"txid": txid.String(),
		}
		return c.Status(fiber.StatusMethodNotAllowed).JSON(response)
	}
}
