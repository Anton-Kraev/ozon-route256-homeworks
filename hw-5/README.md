## Запуск приложения

### Переменные окружения:
В .env файле в корневой директории должны быть переменные:
```
POSTGRES_DB=...
POSTGRES_USER=...
POSTGRES_PASSWORD=...
PGDATA=...
DATABASE_URL=...
```
### Запуск миграций:
```shell
./migration.sh up
```

### Сборка и запуск приложения в докере
```shell
docker-compose up --build -d postgres app
```
