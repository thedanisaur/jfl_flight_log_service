package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/thedanisaur/jfl_platform/types"
	"github.com/thedanisaur/jfl_platform/util"
)

func CreateAircrew(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(CreateAircrew))
		response := fiber.Map{
			"txid": txid.String(),
		}
		return c.Status(fiber.StatusMethodNotAllowed).JSON(response)
	}
}
