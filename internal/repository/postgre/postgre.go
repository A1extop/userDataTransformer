package postgre

import (
	"context"
	"userDataTransformer/internal/models"
	interfaces1 "userDataTransformer/internal/repository"
)

// Если вдруг когда-нибудь понадобится в сервисе этом что-то хранить. Пока что заглушка

type StubParserRepository struct{}

func NewStubParserRepository() interfaces1.IParserRepository {
	return &StubParserRepository{}
}

func (s *StubParserRepository) LogRequest(ctx context.Context, user models.JSONUser) error {
	return nil
}
