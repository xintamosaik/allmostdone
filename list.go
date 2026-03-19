package main

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
)

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

		if err := TodoList(todos).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
