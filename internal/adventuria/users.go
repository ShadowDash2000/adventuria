package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"iter"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

type Users struct {
	users *MemoryCacheWithClose[string, User]
}

func NewUsers(ctx AppContext) *Users {
	u := &Users{
		users: NewMemoryCacheWithClose[string, User](ctx, time.Hour, false),
	}
	u.bindHooks(ctx)
	return u
}

func (u *Users) bindHooks(ctx AppContext) {
	ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionUsers).BindFunc(func(e *core.RecordEvent) error {
		user, err := u.GetByID(AppContext{App: e.App}, e.Record.Id)
		if err != nil {
			return e.Next()
		}
		user.SetProxyRecord(e.Record)

		return e.Next()
	})
}

func (u *Users) GetByID(ctx AppContext, userId string) (User, error) {
	user, ok := u.users.Get(userId)
	if ok {
		return user, nil
	}

	user, err := NewUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	u.users.Set(userId, user)
	return user, nil
}

func (u *Users) GetByName(ctx AppContext, name string) (User, error) {
	for _, user := range u.users.GetAll() {
		if name == user.Name() {
			return user, nil
		}
	}

	user, err := NewUserFromName(ctx, name)
	if err != nil {
		return nil, err
	}

	u.users.Set(user.ID(), user)
	return user, nil
}

func (u *Users) GetAll(ctx AppContext) (iter.Seq2[string, User], error) {
	var users []struct {
		Id string `db:"id"`
	}
	err := ctx.App.RecordQuery(schema.CollectionUsers).
		Select(schema.UserSchema.Id).
		All(&users)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		_, ok := u.users.Get(user.Id)
		if ok {
			continue
		}

		user, err := NewUser(ctx, user.Id)
		if err != nil {
			return nil, err
		}

		u.users.Set(user.ID(), user)
	}

	return u.GetAllInMemory(), nil
}

func (u *Users) GetAllInMemory() iter.Seq2[string, User] {
	return u.users.GetAll()
}
