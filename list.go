package main

import (
	"context"
	"html/template"
	"net/http"

	"github.com/jackc/pgx/v5"
)

const todoList = `
<h1>Todo List</h1>
<table>
    <tr>
        <th>ID</th>
        <th>Short</th>
        <th>Description</th>
        <th>Due Date</th>
        <th>Cost of Delay</th>
        <th>Effort</th>
        <th>Created At</th>
        <th>Updated At</th>
        <th>Actions</th>
    </tr>
    {{range .}}
    <tr>
        <td>{{.ID}}</td>
        <td>{{.Short}}</td>
        <td>{{.Description}}</td>
        <td>{{if .DueDate}}{{.DueDate.Format "2006-01-02"}}{{else}}N/A{{end}}</td>
        <td>{{.CostOfDelay}}</td>
        <td>{{.Effort}}</td>
        <td>{{.CreatedAt.Format "2006-01-02 15:04:05"}}</td>
        <td>{{.UpdatedAt.Format "2006-01-02 15:04:05"}}</td>
        <td>
            <button fx-action="/todos/{{.ID}}/edit" fx-target="#output" fx-swap="innerHTML">
                Edit
            </button>
        </td>
    </tr>
    {{end}}
</table>
`

func getTodos(conn *pgx.Conn) ([]Todo, error) {
	rows, err := conn.Query(context.Background(),
		`SELECT id, short, description, due_date, cost_of_delay, effort, created_at, updated_at
         FROM todos
         ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var t Todo
		err := rows.Scan(&t.ID, &t.Short, &t.Description, &t.DueDate, &t.CostOfDelay, &t.Effort, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func listHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := getTodos(conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t, err := template.New("webpage").Parse(todoList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := t.Execute(w, todos); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
