package user

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	FindOne(ctx context.Context, id string) (User, error) //Метод пользователя по ID
	FindAll(ctx context.Context) (u []User, err error)
	Update(ctx context.Context, user string) error
	Delete(ctx context.Context, id string) error
}
