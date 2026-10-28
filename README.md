# Electronics store — Go API


## Запуск в Windows
PowerShell:
```powershell
make deps
make up
make run
```

Без Make:

```powershell
go mod download
docker compose up -d todoapp-postgres port-forwarder
docker compose run --rm todoapp-postgres-migrate up
powershell -NoProfile -ExecutionPolicy Bypass -File scripts/run.ps1
```

API: http://localhost:8080/api/v1/products
Swagger: http://localhost:8080/swagger/index.html
Health: http://localhost:8080/health
PostgreSQL: localhost:5433, 

## Структура

```text
cmd/app                         точка входа, сборка зависимостей
docs                            docs.go, swagger.json, swagger.yaml
internal/core/config            envconfig
internal/core/domain            модели
internal/core/errors            ошибки приложения
internal/core/logger            zap
internal/core/repository/postgres/pool   подключение pgxpool
internal/core/repository/postgres/pgx    преобразование ошибок
internal/core/transport/http     middleware, request, response, server, types
internal/features/catalog       repository → service → transport
internal/features/orders        repository → service → transport
migrations                      последовательные up/down SQL
out/logs                        журнал приложения
out/pggdata                     файлы PostgreSQL
```


## API

| Метод | Путь /api/v1 | Доступ |
|---|---|---|
| GET | /products | Все товары |
| GET | /products/{uuid} | Один товар |
| GET | /categories | Категории |
| POST | /products | Администратор |
| PUT | /products/{uuid} | Администратор |
| POST | /orders | Создать заказ |
| GET | /orders | Администратор, последние 500 |
| PATCH | /orders/{uuid}/status | Администратор |

