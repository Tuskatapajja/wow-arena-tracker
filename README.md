# WoW Arena Tracker

REST API для отслеживания статистики арены в WoW WotLK.

## Технологии
- Go 1.21
- PostgreSQL
- Docker & Docker Compose

## Эндпоинты

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/api/match` | Сохранить результат боя |
| GET | `/api/stats` | Получить статистику |

## Запуск

```bash
docker compose up -d
curl http://localhost:8080/api/stats
