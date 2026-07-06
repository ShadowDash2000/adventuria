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

Предметы в игре представляют собой набор эффектов, которые подписываются на игровые события.
Существует множество готовых эффектов, имплементация которых лежит здесь: `internal/adventuria/effects/custom`

Для реализации своего эффекта структура должна имплементировать интерфейс
`model.Effect internal/adventuria/model/effect.go`.
Далее нужно зарегистрировать новый эффект при старте приложения:

```go
import (
    "adventuria/internal/adventuria/effects"
    "adventuria/internal/adventuria/model"
)

effects.Register(
    effects.NewEffectDef(
        "effect_name",
        func (effect model.EffectInfo) model.Effect {
            return &SomeNewEffect{
                EffectBase: effects.NewEffectBase(effect),
            }
        },
    ),
)
```

В большинстве случаев эффектам не нужно самостоятельно вызывать сохранение данных `player`, так как в конце каждого
действия `action` вызывается сохранение `player.progress` и его `player.lastAction`.
Есть случаи, когда эффект создаёт новое действие, поэтому ему нужно самостоятельно вызвать сохранение
`player.lasAction`.

## Клетки ♿

Примеры готовых игровых клеток: `internal/adventuria/cells/custom`.\
Клетки должны имплементировать следующие методы:

```go
// Вызывается в момент, когда игрок наступает на клетку
OnCellReached(ctx context.Context, events *Events, player *Player, reachedCtx *ReachedContext) error

// Вызывается в момент, когда игрок покидает клетку
OnCellLeft(ctx context.Context, events *Events, player *Player) error

// Опционально: Вызывается при сохранении клетки в PocketBase для проверки значения в поле "value"
Verify(ctx context.Context, value string) error
```

В клетках можно вызывать сохранение `player.progress` и `player.lastAction`, если на то есть причина. Например, если
клетке нужно
создать новый ход, то это можно сделать таким образом:

```go
// Полный пример: internal/adventuria/actions/custom/reroll/action.go

lastAction := player.LastAction()
lastAction.SetType(Type)
lastAction.SetReview(review.ID())
_, err = r.actions.Save(ctx, lastAction)
if err != nil {
    return nil, err
}

newAction, err := model.NewAction(uuid.New(), model.ActionCreate{
    Player: player.ID(),
    Cell:   currentCell.Data().ID(),
    Type:   Type,
})
if err != nil {
    return nil, err
}

player.SetLastAction(newAction)
```

После выполнения `OnCellReached` и `OnCellLeft`, так же сохраняются поля `player.progress` и его `player.lastAction`.

> [!WARNING]
> Если клетка не представляет собой цепочку вызова `action`, то в `OnCellReached` обязательно нужно
> вызывать `player.LastAction().SetCanMove(true)` для того, чтобы игрок мог идти дальше.

## Действия 🎲

Примеры реализации действий: `internal/adventuria/actions/custom`.\
В своей основе действия нужны для манипуляции над данными игрока, например, чтобы завершить прохождение игры на клетке,
дропнуть, купить предмет в магазине и т.д. Действия обязаны имплементировать два метода:

```go
// Проверяет, может ли игрок выполнить данное действие в текущий момент
CanDo(ctx context.Context, events *Events, player *Player) bool

Do(ctx context.Context, events *Events, player *Player, actionReq ActionRequest) (any, error)
```

После выполнения `Do()` так же, вызывается сохранение `player.progress` и его `player.lastAction`.

## Общее

Если действию или эффекту нужно передать какую-то view информация для фронта, то для этого они должны имплементировать
интерфейс `WithView`.

```go
GetView(ctx context.Context, events *Events, player *Player) (any, error)
```

Пример реализации на эффекте: `internal/adventuria/effects/custom/paid_movement_in_radius/view.go`\
Пример реализации на действии: `internal/adventuria/actions/custom/buy/view.go`

## Тестирование 🔧

В процессе...

## Taskfile

Для генерации бойлерплейт кода используются команды из Taskfile.yml.
Если установлен Go, то можно скачать и собрать из source-кода через команду:

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

В остальных случаях смотреть сюда: https://taskfile.dev/docs/installation.

После установки просто выполняем команду `task`, которая покажет все доступные команды.

## Остальное

Frontend: https://github.com/ShadowDash2000/adventuria-react