# auth_service

[![Coverage Status](https://coveralls.io/repos/github/Sugyk/auth_service/badge.svg?branch=master)](https://coveralls.io/github/Sugyk/auth_service)
[![Go version](https://img.shields.io/badge/go-1.25-00ADD8?logo=go)](go.mod)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

HTTP-сервис аутентификации пользователей на Go. Реализует регистрацию, вход и выдачу JWT-токенов. Написан как pet-проект для изучения best practice в Go разработке: слоистая архитектура, миграции БД, контейнеризация, unit- и интеграционные тесты.

---

## Возможности

- Регистрация пользователя с хешированием пароля (bcrypt, cost = 12)
- Аутентификация и выдача подписанного JWT-токена
- Защита `/login` от брутфорса: блокировка логина в Redis после серии неудачных попыток
- HTTP и gRPC API поверх одной и той же бизнес-логики (транспорты не дублируют её)
- Настраиваемый TTL токена
- Управление пулом соединений с PostgreSQL (pgx v5)
- Автоматические миграции при старте через `migrate`
- Конфигурация через YAML-файл и переменные окружения

---

## Стек

| Слой | Технология |
|---|---|
| Язык | Go 1.25 |
| База данных | PostgreSQL 16 |
| Драйвер БД | [pgx/v5](https://github.com/jackc/pgx) |
| Троттлинг логина | Redis 7 + [go-redis/v9](https://github.com/redis/go-redis) |
| gRPC | [grpc-go](https://github.com/grpc/grpc-go) + Protocol Buffers |
| JWT | [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt) |
| Хеширование | bcrypt (`golang.org/x/crypto`) |
| Конфигурация | [Viper](https://github.com/spf13/viper) + godotenv |
| Миграции | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Моки | [go.uber.org/mock](https://github.com/uber-go/mock) |
| Тесты | testify |
| Контейнеры | Docker Compose |

---

## Архитектура

```
auth_service/
├── cmd/auth_service/     # точка входа, сборка зависимостей (DI)
├── internal/
│   ├── api/http/handler/ # HTTP-обработчики
│   ├── api/grpc/         # gRPC-сервер (тот же service-слой, что и у HTTP)
│   ├── service/          # бизнес-логика
│   ├── repository/       # работа с БД
│   ├── migrations/       # SQL-миграции
│   └── ...
├── pkg/                  # переиспользуемые пакеты (txmanager, redis, hasher, jwt)
├── proto/                # protobuf-контракт gRPC API
├── config/               # config.yaml
├── tests/                # интеграционные тесты
└── docker/               # Dockerfile
```

Сервис следует трёхслойной архитектуре `handler → service → repository`. Каждый слой работает через интерфейс — это позволяет подменять реализации в тестах через моки (go.uber.org/mock). HTTP- и gRPC-серверы — два независимых транспорта поверх одного и того же `service.Service`.

---

## Быстрый старт

### Требования

- Docker и Docker Compose

### Запуск одной командой

```bash
git clone https://github.com/Sugyk/auth_service.git
cd auth_service
make run
```

Команда соберёт образ, запустит PostgreSQL, применит миграции и поднимет сервис на `localhost:8080`. Порядок запуска контролируется через `healthcheck` и `depends_on` в `compose.yaml`.

### Конфигурация

Скопируйте `.env.example` в `.env` и задайте свои значения:

```bash
cp .env.example .env
```

| Переменная | Описание | Пример |
|---|---|---|
| `APP_PG_CONNSTR` | DSN строка подключения к PostgreSQL | `postgres://user:pass@localhost:5432/mydb?sslmode=disable` |
| `APP_PG_MAX_CONNS` | Максимум соединений в пуле | `25` |
| `APP_PG_MIN_CONNS` | Минимум соединений в пуле | `2` |
| `APP_PG_MAX_CONN_LIFETIME` | Максимальное время жизни соединения (сек) | `1800` |
| `APP_PG_MAX_CONN_IDLE_TIME` | Максимальное время простоя соединения (сек) | `300` |
| `APP_HASHER_COST` | Стоимость bcrypt (рекомендуется ≥ 12) | `12` |
| `JWT_SECRET` | Секретный ключ для подписи токенов | `your-secret-key` |
| `APP_JWT_TTL` | Время жизни JWT-токена | `24h` |
| `APP_GRPC_ADDR` | Адрес, на котором слушает gRPC-сервер | `:50051` |
| `APP_REDIS_ADDR` | Адрес Redis, используемого для троттлинга логина | `localhost:6379` |
| `APP_REDIS_PASSWORD` | Пароль Redis | (пусто) |
| `APP_REDIS_DB` | Номер базы Redis | `0` |
| `APP_THROTTLE_MAX_ATTEMPTS` | Сколько неудачных попыток логина допускается перед блокировкой | `10` |
| `APP_THROTTLE_BLOCK_DURATION` | На сколько блокируется логин после превышения лимита | `5m` |

> ⚠️ Никогда не коммитьте `.env` с реальными значениями. В production `JWT_SECRET` должен быть не менее 32 символов и храниться в секрет-менеджере (Vault, AWS SSM и т.д.).

---

## API

Сервис доступен на `http://localhost:8080/api/v1`.

Интерактивная Swagger-документация (Swagger UI) доступна на `http://localhost:8080/swagger/index.html`, JSON-спецификация — на `http://localhost:8080/swagger/doc.json`. Аннотации живут в `cmd/auth_service/main.go` (общая информация об API) и в `internal/api/http/handlers/handler.go` (описание конкретных эндпоинтов). После их изменения нужно перегенерировать `docs/`:

```bash
make swagger
```

### `POST /auth/reg`

Регистрация нового пользователя. Пароль должен содержать не менее 16 символов.

**Тело запроса:**
```json
{
  "login": "john",
  "password": "StrongPass12345678!"
}
```

**Ответ `201 Created`:**
```json
{
  "message": "user with login 'john' created"
}
```

### `POST /auth/login`

Аутентификация и получение JWT-токена.

**Тело запроса:**
```json
{
  "login": "john",
  "password": "StrongPass12345678!"
}
```

**Ответ `200 OK`:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Токен необходимо передавать в заголовке `Authorization: Bearer <token>` для защищённых эндпоинтов.

После `APP_THROTTLE_MAX_ATTEMPTS` неудачных попыток логина подряд для одного `login` дальнейшие запросы блокируются на `APP_THROTTLE_BLOCK_DURATION` и возвращают `429 Too Many Requests` — независимо от того, верен пароль или нет. Успешный логин сбрасывает счётчик попыток.

### Ошибки

Ошибки возвращаются в едином формате:
```json
{
  "error": "Wrong credentials",
  "details": "There is incorrect login or password"
}
```

| `error` | HTTP-статус |
|---|---|
| `Validation error` | `400 Bad Request` |
| `Wrong credentials` | `401 Unauthorized` |
| `Duplicate login` | `409 Conflict` |
| `Too many attempts` | `429 Too Many Requests` |
| `Internal error` | `500 Internal Server Error` |

---

## gRPC API

Помимо HTTP, сервис поднимает gRPC-сервер на отдельном порту (по умолчанию `:50051`, настраивается через `APP_GRPC_ADDR` / `grpc.addr` в `config.yaml`). Это тот же `service.Service`, что и у HTTP-хендлеров — бизнес-логика не дублируется, отличается только транспорт.

Контракт описан в [`proto/auth.proto`](proto/auth.proto) — сервис `auth.v1.AuthService` с двумя RPC: `Register` и `Login`. После изменения `.proto` перегенерировать стабы: `make proto`.

Сервер регистрирует [server reflection](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md) и стандартный [health-checking протокол](https://github.com/grpc/grpc/blob/master/doc/health-checking.md), поэтому его можно опрашивать через [`grpcurl`](https://github.com/fullstorydev/grpcurl) без локальной копии `.proto`:

```bash
# список сервисов
grpcurl -plaintext localhost:50051 list

# health-check
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check

# регистрация
grpcurl -plaintext -d '{"login":"john","password":"StrongPass12345678!"}' \
  localhost:50051 auth.v1.AuthService/Register

# логин
grpcurl -plaintext -d '{"login":"john","password":"StrongPass12345678!"}' \
  localhost:50051 auth.v1.AuthService/Login
```

Ошибки возвращаются как стандартные gRPC status-коды, а не HTTP-статусы:

| `models.ErrorCode` | gRPC `codes.Code` |
|---|---|
| `CodeValidationError` | `InvalidArgument` |
| `CodeErrDuplicate` | `AlreadyExists` |
| `CodeWrongCredentials` | `Unauthenticated` |
| `CodeTooManyAttempts` | `ResourceExhausted` |
| `CodeInternalError` | `Internal` |

---

## Тестирование

Проект покрыт двумя видами тестов.

### Unit-тесты

Тестируют каждый слой изолированно. Репозиторий и внешние зависимости подменяются моками, сгенерированными через `go.uber.org/mock`.

```bash
make unit
# запускает: go test ./internal/... с профилем покрытия
```

### Интеграционные тесты

Поднимают реальную PostgreSQL через отдельный `docker compose` и проверяют сквозные сценарии. Redis-троттлинг логина тестируется через [`miniredis`](https://github.com/alicebob/miniredis) (in-process фейк) — отдельного Redis-контейнера для тестов не требуется.

```bash
make integration
# поднимает tests/docker/compose.yaml → запускает тесты → останавливает контейнеры
```

### Просмотр покрытия

```bash
# unit
make unit
make cover f=coverage_unit.out

# integration
make integration
make cover f=coverage_integration.out
```

#### Известное ограничение тестов

`TestTxManager` не выполняет реальный коммит транзакций. Это означает, что тест-сценарии, использующие несколько транзакций над одними и теми же данными (в том числе негативные кейсы с откатом), корректно проверить через него нельзя. Планируется заменить на полноценный `testcontainers-go` с реальным rollback.

---

## Разработка

### Миграции

Файлы миграций хранятся в `internal/migrations/`. Применяются автоматически при старте `docker compose`. Для ручного запуска используется образ `migrate/migrate:4`.

### Линтинг

```bash
golangci-lint run ./...
```

---

## Лицензия

MIT