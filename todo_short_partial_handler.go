package main

import (
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

// saveShortInlineHandler handles inline updates to the short field. Path: /todos/{id}/update/short
func (a *App) saveShortInlineHandler() http.HandlerFunc {
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

		short := strings.TrimSpace(r.FormValue("short"))
		if short == "" {
			http.Error(w, "short is required", http.StatusBadRequest)
			return
		}

		if err = updateTodoPartial(r.Context(), a.DB, id, TodoPartialUpdate{
			Short:    short,
			ShortSet: true,
		}); err != nil {
			if err == pgx.ErrNoRows {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = ShortDisplay(id, short).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// showShortEditorHandler swaps the short field into input mode. Path: /todos/{id}/edit/short
func (a *App) showShortEditorHandler() http.HandlerFunc {
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

		if err = ShortEditor(todo.ID, todo.Short).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
