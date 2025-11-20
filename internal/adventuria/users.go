package adventuria

import (
	"adventuria/pkg/cache"
	"time"
)

type Users struct {
	users cache.Cache[string, User]
}

func NewUsers() *Users {
	return &Users{
		users: cache.NewMemoryCacheWithClose[string, User](time.Hour, false),
	}
}

func (u *Users) GetByID(userId string) (User, error) {
	user, ok := u.users.Get(userId)
	if ok {
		return user, nil
	}

	user, err := NewUser(userId)
	if err != nil {
		return nil, err
	}

	u.users.Set(userId, user)
	return user, nil
}

func (u *Users) GetByName(name string) (User, error) {
	for _, user := range u.users.GetAll() {
		if name == user.Name() {
			return user, nil
		}
	}

	user, err := NewUserFromName(name)
	if err != nil {
		return nil, err
	}

	u.users.Set(user.ID(), user)
	return user, nil
}
