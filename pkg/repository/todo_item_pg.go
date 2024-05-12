package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/yourchik/todo-app"
	"strings"
)

type ToDoItemPostgres struct {
	db *sqlx.DB
}

func NewToDoItemPostgres(db *sqlx.DB) *ToDoItemPostgres {
	return &ToDoItemPostgres{db: db}
}

func (r *ToDoItemPostgres) Create(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var itemId int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id",
		todoItemsTable)

	row := tx.QueryRow(createListQuery, item.Title, item.Description)
	if err := row.Scan(&itemId); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)",
		listsItemsTable)
	_, err = tx.Exec(createUsersListQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return itemId, tx.Commit()
}

func (r *ToDoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on ti.id = li.item_id 
									INNER JOIN %s ul on ul.list_id = li.list_id 
                                     WHERE li.list_id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Select(&items, query, listId); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ToDoItemPostgres) GetById(userId, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on ti.id = li.item_id 
									INNER JOIN %s ul on ul.list_id = li.list_id
									 WHERE ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Get(&item, query, itemId); err != nil {
		return item, err
	}
	return item, nil
}

func (r *ToDoItemPostgres) Delete(userId, itemId int) error {
	_, err := r.db.Exec(fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul 
												WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND li.item_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable), userId, itemId)
	return err
}

func (r *ToDoItemPostgres) Update(userId, itemId int, input todo.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	serQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul 
									WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND li.item_id = $%d`,
		todoItemsTable, serQuery, listsItemsTable, usersListsTable, argId, argId+1)
	args = append(args, userId, itemId)

	logrus.Debugf("update query: %s", query)
	logrus.Debugf("update args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err

}
