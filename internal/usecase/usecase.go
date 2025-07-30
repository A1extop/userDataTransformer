package usecase

import (
	"context"
	"go.uber.org/zap"
	"time"
	"userDataTransformer/internal/domain"

	//"userDataTransformer/internal/domain"
	"userDataTransformer/internal/localStore"
	"userDataTransformer/internal/models"
	"userDataTransformer/internal/repository"
	"userDataTransformer/internal/sender"
)

const MaxConcurrent = 100

type ProviderUsecase struct {
	repo       repository.IParserRepository
	logger     *zap.Logger // если используешь интерфейс логгера
	sender     sender.IRemoteSender
	localStore localStore.ILocalStore
}

func NewProviderUsecase(repo repository.IParserRepository, logger *zap.Logger, sender sender.IRemoteSender, localStore localStore.ILocalStore) IProviderUsecase {

	return &ProviderUsecase{
		repo:       repo,
		logger:     logger,
		sender:     sender,
		localStore: localStore,
	}
}

// ProcessImport - преобразует данные из XML в JSON и отправляет в локальное хранилище
// (в целом, я здесь мог использовать горутины, чтобы распараллелить ещё дополнительно преобразование каждой, но счёл это избыточным действием в данном контексте)
func (u *ProviderUsecase) ProcessImport(ctx context.Context, xmlUsers *models.XMLUsers) error {
	u.logger.Info("Начало обработки импорта", zap.Int("users_count", len(xmlUsers.Users)))
	if len(xmlUsers.Users) == 0 {
		u.logger.Error("Ошибка обработки импорта", zap.Error(domain.ErrNoData))
		return domain.ErrNoData
	}

	var items []models.RetryItem
	for _, user := range xmlUsers.Users {
		jsonUser := user.ToJSONUser()
		if user.Age <= 0 {
			continue
		}
		items = append(items, models.RetryItem{
			Data:     jsonUser,
			Attempts: 0,
		})
	}
	u.logger.Info("Фильтрация пользователей завершена")

	if err := u.localStore.Add(ctx, items); err != nil {
		u.logger.Error("Ошибка добавления в локальное хранилище",
			zap.Error(err),
			zap.Int("items_count", len(items)),
		)
		return domain.ErrFailedAdd
	}
	u.logger.Info("Данные успешно добавлены в локальное хранилище")
	return nil
}

// RunWorkerPool - раз в определенное время забирает данные из локального хранилища обрабатывает.
// Использован паттерн воркер пул для ограничения количества горутин
// В одной части мы асинхронно забираем данные из локалки и отправляет асинхронно. Воркер пул после асинхронно через канал принимает данные и обрабатывает
func (u *ProviderUsecase) RunWorkerPool(ctx context.Context) {
	u.logger.Info("Запуск воркер пула", zap.Int("workers", MaxConcurrent))
	taskCh := make(chan models.RetryItem, 100)
	retryCh := make(chan models.RetryItem, 100)

	go func() {
		u.logger.Debug("Запуск обработчика повторов")
		var buffer []models.RetryItem
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				u.logger.Debug("Остановка обработчика повторов (контекст завершен)")
				return
			case item := <-retryCh:
				if item.Attempts < 5 {
					buffer = append(buffer, item)
				}
			case <-ticker.C:
				if len(buffer) > 0 {
					if err := u.localStore.Add(ctx, buffer); err != nil {
						u.logger.Error("Ошибка сохранения элементов повтора",
							zap.Error(err),
							zap.Int("count", len(buffer)),
						)
					} else {
						buffer = buffer[:0]
					}

				}
			}
		}
	}()

	go func() {
		u.logger.Debug("Запуск поставщика задач")
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				u.logger.Debug("Остановка поставщика задач (контекст завершен)")
				close(taskCh)
				return
			case <-ticker.C:
				items, err := u.localStore.GetAll(ctx)
				if err != nil {
					u.logger.Error("Ошибка получения задач из хранилища", zap.Error(err))
					continue
				}
				if items == nil {
					continue
				}
				for _, item := range items {
					select {
					case taskCh <- item:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	for i := 0; i < MaxConcurrent; i++ {
		go func(workerID int) {
			u.logger.Debug("Запуск воркера", zap.Int("worker_id", workerID))
			for {
				select {
				case <-ctx.Done():
					u.logger.Debug("Остановка воркера", zap.Int("worker_id", workerID))
					return
				case item, ok := <-taskCh:

					if !ok {
						u.logger.Debug("Канал задач закрыт", zap.Int("worker_id", workerID))
						return
					}
					u.logger.Debug("Обработка задачи",
						zap.Int("worker_id", workerID),
						zap.Any("user_id", item.Data.ID),
						zap.Int("attempt", item.Attempts+1),
					)
					err := u.sender.SendUser(ctx, *item.Data)
					if err != nil {
						u.logger.Error("Ошибка отправки пользователя",
							zap.Int("worker_id", workerID),
							zap.Error(err),
							zap.Any("user_id", item.Data.ID),
							zap.Int("attempt", item.Attempts+1),
						)
						item.Attempts++
						retryCh <- item
					}
				}
			}
		}(i)
	}
	u.logger.Info("Воркер пул запущен")
	<-ctx.Done()
	close(taskCh)
	close(retryCh)
	u.logger.Info("Воркер пул остановлен")
}
