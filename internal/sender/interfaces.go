package sender

import (
	"context"
	"userDataTransformer/internal/models"
)

type IRemoteSender interface {
	SendUser(ctx context.Context, user models.JSONUser) error
}
