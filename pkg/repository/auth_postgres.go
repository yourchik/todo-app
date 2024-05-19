package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/yourchik/todo-app/domain"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user domain.User) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (name, username, password_hash) 
                      SELECT $1, LOWER($2), $3 
                      WHERE NOT EXISTS (
                          SELECT 1 FROM %s WHERE LOWER(username) = LOWER($2)) 
                      RETURNING id`, usersTable, usersTable)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)
	return user, err
}

func (r *AuthPostgres) GetAll(userId int) ([]domain.TodoList, error) {
	var lists []domain.TodoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable)
	if err := r.db.Select(&lists, query, userId); err != nil {
		return nil, err
	}
	return lists, nil
}
