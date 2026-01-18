package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/thedanisaur/jfl_platform/types"
	"github.com/thedanisaur/jfl_platform/util"
)

func CreateFlightlogComment(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := uuid.New()
		log.Printf("%s | %s\n", util.GetFunctionName(CreateFlightlogComment), txid.String())
		response := fiber.Map{
			"txid": txid.String(),
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}
