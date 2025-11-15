# Adventuria / Приключпопия

## Запуск

#### Запуск приложения:

```bash
go run cmd/main.go serve
```

При первом запуске в консоли будет выведена ссылка для создания superuser в PocketBase.\
Либо можно создать пользователя через консоль:

```bash
go run cmd/main.go superuser create EMAIL PASS
```

#### Запуск миграции:

```bash
go run cmd/main.go migrate up
```

## Docker

#### Запуск через docker-compose:

```bash
docker-compose up --build -d
```

Frontend: https://github.com/ShadowDash2000/adventuria-react