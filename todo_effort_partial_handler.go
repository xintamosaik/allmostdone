package main

import (
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

// saveEffortInlineHandler handles inline updates to the effort field. Path: /todos/{id}/update/effort
func (a *App) saveEffortInlineHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := todoIDFromRequest(r)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if err = parseInlineForm(r); err != nil {
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		effort := strings.TrimSpace(r.FormValue("effort"))
		if !isAllowedEffort(effort) {
			http.Error(w, "invalid effort", http.StatusBadRequest)
			return
		}

		if err = updateTodoPartial(r.Context(), a.DB, id, TodoPartialUpdate{
			Effort:    effort,
			EffortSet: true,
		}); err != nil {
			if err == pgx.ErrNoRows {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = EffortDisplay(id, effort).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// showEffortEditorHandler swaps the effort field into input mode. Path: /todos/{id}/edit/effort
func (a *App) showEffortEditorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := todoIDFromRequest(r)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		todo, err := getTodo(r.Context(), a.DB, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = EffortEditor(todo.ID, todo.Effort).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
