package main

import (
	"fmt"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/router"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("Start webhook application")
	engine := setupServer(logger)

	port := utils.GetEnvPortOrDefault("8081")
	if err := engine.Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.Fatal("Error starting webhook application", zap.Error(err))
	}
}

func setupServer(logger *zap.Logger) *gin.Engine {
	engine := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"*"},
	}))

	go router.Init(engine, logger)

	return engine
}
