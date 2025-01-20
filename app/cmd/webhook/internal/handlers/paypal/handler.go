package paypal

import (
	"net/http"
	"os"

	paypalService "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/services/paypal"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/services/paypal/processor"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PayPalHandler struct {
	logger        *zap.Logger
	paypalService paypalService.PaypalService
}

// New creates a new instance of PayPalHandler with the provided logger and PayPal service.
// It returns a pointer to the newly created PayPalHandler.
//
// Parameters:
//   - logger: A zap.Logger instance used for logging.
//   - paypalService: An instance of paypalService.PaypalService to handle PayPal-related operations.
//
// Returns:
//   - A pointer to a PayPalHandler instance.
func New(logger *zap.Logger, paypalService paypalService.PaypalService) *PayPalHandler {
	return &PayPalHandler{
		logger:        logger,
		paypalService: paypalService,
	}
}

// WebhookHandler handles incoming webhook requests from PayPal.
// It reads the request body, verifies the webhook signature, and processes the event.
// If the event type is supported, it processes the payment using the PayPal service.
// It responds with appropriate HTTP status codes based on the success or failure of each step.
//
// @param ctx *gin.Context - the context for the request
//
// Environment Variables:
// - PAYPAL_WEBHOOK_KEY: the secret key used to verify the PayPal webhook signature
//
// Possible HTTP Status Codes:
// - 400 Bad Request: if there is an error reading the request body or verifying the webhook signature
// - 500 Internal Server Error: if there is an error processing the payment
// - 200 OK: if the webhook event is successfully processed
func (c *PayPalHandler) WebhookHandler(ctx *gin.Context) {

	paypalWebhookSecret := os.Getenv("PAYPAL_WEBHOOK_KEY")

	event := utils.GenerateGUID()

	c.logger.Info("Webhook event received", zap.String("event_id", event), zap.String("event_type", event), zap.String("webhook_secret", paypalWebhookSecret))

	processorType := processor.PayPalProcessType(event)
	res, err := processor.NewProcessor(processorType)

	if err != nil {
		c.logger.Error("Unsupported payment gateway type", zap.String("event_id", event), zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = res.Process(c.paypalService, event)
	if err != nil {
		c.logger.Error("Error processing payment", zap.String("event_id", event), zap.Error(err))
		utils.ApiResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	c.logger.Info("Successfully processed PayPal request", zap.String("event_id", event))
	utils.ApiResponse(ctx, http.StatusOK, nil)
}
