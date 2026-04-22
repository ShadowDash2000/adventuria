# Adventuria / Приключпопия 🎡

## Что это?

Бэкенд для проведения ивента по прохождению рандомных игр, где игроки перемещаются по игровому полю с помощью
броска кубиков. \
Список игр подгружается с IGDB.

## Запуск 🚀

#### Запуск приложения:

```bash
go run cmd/main.go serve
```

При первом запуске в консоли будет выведена ссылка для создания superuser в PocketBase.
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

При использовании `docker-compose` будут созданы два контейнера. Один для запуска приложения ("api"), а второй
для выполнения миграции ("migrator").

## Переменные окружения

Для работы некоторых компонентов требуются переменные из env.
Можно создать .env файл в корне проекта и он автоматически подгрузится.

Ключ Twitch используются для парсинга игр с IGDB, а также для получения статуса стримов игроков.
Ключ можно получить здесь: https://dev.twitch.tv/console/apps/create 
```
TWITCH_CLIENT_ID=***
TWITCH_CLIENT_SECRET=***
```
Ключ YouTube API используется для получения статуса стримов с YouTube.
```
YOUTUBE_API_KEY=***
```
Строка из этого параметра применяется в качестве "where = ..." для фильтрации игр при парсинге.
Переменные для фильтрации https://api-docs.igdb.com/#game
```
IGDB_PARSE_FILTER="game_type = 0 & platforms = (6)"
```

## Предметы и эффекты 📦✨

Предметы в игре представляют собой набор эффектов, которые подписываются на события игроков (`player`).
Существует множество готовых эффектов, имплементация которых лежит здесь: `internal/adventuria/effects`

Для реализации своего эффекта структура должна имплементировать интерфейс `Effect internal/adventuria/effect.go`.
Далее нужно зарегистрировать новый эффект при старте приложения:
```go
adventuria.RegisterEffects(map[string]adventuria.EffectCreator{
    "myNewEffect": adventuria.NewEffect(func() adventuria.Effect { return &MyNewEffect{} }),
})
```

В большинстве случаев эффектам не нужно самостоятельно вызывать сохранение данных `player`, так как в конце каждого
действия `action` вызывается сохранение `player.progress` и его `player.lastAction`.
Есть случаи, когда эффект создаёт новое действие, поэтому ему нужно самостоятельно вызвать сохранение
`player.lasAction`.

## Клетки ♿

Примеры готовых игровых клеток: `internal/adventuria/cells`.\
Клетки должны имплементировать следующие методы:
```go
// Вызывается в момент, когда игрок наступает на клетку
OnCellReached(*CellReachedContext) error
// Вызывается в момент, когда игрок покидает клетку
OnCellLeft(*CellLeftContext) error
// Опционально: Вызывается при сохранении клетки в PocketBase для проверки значения в поле "value"
Verify(string) error
```
В клетках можно вызывать сохранение `player.progress` и `player.lastAction`, если на то есть причина. Например, если клетке нужно
создать новый ход, то это можно сделать таким образом:

```go
action := player.LastAction()
action.SetType("SOME_ACTION_TYPE")
err := adventuria.PocketBase.Save(action.ProxyRecord())
if err != nil {
    return err
}
// Помечаем record экшена, как новый, чтобы создать новую запись
action.MarkAsNew()
```

После выполнения `OnCellReached` и `OnCellLeft`, так же сохраняются поля `player.progress` и его `player.lastAction`.

> [!WARNING]
> Если клетка не представляет собой цепочку вызова `action`, то в `OnCellReached` обязательно нужно
> вызывать `player.LastAction().SetCanMove(true)` для того, чтобы игрок мог идти дальше.

## Действия 🎲

Примеры реализации действий: `internal/adventuria/actions`.\
В своей основе действия нужны для манипуляции над данными игрока, например, чтобы завершить прохождение игры на клетке,
дропнуть, купить предмет в магазине и т.д. Действия обязаны имплементировать два метода:
```go
CanDo(ActionContext) bool
Do(ActionContext, ActionRequest) (*ActionResult, error)
// Этот метод используется для получения JSON'а фронтом.
// В некоторых случаях перед выдачей ответа нужно модифицировать данные.
// Например, в магазине предметов требуется применить активные эффекты для
// скидки на товары (internal/adventuria/actions/buy.go).
GetVariants(ActionContext) any
```

После выполнения `Do()` так же, вызывается сохранение `player.progress` и его `player.lastAction`.

## Тестирование 🔧

Для написания тестов существует отдельная структура `GameTest`, которая запускает приложение
так же, как и при обычном запуске, но без `http` сервера.

```go
game, err := tests.NewGameTest()
if err != nil {
    t.Fatal(err)
}
```

Frontend: https://github.com/ShadowDash2000/adventuria-react