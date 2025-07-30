package localStore

import (
	"context"
	"userDataTransformer/internal/models"
)

type ILocalStore interface {
	Add(ctx context.Context, items []models.RetryItem) error
	GetAll(ctx context.Context) ([]models.RetryItem, error)
}
