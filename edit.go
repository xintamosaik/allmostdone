package main

import (
	"context"
	"html/template"
	"net/http"
	"strconv"

	"time"

	"github.com/jackc/pgx/v5"
)

const editFormHTML = `
<h1>Edit Todo</h1>
<form
  action="/todos/{{.ID}}/update"
  method="post"
  fx-action="/todos/{{.ID}}/update"
  fx-method="POST"
  fx-target="#output"
  fx-swap="innerHTML">
  <label for="short">Short:</label><br>
  <input type="text" id="short" name="short" value="{{.Short}}"><br>

  <label for="description">Description:</label><br>
  <textarea id="description" name="description">{{.Description}}</textarea><br>

  <label for="due_date">Due Date (YYYY-MM-DD):</label><br>
  <input type="date" id="due_date" name="due_date" value="{{if .DueDate}}{{.DueDate.Format "2006-01-02"}}{{end}}">

  <label for="cost_of_delay">Cost of Delay:</label><br>
  <input type="number" id="cost_of_delay" name="cost_of_delay" min="-2" max="2" value="{{.CostOfDelay}}"><br>

  <label for="effort">Effort:</label><br>
  <select id="effort" name="effort">
    <option value="mins" {{if eq .Effort "mins"}}selected{{end}}>mins</option>
    <option value="hours" {{if eq .Effort "hours"}}selected{{end}}>hours</option>
    <option value="days" {{if eq .Effort "days"}}selected{{end}}>days</option>
    <option value="weeks" {{if eq .Effort "weeks"}}selected{{end}}>weeks</option>
    <option value="months" {{if eq .Effort "months"}}selected{{end}}>months</option>
  </select><br><br>

  <input type="submit" value="Update">
    {{template "backButton"}}
</form>
`
var editForm = template.Must(template.New("editForm").Parse(editFormHTML + backButton))
		
func updateTodo(conn *pgx.Conn, id int, short string, description string, dueDate *time.Time, costOfDelay int16, effort string) error {
	tag, err := conn.Exec(
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
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
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
			if err == pgx.ErrNoRows {
				http.NotFound(w, r)
				return
			}
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
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.NotFound(w, r)
			return
		}

		todo, err := getTodo(conn, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := editForm.Execute(w, todo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
