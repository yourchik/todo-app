package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/yourchik/todo-app"
)

type TodoList interface {
}

type Authorization interface {
	CreateUser(user todo.User) (int, error)
}

type TodoItem interface {
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
