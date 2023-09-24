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
	rows, _ := b.db.Query("SELECT * FROM todo")
	defer rows.Close()
	// print rows as json
	todos := []ToDo{}

	for rows.Next() {
		var todo ToDo
		if err := rows.Scan(&todo.Id, &todo.Text, &todo.Done); err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
		fmt.Printf("%d is %v and %v \n", todo.Id, todo.Text, todo.Done)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return todos, nil
}

func (b *BaseModel) InsertToDo(todo string) {
	_, err := b.db.Exec("INSERT INTO todo (todo, done) VALUES ($1, $2)", todo, false)
	if err != nil {
		fmt.Println(err)
	}
}
