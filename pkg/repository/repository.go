package repository

import "github.com/jmoiron/sqlx"

type TodoList interface {
}

type Authorization interface {
}

type TodoItem interface {
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(*sqlx.DB) *Repository {
	return &Repository{}
}
