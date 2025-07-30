package usecase

import (
	"context"
	"userDataTransformer/internal/models"
)

type IProviderUsecase interface {
	ProcessImport(ctx context.Context, xmlUsers *models.XMLUsers) error
	RunWorkerPool(ctx context.Context)
}
