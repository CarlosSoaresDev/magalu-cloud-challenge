package router

import (
	"net/http"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/cache"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	paypalHandler "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/handlers/paypal"
	stripeHandler "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/handlers/stripe"

	paypalService "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/services/paypal"
	stripeService "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/services/stripe"
)

func Init(route *gin.Engine, logger *zap.Logger) {
	cacheClient := cache.New()

	stripeService := stripeService.New(cacheClient)
	stripeHandler := stripeHandler.New(logger, stripeService)

	paypalService := paypalService.New(cacheClient)
	paypalHandler := paypalHandler.New(logger, paypalService)

	groupRoute := route.Group("/api/v1")

	stripeGroup := groupRoute.Group("/stripe")
	{
		stripeGroup.POST("/webhook", stripeHandler.WebhookHandler)
	}

	gatewayGroup := groupRoute.Group("/paypal")
	{
		gatewayGroup.POST("/webhook", paypalHandler.WebhookHandler)
	}

	route.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})
}
