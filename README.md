# Контакт — Go API

Go 1.25.5, PostgreSQL, pgxpool, net/http, validator, envconfig, zap, Swagger.
Версии обязательных библиотек закреплены в go.mod; транзитивные зависимости — в go.sum.

## Запуск в Windows / VS Code

Откройте эту папку в VS Code. Требуются Go 1.25.5+, Docker Desktop с Linux containers, GNU Make (для команд make).
Запустите Docker Desktop, затем в PowerShell:

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
PostgreSQL: localhost:5433, база kontakt. Локальные параметры в .env; .env не следует коммитить. Токен администратора предназначен для локальной разработки, замените его перед размещением приложения.

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

Сервисы зависят от интерфейсов repository. Транспорт не выполняет SQL. Заказ не вызывает HTTP каталога: репозиторий заказов читает и блокирует товары в той же транзакции, что создаёт заказ. Это обеспечивает атомарность списания остатков. Состав заказа хранится как JSONB-снимок названий, количества и цены на момент покупки; последующее редактирование товара не меняет историю.

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

Для администратора заголовок `Authorization: Bearer <APP_ADMIN_TOKEN>`. В Swagger нажмите Authorize и введите `Bearer ` и значение из .env. Общего публичного списка заказов нет, чтобы не раскрывать телефоны и адреса покупателей.

Пример заказа:

```json
{"name":"Тестовый покупатель","phone":"+70000000000","delivery":"pickup","address":"","items":[{"id":"00000000-0000-4000-8000-000000000002","qty":1}]}
```

Цены, наименования и итог сервер получает из PostgreSQL. При недостатке остатка возвращает 409; вся транзакция откатывается. Повторяющиеся позиции в одном заказе отклоняются. Порядок блокировки товаров одинаков для конкурентных заказов. Статусы: Новый, Подтверждён, Собирается, Готов к выдаче, Завершён. Отмена с возвратом остатков не реализована.

## Версии базы и команды

Миграции 000001 создают таблицы, 000002 добавляют демонстрационные товары. migrate хранит текущую версию и dirty-флаг в schema_migrations. Номер миграции — версия структуры базы; тег postgres:18.4-bookworm — версия СУБД.

```powershell
make migrate-version
make migrate-up
make migrate-down
make build
make swagger
make test
make logs
make down
```

`migrate-down` откатывает одну миграцию: откат seed удаляет демонстрационные товары, откат schema удаляет таблицы. `down` останавливает контейнеры, файлы базы остаются в out/pggdata. Не изменяйте уже применённые миграции: добавляйте следующий номер с парой up/down. Для PostgreSQL 18 каталог монтируется в /var/lib/postgresql согласно официальному образу: https://hub.docker.com/_/postgres .

## Подключение React

Фронтенд из предыдущего этапа пока использует localStorage; этот архив содержит отдельный сервер. Для подключения замените чтение товаров на GET /api/v1/products, отправку заказа — на POST /api/v1/orders. Используйте UUID из ответа API вместо числовых идентификаторов демофронтенда. Корзина и избранное могут остаться в localStorage. CORS разрешает один origin из APP_CORS_ORIGIN: при запуске React на http://127.0.0.1:5173 измените .env соответствующим образом. Администраторский токен нельзя вшивать в общедоступную сборку React. Личные кабинеты и пользовательская авторизация требуют отдельного этапа; текущий API не предоставляет покупателю историю заказов между устройствами.

## Проверка

Сборка Go и проверка конфигурации Compose выполнены. Автоматические тесты проверяют валидацию HTTP-входа и защиту административных маршрутов. Проверка реальных миграций и транзакций требует запущенного Docker Engine.
