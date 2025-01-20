package currency

import (
	"net/http"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/models"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/services/currency"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CurrencyHandler struct {
	logger          *zap.Logger
	currencyService currency.CurrencyService
}

// New creates a new instance of CurrencyHandler with the provided logger and currency service.
// Parameters:
//   - logger: an instance of zap.Logger used for logging within the handler.
//   - currencyService: an instance of currency.CurrencyService that provides currency-related operations.
//
// Returns:
//   - A pointer to a newly created CurrencyHandler.
func New(logger *zap.Logger, currencyService currency.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{
		logger:          logger,
		currencyService: currencyService,
	}
}

// GetAllCurrencyHandler handles the request to retrieve all currencies.
// It extracts the correlation ID from the context, logs the request initiation,
// calls the currency service to get all currencies, and returns the result in the response.
// If any error occurs during the process, it logs the error and returns an appropriate
// HTTP response with the error message.
//
// @Summary Retrieve all currencies
// @Description Get a list of all available currencies
// @Tags currency
// @Accept json
// @Produce json
// @Success 200 {object} []Currency
// @Failure 400 {object} ErrorResponse
// @Router /currencies [get]
func (c *CurrencyHandler) GetAllCurrencyHandler(ctx *gin.Context) {

	correlationId, err := utils.GetCorrelationId(ctx)

	if err != nil {
		c.logger.Error("Failed to get correlation ID", zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	c.logger.Info("Initialize request to get all currency", zap.String("correlation_id", correlationId))

	result, err := c.currencyService.GetAllCurrency()

	if err != nil {
		c.logger.Error("Failed to get all currency", zap.String("CorrelationId", correlationId), zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, "We were unable to process your request, please try later")
		return
	}

	utils.ApiResponse(ctx, http.StatusOK, result)
	c.logger.Info("Successfully retrieved all currency", zap.String("CorrelationId", correlationId))
}

// ConvertExchangeRateHandler handles the request to convert currency exchange rates.
// It retrieves the correlation ID from the context, binds the JSON payload to the CurrencyConvert model,
// and calls the currencyService to perform the conversion. If any error occurs during these steps,
// it logs the error and sends an appropriate HTTP response. On success, it returns the conversion result
// and logs the successful completion of the request.
//
// @Summary Convert currency exchange rate
// @Description Converts the currency exchange rate based on the provided payload
// @Tags currency
// @Accept json
// @Produce json
// @Param payload body models.CurrencyConvert true "Currency conversion payload"
// @Success 200 {object} models.CurrencyConvertResponse
// @Failure 400 {object} utils.ApiErrorResponse
// @Router /currency/convert [post]
func (c *CurrencyHandler) ConvertExchangeRateHandler(ctx *gin.Context) {

	correlationId, err := utils.GetCorrelationId(ctx)
	if err != nil {
		c.logger.Error("Failed to get correlation ID", zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var payload models.CurrencyConvert
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		c.logger.Error("Failed to bind JSON payload", zap.String("correlation_id", correlationId), zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, utils.ValidatorError(err))
		return
	}

	c.logger.Info("Starting currency conversion request", zap.String("correlation_id", correlationId))

	res, err := c.currencyService.ConvertExchangeRate(payload)

	if err != nil {
		c.logger.Error("Currency conversion failed", zap.String("correlation_id", correlationId), zap.Error(err))
		utils.ApiResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.ApiResponse(ctx, http.StatusOK, res)
	c.logger.Info("Currency conversion completed successfully", zap.String("correlation_id", correlationId))
}
