package models

import (
	"database/sql"
)

type ToDo struct {
	Id   int64
	Text string
	Done bool
}

type BaseModel struct {
	db *sql.DB
}

func NewBaseModel(db *sql.DB) *BaseModel {
	return &BaseModel{
		db: db,
	}
}

func (b *BaseModel) GetTodos() ([]ToDo, error) {
	rows, _ := b.db.Query("SELECT * FROM todo order by id desc;")
	defer rows.Close()
	// print rows as json
	todos := []ToDo{}

	for rows.Next() {
		var todo ToDo
		if err := rows.Scan(&todo.Id, &todo.Text, &todo.Done); err != nil {
			return []ToDo{}, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return []ToDo{}, err
	}

	return todos, nil
}

func (b *BaseModel) GetSingleToDo(todoId string) (ToDo, error) {
	row := b.db.QueryRow("SELECT * FROM todo where id=$1;", todoId)

	// print rows as json
	todo := ToDo{}

	if err := row.Scan(&todo.Id, &todo.Text, &todo.Done); err != nil {
		return ToDo{}, err
	}

	if err := row.Err(); err != nil {
		return ToDo{}, err
	}

	return todo, nil
}

func (b *BaseModel) UpdateSingleToDo(todoId string, text string, done bool) (ToDo, error) {
	todo := ToDo{}

	row := b.db.QueryRow(`UPDATE todo
			SET text = $1, done = $2
			WHERE id = $3 RETURNING *;`, text, done, todoId)

	if err := row.Scan(&todo.Id, &todo.Text, &todo.Done); err != nil {
		return ToDo{}, err
	}

	return todo, nil
}

func (b *BaseModel) InsertToDo(todoText string) (ToDo, int, error) {
	todo := ToDo{}
	row := b.db.QueryRow(
		"INSERT INTO todo (text, done) VALUES ($1, $2) RETURNING *;",
		todoText,
		false,
	)
	if err := row.Scan(&todo.Id, &todo.Text, &todo.Done); err != nil {
		return ToDo{}, 0, err
	}

	todoCount := b.db.QueryRow("SELECT COUNT(*) FROM todo;")

	var count int
	err := todoCount.Scan(&count)
	if err != nil {
		return ToDo{}, 0, err
	}

	return todo, count, nil
}

func (b *BaseModel) ToggleTodo(todoId string) (ToDo, error) {
	todo := ToDo{}
	row := b.db.QueryRow(`UPDATE todo
			SET done = NOT done
			WHERE id = $1 RETURNING *;`, todoId)

	if err := row.Scan(&todo.Id, &todo.Text, &todo.Done); err != nil {
		return ToDo{}, err
	}
	return todo, nil
}

func (b *BaseModel) DeleteTodo(todoId string) (int, error) {
	// Deleting the specified todo
	_, err := b.db.Exec("DELETE FROM todo WHERE id = $1", todoId)
	if err != nil {
		return 0, err
	}

	// Getting the count of remaining todos
	var remainingCount int
	row := b.db.QueryRow("SELECT COUNT(*) FROM todo")

	if err := row.Scan(&remainingCount); err != nil {
		return 0, err
	}

	return remainingCount, nil
}

func (b *BaseModel) DeleteAll() error {
	// Deleting the specified todo
	_, err := b.db.Exec("DELETE FROM todo")
	if err != nil {
		return err
	}

	return nil
}

func (b *BaseModel) BulkInsertToDos(bulkTodos []ToDo) ([]ToDo, error) {
	tx, err := b.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO todo (text, done) VALUES ($1, $2)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var todos []ToDo
	for _, bulkTodo := range bulkTodos {

		if _, err := stmt.Exec(bulkTodo.Text, false); err != nil {
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
