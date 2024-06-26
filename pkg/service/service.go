package service

import (
	"github.com/yourchik/todo-app/domain"
	"github.com/yourchik/todo-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoItem interface {
	Create(userId, listId int, item domain.TodoItem) (int, error)
	GetAll(userId, listId int) ([]domain.TodoItem, error)
	GetById(userId, itemId int) (domain.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input domain.UpdateItemInput) error
}

type TodoList interface {
	Create(userId int, list domain.TodoList) (int, error)
	GetAll(userId int) ([]domain.TodoList, error)
	GetById(userId, listId int) (domain.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input domain.UpdateListInput) error
}
type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
