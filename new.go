package main

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func createTodo(ctx context.Context, db *pgxpool.Pool, in TodoInput) (Todo, error) {
	var t Todo

	err := db.QueryRow(
		ctx,
		`INSERT INTO todos (short, description, due_date, cost_of_delay, effort)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, short, description, due_date, cost_of_delay, effort, created_at, updated_at`,
		in.Short,
		in.Description,
		in.DueDate,
		in.CostOfDelay,
		in.Effort,
	).Scan(&t.ID, &t.Short, &t.Description, &t.DueDate, &t.CostOfDelay, &t.Effort, &t.CreatedAt, &t.UpdatedAt)

	return t, err
}

func (a App) createHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, err := parseTodoForm(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = createTodo(r.Context(), a.DB, in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// after create show updated list
		a.listHandler()(w, r)
	}
}

func (_ App) newHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := NewForm().Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
