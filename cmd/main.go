package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"userDataTransformer/internal/config"
	v1 "userDataTransformer/internal/controller/http/v1"
	localstore "userDataTransformer/internal/localStore/data_slice"
	"userDataTransformer/internal/middleware"
	repos1 "userDataTransformer/internal/repository/postgre"
	"userDataTransformer/internal/sender"
	usecase1 "userDataTransformer/internal/usecase"
	log "userDataTransformer/pkg/logger"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.New()
	logger, err := log.SetupLogger(cfg)
	if err != nil {
		return
	}
	gin.SetMode(cfg.App.Mode)

	router := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}

	router.Use(cors.New(corsConfig))

	repos1 := repos1.NewStubParserRepository() // затычка под возможность использовать бд
	sender := sender.NewRemoteSender(cfg.Sender.Endpoint)
	store := localstore.NewMemoryStorage()
	usecase1 := usecase1.NewProviderUsecase(repos1, logger, sender, store)
	mware := middleware.NewMiddlewareService(cfg)
	api := router.Group("/api")
	{
		v1.NewProviderHandler(ctx, cfg, api, usecase1, logger, mware)
	}
	Run(ctx, cfg, logger, router)
}
func Run(ctx context.Context, config *config.Config, logger *zap.Logger, router *gin.Engine) {
	srv := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
		Addr:         config.App.Host + ":" + config.App.Port,
		Handler:      router,
	}

	logger.Info("listen: " + config.App.Host + ":" + config.App.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen: " + err.Error())
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down graceful")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown: " + err.Error())
	}

	logger.Info("Server exiting")
}
