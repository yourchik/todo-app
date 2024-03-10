package service

import "github.com/yourchik/todo-app/pkg/repository"

type TodoList interface {
}

type Authorization interface {
}

type TodoItem interface {
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
