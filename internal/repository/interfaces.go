package repository

import (
	"context"

	"userDataTransformer/internal/models"
)

// Если вдруг когда-нибудь понадобится в сервисе этом что-то хранить. Пока что заглушка
type IParserRepository interface {
	LogRequest(ctx context.Context, user models.JSONUser) error
}
