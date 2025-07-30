// локальное хранилище данных
// сделано именно так, потому что по тз нужно отправлять в другой сервис обязательно через Rest Api. Вообще использовал бы лучше шину данных
// вполне можно для возможность сохранения взять редис, а так это упрощённый вариант
package data_slice

import (
	"context"
	"sync"
	"userDataTransformer/internal/localStore"
	"userDataTransformer/internal/models"
)

type MemoryStorage struct {
	mu    *sync.Mutex
	items []models.RetryItem
}

func NewMemoryStorage() localStore.ILocalStore {
	return &MemoryStorage{
		mu:    &sync.Mutex{},
		items: []models.RetryItem{},
	}
}
func (m *MemoryStorage) Add(ctx context.Context, items []models.RetryItem) error { // Принимаем срез, а не указатель
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = append(m.items, items...) // Простое добавление элементов

	return nil
}

func (s *MemoryStorage) GetAll(ctx context.Context) ([]models.RetryItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.items) == 0 {
		return []models.RetryItem{}, nil
	}

	items := s.items
	s.items = nil
	return items, nil
}
