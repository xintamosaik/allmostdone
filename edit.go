package main

import (
	"context"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

 
const editForm = `
<h1>Edit Todo</h1>
<form action="/todos/{{.ID}}/update" method="post">
  <label for="short">Short:</label><br>
  <input type="text" id="short" name="short" value="{{.Short}}"><br>

  <label for="description">Description:</label><br>
  <textarea id="description" name="description">{{.Description}}</textarea><br>

  <label for="due_date">Due Date (YYYY-MM-DD):</label><br>
  <input type="text" id="due_date" name="due_date" value="{{if .DueDate}}{{.DueDate.Format "2006-01-02"}}{{end}}"><br>

  <label for="cost_of_delay">Cost of Delay:</label><br>
  <input type="number" id="cost_of_delay" name="cost_of_delay" value="{{.CostOfDelay}}"><br>

  <label for="effort">Effort:</label><br>
  <input type="text" id="effort" name="effort" value="{{.Effort}}"><br><br>

  <input type="submit" value="Update">
</form>
`
 
func updateTodo(conn *pgx.Conn, id int, short string, description string, dueDate *time.Time, costOfDelay int16, effort string) error {
	_, err := conn.Exec(
		context.Background(),
		`UPDATE todos
         SET short=$1,
             description=$2,
             due_date=$3,
             cost_of_delay=$4,
             effort=$5,
             updated_at=now()
         WHERE id=$6`,
		short,
		description,
		dueDate,
		costOfDelay,
		effort,
		id,
	)

	return err
}
func updateHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.NotFound(w, r)
			return
		}

		short, description, dueDate, costOfDelay, effort, err := parseTodoForm(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = updateTodo(conn, id, short, description, dueDate, costOfDelay, effort); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		listHandler(conn)(w, r)
	}
}
 
func getTodo(conn *pgx.Conn, id int) (Todo, error) {
	var t Todo

	err := conn.QueryRow(
		context.Background(),
		`SELECT id, short, description, due_date, cost_of_delay, effort, created_at, updated_at
         FROM todos
         WHERE id=$1`,
		id,
	).Scan(&t.ID, &t.Short, &t.Description, &t.DueDate, &t.CostOfDelay, &t.Effort, &t.CreatedAt, &t.UpdatedAt)

	return t, err
}

// editHandler serves the form to edit an existing todo item. Path: /todos/{id}/edit
func editHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 3 {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			http.NotFound(w, r)
			return
		}
		todo, err := getTodo(conn, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t, err := template.New("webpage").Parse(editForm)
		check(err)
		t.Execute(w, todo)
	}
}
