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

func GetFlightlog(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", util.GetFunctionName(CreateFlightlog), txid.String())
		response := fiber.Map{
			"txid": txid.String(),
		}
		return c.Status(fiber.StatusMethodNotAllowed).JSON(response)
	}
}

func GetFlightlogs(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", util.GetFunctionName(CreateFlightlog), txid.String())
		response := fiber.Map{
			"txid": txid.String(),
		}
		return c.Status(fiber.StatusMethodNotAllowed).JSON(response)
	}
}

func UpdateFlightlog(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", util.GetFunctionName(CreateFlightlog), txid.String())
		response := fiber.Map{
			"txid": txid.String(),
		}
		return c.Status(fiber.StatusMethodNotAllowed).JSON(response)
	}
}
