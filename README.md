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

## Переменные окружения

Для работы некоторых компонентов требуются переменные из env.\
Можно создать .env файл в корне проекта и он автоматически подгрузится.

Ключ Twitch используются для парсинга игр с IGDB, а так же для получения\
статуса стримов игроков.
Ключ можно получить здесь: https://dev.twitch.tv/console/apps/create 
```
TWITCH_CLIENT_ID=***
TWITCH_CLIENT_SECRET=***
```
Строка из этого параметра применяется в качестве "where = ..." для фильтрации игр при парсинге.\
Переменные для фильтрации https://api-docs.igdb.com/#game
```
IGDB_PARSE_FILTER="game_type = 0 & platforms = (6)"
```

Frontend: https://github.com/ShadowDash2000/adventuria-react