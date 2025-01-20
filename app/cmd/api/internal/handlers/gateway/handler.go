package gateway

import (
	"net/http"
	"time"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/models"
	gatewayService "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/services/gateway"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/services/gateway/provider"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GatewayHandler struct {
	logger         *zap.Logger
	gatewayService gatewayService.GatewayService
}

// New creates a new instance of GatewayHandler with the provided logger and gateway service.
// Parameters:
//   - logger: an instance of zap.Logger for logging purposes.
//   - gatewayService: an instance of GatewayService to handle gateway operations.
//
// Returns:
//   - A pointer to a newly created GatewayHandler.
func New(logger *zap.Logger, gatewayService gatewayService.GatewayService) *GatewayHandler {
	return &GatewayHandler{
		logger:         logger,
		gatewayService: gatewayService,
	}
}

// GetAllAvaiablesGateways handles the request to retrieve all available gateways.
// It extracts the correlation ID from the context, logs the start of the request,
// calls the gateway service to get all available gateways, and returns the result
// in the API response. If there is an error in extracting the correlation ID, it
// responds with a bad request status and logs the error.
//
// @Summary Retrieve all available gateways
// @Description Retrieves a list of all available gateways from the gateway service
// @Tags gateways
// @Accept json
// @Produce json
// @Success 200 {object} []Gateway "List of available gateways"
// @Failure 400 {object} string "Bad request"
// @Router /gateways [get]
func (c *GatewayHandler) GetAllAvaiablesGateways(ctx *gin.Context) {

	correlationId, err := utils.GetCorrelationId(ctx)
	if err != nil {
		c.logger.Error("Failed to get correlation ID", zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	c.logger.Info("Starting request to get all gateways", zap.String("correlation_id", correlationId))

	result := c.gatewayService.GetAllAvaiablesGateways()

	utils.ApiResponse(ctx, http.StatusOK, result)
	c.logger.Info("Successfully retrieved all gateways", zap.String("correlation_id", correlationId))
}

// GetAllTransactionsByDateHandler handles the request to retrieve all transactions by a specific date.
// It expects a query parameter "date" in the format "dd_mm_yyyy". If the date is not provided, it defaults to the current date.
// The function retrieves the correlation ID from the context for logging purposes.
// It logs the start of the request, calls the gateway service to get all transactions by the specified date,
// and returns the result in the response. If any error occurs during the process, it logs the error and returns an appropriate HTTP status code and message.
//
// @Summary Get all transactions by date
// @Description Retrieve all transactions for a given date. If no date is provided, the current date is used.
// @Tags transactions
// @Accept json
// @Produce json
// @Param date query string false "Date in format dd_mm_yyyy"
// @Success 200 {object} []Transaction "List of transactions"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /transactions [get]
func (c *GatewayHandler) GetAllTransactionsByDateHandler(ctx *gin.Context) {

	correlationId, err := utils.GetCorrelationId(ctx)

	if err != nil {
		c.logger.Error("Failed to get correlation ID", zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	date := ctx.Query("date")

	if utils.IsEmptyOrNull(date) {
		date = time.Now().Format("02_01_2006")
	}

	c.logger.Info("Starting request to get all gateways", zap.String("correlation_id", correlationId))

	result, err := c.gatewayService.GetAllTransactionsByDate(date)
	if err != nil {
		c.logger.Error("Failed to get all gateways", zap.String("correlation_id", correlationId), zap.Error(err))
		utils.ApiResponse(ctx, http.StatusInternalServerError, "Unable to process your request, please try again later")
		return
	}

	utils.ApiResponse(ctx, http.StatusOK, result)
	c.logger.Info("Successfully retrieved all gateways", zap.String("correlation_id", correlationId), zap.Int("GatewayCount", len(*result)))
}

// PaymentHandler handles payment requests by processing the payment through the specified gateway provider.
// It retrieves the correlation ID from the context, binds the JSON payload to the Gateway model, and logs the start of the payment request.
// The handler then initializes the appropriate payment provider based on the payload's gateway type and processes the payment.
// If any errors occur during these steps, appropriate error responses are returned to the client.
// Upon successful payment processing, the transaction is added to the gateway service, and a no-content response is returned.
//
// @Summary Process payment request
// @Description Processes a payment request through the specified gateway provider
// @Tags Payment
// @Accept json
// @Produce json
// @Param payload body models.Gateway true "Payment payload"
// @Success 204 "No Content"
// @Failure 400 {object} utils.ApiError "Bad Request"
// @Router /payment [post]
func (c *GatewayHandler) PaymentHandler(ctx *gin.Context) {

	correlationId, err := utils.GetCorrelationId(ctx)
	if err != nil {
		c.logger.Error("Failed to get correlation ID", zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var payload models.Gateway
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		c.logger.Error("Failed to bind JSON payload", zap.String("correlation_id", correlationId), zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, utils.ValidatorError(err))
		return
	}

	c.logger.Info("Starting payment request", zap.String("correlation_id", correlationId))

	providerType := provider.ProviderType(payload.Gateway)
	provider, err := provider.NewProvider(providerType)
	if err != nil {
		c.logger.Error("Unsupported payment gateway type", zap.String("correlation_id", correlationId), zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := provider.ProcessPayment(payload, correlationId)

	if err != nil {
		c.logger.Error("Payment processing failed", zap.String("correlation_id", correlationId), zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = c.gatewayService.AddTransaction(*res, payload)

	if err != nil {
		c.logger.Error("Payment processing failed", zap.String("correlation_id", correlationId), zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, "Payment processing failed")
		return
	}

	utils.ApiResponse(ctx, http.StatusNoContent, nil)
	c.logger.Info("Payment request completed successfully", zap.String("correlation_id", correlationId))
}
