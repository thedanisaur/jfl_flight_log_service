package handlers

import (
	"flight_log_service/db"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/thedanisaur/jfl_platform/types"
	"github.com/thedanisaur/jfl_platform/util"
)

func CreateFlightlogComment(config types.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		txid := c.Locals("transaction_id").(uuid.UUID)
		log.Printf("%s | %s\n", txid.String(), util.GetFunctionName(CreateFlightlogComment))

		var flight_log_comment types.FlightLogCommentDTO
		err := c.BodyParser(&flight_log_comment)
		if err != nil {
			log.Printf("Failed to parse flight log data\n%s\n", err.Error())
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Failed to parse flight log data: %s\n", txid.String()))
		}
		/* Get the target user/flight log info */
		// TODO [drd] should we make sure we're commenting on a flight log owned by the target user or is any uuid a valid address?
		flight_log_id, err := uuid.Parse(c.Params("flight_log_id"))
		if err != nil {
			log.Printf("Failed to parse flight log id: %s\n", c.Params("flight_log_id"))
			return c.Status(fiber.StatusServiceUnavailable).SendString("invalid flight log")
		}
		if flight_log_id != flight_log_comment.FlightLogID {
			log.Printf("Attempted write comment to wrong flight log. URL ID: %s. Flight Log Comment - Flight Log ID: %s\n", flight_log_id, flight_log_comment.FlightLogID)
			return c.Status(fiber.StatusBadRequest).SendString("malformed flight log comment")
		}
		/* Get the requesting user */
		request_user := c.Locals("user_claims").(types.UserClaims)

		/* Now start inserting the flight log comment */
		comment_id, err := db.InsertFlightLogComment(txid, request_user.UserID, flight_log_comment)
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString(err.Error())
		}
		response := fiber.Map{
			"txid":       txid.String(),
			"comment_id": comment_id,
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}
