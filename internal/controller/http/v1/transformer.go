package v1

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"userDataTransformer/internal/config"
	"userDataTransformer/internal/models"

	"github.com/gin-gonic/gin"
	"userDataTransformer/internal/domain"
	"userDataTransformer/internal/middleware"
	"userDataTransformer/internal/usecase"
)

type ProviderHandler struct {
	config  *config.Config
	service usecase.IProviderUsecase
	logger  *zap.Logger
}

func NewProviderHandler(ctx context.Context, config *config.Config, router *gin.RouterGroup, service usecase.IProviderUsecase, logger *zap.Logger, middleware middleware.IMiddlewareService) {
	handler := ProviderHandler{
		config:  config,
		service: service,
		logger:  logger,
	}
	providerGroup := router.Group("/v1/provider")
	{
		//после url был бы middleware для проверки авторизации
		providerGroup.POST("/import-users", handler.ImportUsers)
	}
	go service.RunWorkerPool(ctx)
}

// ImportUsers - обрабатывает POST-запрос с XML-документом в теле и ключом авторизации в хедерах.
// Для упрощения проверки сложная авторизация убрана, можно посмотреть пакет middleware
func (h *ProviderHandler) ImportUsers(ctx *gin.Context) {
	h.logger.Info("Начало обработки запроса ImportUsers")
	authKey := ctx.GetHeader("Authorization")
	if authKey == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": domain.ErrNoAuthToken.Error()})
		return
	}
	h.logger.Debug("Получен заголовок авторизации")
	var users models.XMLUsers
	err := ctx.ShouldBindXML(&users)
	if err != nil {
		h.logger.Error("Ошибка привязки XML",
			zap.Error(err),
			zap.String("content_type", ctx.ContentType()),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidRequest.Error()})
		return
	}
	if err = h.service.ProcessImport(ctx, &users); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Обработка импорта запущена")
	ctx.JSON(http.StatusAccepted, gin.H{"status": "processing started"})
}
