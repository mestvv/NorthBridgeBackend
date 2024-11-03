package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiHttp "github.com/mestvv/NorthBridgeBackend/internal/api/http"
	"github.com/mestvv/NorthBridgeBackend/internal/config"
	"github.com/mestvv/NorthBridgeBackend/internal/db"
	"github.com/mestvv/NorthBridgeBackend/internal/repository"
	"github.com/mestvv/NorthBridgeBackend/internal/server"
	"github.com/mestvv/NorthBridgeBackend/internal/service"
	"github.com/mestvv/NorthBridgeBackend/pkg/auth"
	"github.com/mestvv/NorthBridgeBackend/pkg/email/smtp"
	"github.com/mestvv/NorthBridgeBackend/pkg/hash"
	log "github.com/mestvv/NorthBridgeBackend/pkg/logger"
	"github.com/mestvv/NorthBridgeBackend/pkg/otp"
)

const configPath = "config/config.yaml"

func main() {
	// Init cfg
	cfg := config.MustLoad(configPath)

	// Dependencies
	logger := log.SetupLogger(cfg.Env)

	logger.Info("starting backend api", "env", cfg.Env)
	logger.Debug("debug messages are enabled")

	// Init database
	dbMySQL, err := db.New(cfg.Database)
	if err != nil {
		logger.Error("mysql connect problem", "error", err)
		os.Exit(1)
	}
	defer func() {
		err = dbMySQL.Close()
		if err != nil {
			logger.Error("error when closing", "error", err)
		}
	}()
	logger.Info("mysql connection done")

	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	emailSender, err := smtp.NewSMTPSender(cfg.SMTP.From, cfg.SMTP.Pass, cfg.SMTP.Host, cfg.SMTP.Port)
	if err != nil {
		logger.Error("smtp sender creation failed", err)
		return
	}

	tokenManager, err := auth.NewManager(cfg.Auth.JWT)
	if err != nil {
		logger.Error("auth manager creation err", err)
		return
	}

	otpGenerator := otp.NewGOTPGenerator()

	// Services, Repos & API Handlers
	repos := repository.NewRepositories(dbMySQL)
	services := service.NewServices(service.Deps{
		Logger:       logger,
		Config:       cfg,
		Hasher:       hasher,
		TokenManager: tokenManager,
		OtpGenerator: otpGenerator,
		EmailSender:  emailSender,
		Repos:        repos,
	})
	handlers := apiHttp.NewHandlers(services, logger, tokenManager)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg))
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error occurred while running http server", "error", err)
		}
	}()
	logger.Info("server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Error("failed to stop server", "error", err)
	}

	logger.Info("app stopped")
}
