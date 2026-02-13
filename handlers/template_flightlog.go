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

		var template_flight_log types.TemplateFlightLogDTO
		err := c.BodyParser(&template_flight_log)
		if err != nil {
			log.Printf("Failed to parse template flight log data\n%s\n", err.Error())
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Failed to parse template flight log data: %s\n", txid.String()))
		}
		/* Get the requesting user */
		request_user := c.Locals("user_claims").(types.UserClaims)

		/* Now start inserting the flight log */
		template_flight_log.ID, err = db.InsertTemplateFlightLog(txid, request_user.UserID, template_flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		mission_ids, err := db.InsertTemplateMissions(txid, template_flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		aircrew_ids, err := db.InsertTemplateAircrews(txid, template_flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		response := fiber.Map{
			"txid":                   txid.String(),
			"template_flight_log_id": template_flight_log.ID.String(),
			"template_mission_ids":   mission_ids,
			"template_aircrew_ids":   aircrew_ids,
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func DeleteTemplateFlightlog(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(DeleteTemplateFlightlog))

		user_id, err := uuid.Parse(c.Params("user_id"))
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString("invalid user")
		}
		template_id, err := uuid.Parse(c.Params("template_id"))
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString("invalid template flight log")
		}

		flight_log, err := db.DeleteTemplateFlightlog(txid, user_id, template_id)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}

		// response := fiber.Map{
		// 	"txid": txid.String(),
		// }

		return c.Status(fiber.StatusOK).JSON(flight_log)
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
		template_id, err := uuid.Parse(c.Params("template_id"))
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString("invalid flight log")
		}

		flight_log, err := db.GetTemplateFlightlog(txid, user_id, template_id)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}

		flight_log.Missions, err = db.GetTemplateMissions(txid, template_id)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}

		flight_log.Aircrew, err = db.GetTemplateAirCrews(txid, template_id)
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

		template_flight_logs, err := db.GetTemplateFlightlogs(txid, user_id)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}

		// for index := range template_flight_logs {
		// 	template_flight_logs[index].Missions, err = db.GetTemplateMissions(txid, template_flight_logs[index].ID)
		// 	if err != nil {
		// 		return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		// 	}

		// 	template_flight_logs[index].Aircrew, err = db.GetTemplateAirCrews(txid, template_flight_logs[index].ID)
		// 	if err != nil {
		// 		return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		// 	}
		// }

		// response := fiber.Map{
		// 	"txid": txid.String(),
		// }

		return c.Status(fiber.StatusOK).JSON(template_flight_logs)
	}
}

func UpdateTemplateFlightlog(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(UpdateTemplateFlightlog))

		var template_flight_log types.TemplateFlightLogDTO
		err := c.BodyParser(&template_flight_log)
		if err != nil {
			log.Printf("Failed to parse template flight log data\n%s\n", err.Error())
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Failed to parse template flight log data: %s\n", txid.String()))
		}
		// TODO [drd] validate that this action is allowed.
		/* Get the requesting user */
		// request_user := c.Locals("user_claims").(types.UserClaims)

		/* Now update the flight log */
		template_id, err := db.UpdateTemplateFlightLog(txid, template_flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		mission_ids, err := db.UpdateTemplateMissions(txid, template_flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		aircrew_ids, err := db.UpdateTemplateAircrews(txid, template_flight_log)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		response := fiber.Map{
			"txid":                   txid.String(),
			"template_flight_log_id": template_id,
			"template_mission_ids":   mission_ids,
			"template_aircrew_ids":   aircrew_ids,
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}
