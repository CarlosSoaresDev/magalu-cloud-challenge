package stripe

import (
	"io"
	"net/http"
	"os"

	stripeService "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/services/stripe"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/services/stripe/processor"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/webhook"
	"go.uber.org/zap"
)

type GatewayHandler struct {
	logger        *zap.Logger
	stripeService stripeService.StripeService
}

// New creates a new instance of GatewayHandler with the provided logger and stripeService.
// Parameters:
//   - logger: an instance of zap.Logger used for logging.
//   - stripeService: an instance of stripeService.StripeService used to interact with Stripe services.
//
// Returns:
//   - A pointer to a GatewayHandler instance.
func New(logger *zap.Logger, stripeService stripeService.StripeService) *GatewayHandler {
	return &GatewayHandler{
		logger:        logger,
		stripeService: stripeService,
	}
}

// WebhookHandler handles incoming Stripe webhook events.
// It reads the request body, verifies the webhook signature, and processes the event.
//
// The function performs the following steps:
// 1. Reads the Stripe webhook secret from environment variables.
// 2. Retrieves the Stripe-Signature header from the request.
// 3. Limits the request body size to 65536 bytes.
// 4. Reads the request body.
// 5. Verifies the webhook signature using the Stripe webhook secret.
// 6. Logs the received event ID and type.
// 7. Determines the processor type based on the event type.
// 8. Creates a new processor for the event type.
// 9. Processes the event using the appropriate processor.
// 10. Logs the success or failure of the event processing.
// 11. Sends an appropriate HTTP response based on the processing result.
//
// Parameters:
// - ctx: The Gin context for the request.
//
// Responses:
// - 400 Bad Request: If there is an error reading the request body, verifying the webhook signature, or if the payment gateway type is unsupported.
// - 500 Internal Server Error: If there is an error processing the payment.
// - 200 OK: If the event is successfully processed.
func (c *GatewayHandler) WebhookHandler(ctx *gin.Context) {

	req := ctx.Request

	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_KEY")
	signatureHeader := req.Header.Get("Stripe-Signature")

	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(ctx.Writer, req.Body, MaxBodyBytes)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		c.logger.Error("Error reading request body", zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	event, err := webhook.ConstructEvent(body, signatureHeader, stripeWebhookSecret)
	if err != nil {
		c.logger.Error("Error verifying webhook signature", zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	c.logger.Info("Webhook event received", zap.String("event_id", event.ID), zap.String("event_type", event.Type))
	processorType := processor.StripeProcessType(event.Type)
	res, err := processor.NewProcessor(processorType)

	if err != nil {
		c.logger.Error("Unsupported payment gateway type", zap.String("event_id", event.ID), zap.String("event_type", event.Type), zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = res.Process(c.stripeService, event)
	if err != nil {
		c.logger.Error("Error processing payment", zap.String("event_id", event.ID), zap.String("event_type", event.Type), zap.Error(err))
		utils.ApiResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	c.logger.Info("Successfully processed stripe request", zap.String("event_id", event.ID), zap.String("event_type", event.Type))
	utils.ApiResponse(ctx, http.StatusOK, nil)
}
