# Wallet Service

## Описание

userDataTransformer Service — это HTTP API, который преобразует из XML в JSON



## Архитектура

### Основные компоненты:
- **HTTP API**: реализован на Gin (`/api/v1/provider`, POST)
- **Хранилище данных**: 


### Статусы операций:
`success`: операция выполнена успешно  
`insufficient_funds`: недостаточно средств для списания  
`wallet_not_found`: кошелек не существует

## Запуск

### 1. Переменные окружения
Создайте файл `config.env` в корне проекта:

```env
HTTP_HOST=0.0.0.0
HTTP_PORT=8080
APP_MODE=debug
LOG_LEVEL=info

PG_HOST=localhost
PG_PORT=5432
PG_USERNAME=user
PG_PASSWORD=password
PG_DATABASE=wallet_db
PG_SSLMODE=disable
PG_MIGRATE=up
```

### 2. Запуск сервиса
- запуск сервиса происходит с помощью Docker compose
- docker compose build app
- docker compose up -d



## API
### Все эндпоинты начинаются с api/v1/wallet/.
- POST / - изменить баланс 
Request:
````
{
  "walletId": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "operationType": "DEPOSIT",
  "amount": 1000
}
````
Response:
````
{
  "result": "updated", 
  "balance": currentSum}
}
````

- GET /:id — Получить баланс кошелька по ID
Response:
````
{
  "uuid": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "balance": 100
}
````
- POST /create - Создать кошелек
  Response:
````
{
  "result": "created", 
  "uuid": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
}
````

## Упрощение
- Некоторые моменты упрощены, так как проект тестовый