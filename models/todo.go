package models

import (
	"database/sql"
)

type ToDo struct {
	Id     int64
	Text   string
	Done   bool
	UserId int64
}

type BaseModel struct {
	DB *sql.DB
}

func NewBaseModel(db *sql.DB) *BaseModel {
	return &BaseModel{
		DB: db,
	}
}

func (b *BaseModel) GetTodos(userId int) ([]ToDo, error) {
	rows, _ := b.DB.Query("SELECT * FROM todos WHERE user_id=$1 ORDER BY id DESC;", userId)
	defer rows.Close()
	// print rows as json
	todos := []ToDo{}

	for rows.Next() {
		var todo ToDo
		if err := rows.Scan(&todo.Id, &todo.Text, &todo.Done, &todo.UserId); err != nil {
			return []ToDo{}, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return []ToDo{}, err
	}

	return todos, nil
}

func (b *BaseModel) GetSingleToDo(todoId string, userId int) (ToDo, error) {
	row := b.DB.QueryRow("SELECT * FROM todos WHERE id=$1 AND user_id=$2;", todoId, userId)

	// print rows as json
	todo := ToDo{}

	if err := row.Scan(&todo.Id, &todo.Text, &todo.Done, &todo.UserId); err != nil {
		return ToDo{}, err
	}

	if err := row.Err(); err != nil {
		return ToDo{}, err
	}

	return todo, nil
}

func (b *BaseModel) UpdateSingleToDo(
	todoId string,
	text string,
	done bool,
	userId int,
) (ToDo, error) {
	todo := ToDo{}

	row := b.DB.QueryRow(`UPDATE todos
			SET text = $1, done = $2
			WHERE id = $3 AND user_id = $4 RETURNING *;`, text, done, todoId, userId)

	if err := row.Scan(&todo.Id, &todo.Text, &todo.Done, &todo.UserId); err != nil {
		return ToDo{}, err
	}

	return todo, nil
}

func (b *BaseModel) InsertToDo(todoText string, userId int) (ToDo, int, error) {
	todo := ToDo{}
	row := b.DB.QueryRow(
		"INSERT INTO todos (text, done, user_id) VALUES ($1, $2, $3) RETURNING *;",
		todoText,
		false,
		userId,
	)
	if err := row.Scan(&todo.Id, &todo.Text, &todo.Done, &todo.UserId); err != nil {
		return ToDo{}, 0, err
	}

	todoCount := b.DB.QueryRow("SELECT COUNT(*) FROM todos where user_id = $1;", userId)

	var count int
	err := todoCount.Scan(&count)
	if err != nil {
		return ToDo{}, 0, err
	}

	return todo, count, nil
}

func (b *BaseModel) ToggleTodo(todoId string, userId int) (ToDo, error) {
	todo := ToDo{}
	row := b.DB.QueryRow(`UPDATE todos
			SET done = NOT done
			WHERE id = $1 AND user_id = $2 RETURNING *;`, todoId, userId)

	if err := row.Scan(&todo.Id, &todo.Text, &todo.Done, &todo.UserId); err != nil {
		return ToDo{}, err
	}
	return todo, nil
}

func (b *BaseModel) DeleteTodo(todoId string, userId int) (int, error) {
	// Deleting the specified todo
	_, err := b.DB.Exec("DELETE FROM todos WHERE id = $1 AND user_id = $2", todoId, userId)
	if err != nil {
		return 0, err
	}

	// Getting the count of remaining todos
	var remainingCount int
	row := b.DB.QueryRow("SELECT COUNT(*) FROM todos WHERE user_id = $1", userId)

	if err := row.Scan(&remainingCount); err != nil {
		return 0, err
	}

	return remainingCount, nil
}

func (b *BaseModel) DeleteAll(userId int) error {
	// Deleting the specified todo
	_, err := b.DB.Exec("DELETE FROM todos where user_id=$1;", userId)
	if err != nil {
		return err
	}

	return nil
}

func (b *BaseModel) BulkInsertToDos(bulkTodos []ToDo, userId int) ([]ToDo, error) {
	tx, err := b.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO todos (text, done, user_id) VALUES ($1, $2, $3)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var todos []ToDo
	for _, bulkTodo := range bulkTodos {

		if _, err := stmt.Exec(bulkTodo.Text, false, userId); err != nil {
			return nil, err
		}

		todo := ToDo{Text: bulkTodo.Text, Done: bulkTodo.Done}
		todos = append(todos, todo)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return todos, nil
}
