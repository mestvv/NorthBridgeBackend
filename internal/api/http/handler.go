package apiHttp

import (
	"log/slog"

	sloggin "github.com/samber/slog-gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/mestvv/NorthBridgeBackend/docs"
	"github.com/mestvv/NorthBridgeBackend/pkg/auth"
	"github.com/mestvv/NorthBridgeBackend/pkg/limiter"
	"github.com/mestvv/NorthBridgeBackend/pkg/validator"

	"github.com/gin-gonic/gin"
	internalV1 "github.com/mestvv/NorthBridgeBackend/internal/api/http/internal/v1"
	"github.com/mestvv/NorthBridgeBackend/internal/config"
	"github.com/mestvv/NorthBridgeBackend/internal/service"
)

type Handler struct {
	services     *service.Services
	logger       *slog.Logger
	tokenManager auth.TokenManager
}

func NewHandlers(services *service.Services, logger *slog.Logger, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		logger:       logger,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	validator.RegisterGinValidator()

	router.Use(
		gin.Recovery(),
		sloggin.NewWithConfig(h.logger, sloggin.Config{
			WithSpanID:  true,
			WithTraceID: true,
		}),
		limiter.Limit(cfg.Limiter.RPS, cfg.Limiter.Burst, cfg.Limiter.TTL, h.logger),
		corsMiddleware,
	)

	if cfg.HttpServer.SwaggerEnabled {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler(), ginSwagger.InstanceName("internal")))
	}

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	internalHandlersV1 := internalV1.NewHandler(h.services, h.logger, h.tokenManager)
	api := router.Group("/api")
	{
		internalHandlersV1.Init(api)
	}
}
