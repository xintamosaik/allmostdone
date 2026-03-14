package main

import (
	"html/template"
	"net/http"

	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

const newForm = `
<h1>Create Todo</h1>
<form
  action="/todos/create"
  method="post"
  fx-action="/todos/create"
  fx-method="POST"
  fx-target="#output"
  fx-swap="innerHTML">
  <label for="short">Short:</label><br>
  <input type="text" id="short" name="short"><br>

  <label for="description">Description:</label><br>
  <textarea id="description" name="description"></textarea><br>

  <label for="due_date">Due Date (YYYY-MM-DD):</label><br>
  <input type="date" id="due_date" name="due_date"><br>

  <label for="cost_of_delay">Cost of Delay:</label><br>
  <input type="number" id="cost_of_delay" name="cost_of_delay" min="-2" max="2"><br>

  <label for="effort">Effort:</label><br>
  <select id="effort" name="effort">
    <option value="mins">mins</option>
    <option value="hours" selected>hours</option>
    <option value="days">days</option>
    <option value="weeks">weeks</option>
    <option value="months">months</option>
  </select><br><br>

  <input type="submit" value="Create">
  
</form>
`

func createTodo(conn *pgx.Conn, short string, description string, dueDate *time.Time, costOfDelay int16, effort string) (Todo, error) {
	var t Todo

	err := conn.QueryRow(
		context.Background(),
		`INSERT INTO todos (short, description, due_date, cost_of_delay, effort)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, short, description, due_date, cost_of_delay, effort, created_at, updated_at`,
		short,
		description,
		dueDate,
		costOfDelay,
		effort,
	).Scan(&t.ID, &t.Short, &t.Description, &t.DueDate, &t.CostOfDelay, &t.Effort, &t.CreatedAt, &t.UpdatedAt)

	return t, err
}

func createHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		short, description, dueDate, costOfDelay, effort, err := parseTodoForm(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = createTodo(conn, short, description, dueDate, costOfDelay, effort)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// after create show updated list
		listHandler(conn)(w, r)
	}
}

func newHandler(_ *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// create an empty Todo so the template never receives a nil pointer
		// (the template itself will also guard against nil values).
		empty := &Todo{}
		t, err := template.New("webpage").Parse(newForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := t.Execute(w, empty); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
