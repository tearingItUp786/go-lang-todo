package models

import (
	"database/sql"
	"fmt"
	"log"
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
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
		panic(err)
	}

	return todos, nil
}

func (b *BaseModel) GetSingleToDo(todoId string) (ToDo, error) {
	row := b.db.QueryRow("SELECT * FROM todo where id=$1;", todoId)

	// print rows as json
	todo := ToDo{}

	if err := row.Scan(&todo.Id, &todo.Text, &todo.Done); err != nil {
		log.Fatal(err)
	}

	if err := row.Err(); err != nil {
		log.Fatal(err)
		panic(err)
	}

	return todo, nil
}

func (b *BaseModel) UpdateSingleToDo(todoId string, text string, done bool) (ToDo, error) {
	todo := ToDo{}

	err := b.db.QueryRow(`UPDATE todo
			SET text = $1, done = $2
			WHERE id = $3 RETURNING *;`, text, done, todoId).Scan(&todo.Id, &todo.Text, &todo.Done)
	if err != nil {
		return ToDo{}, err
	}

	return todo, nil
}

func (b *BaseModel) InsertToDo(todo string) {
	_, err := b.db.Exec("INSERT INTO todo (text, done) VALUES ($1, $2)", todo, false)
	if err != nil {
		fmt.Println(err)
	}
}

func (b *BaseModel) ToggleTodo(todoId string) (ToDo, error) {
	todo := ToDo{}
	err := b.db.QueryRow(`UPDATE todo
			SET done = NOT done
			WHERE id = $1 RETURNING *;`, todoId).Scan(&todo.Id, &todo.Text, &todo.Done)
	if err != nil {
		return ToDo{}, err
	}
	return todo, nil
}

func (b *BaseModel) DeleteTodo(todoId string) error {
	_, err := b.db.Exec(`DELETE FROM todo where id = $1`, todoId)
	if err != nil {
		return err
	}
	return nil
}
