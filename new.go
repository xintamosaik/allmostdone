package main

import (
	"net/http"

	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

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
		if err := NewForm().Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
