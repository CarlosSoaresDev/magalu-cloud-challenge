package router

import (
	"net/http"

	currencyHandler "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/handlers/currency"
	gatewayHandler "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/handlers/gateway"
	currencyService "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/services/currency"
	gatewayService "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/services/gateway"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/cache"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init(route *gin.Engine, logger *zap.Logger) {
	cacheClient := cache.New()

	currencyService := currencyService.New(cacheClient)
	currencyHandler := currencyHandler.New(logger, currencyService)

	gatewayService := gatewayService.New(cacheClient)
	gatewayHandler := gatewayHandler.New(logger, gatewayService)

	groupRoute := route.Group("/api/v1")

	currencyRoute := groupRoute.Group("/currencies")
	{
		currencyRoute.GET("", currencyHandler.GetAllCurrencyHandler)
		currencyRoute.POST("convert", currencyHandler.ConvertExchangeRateHandler)
	}

	gatewayRoute := groupRoute.Group("/gateways")
	{
		gatewayRoute.GET("avaiables", gatewayHandler.GetAllAvaiablesGateways)
		gatewayRoute.GET("transactions", gatewayHandler.GetAllTransactionsByDateHandler)
		gatewayRoute.POST("", gatewayHandler.PaymentHandler)
	}

	route.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})
}
