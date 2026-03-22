package main

import (
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

// saveDescriptionInlineHandler handles inline updates to the description field. Path: /todos/{id}/update/description
func (a *App) saveDescriptionInlineHandler() http.HandlerFunc {
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

		description := strings.TrimSpace(r.FormValue("description"))

		if err = updateTodoPartial(r.Context(), a.DB, id, TodoPartialUpdate{
			Description:    description,
			DescriptionSet: true,
		}); err != nil {
			if err == pgx.ErrNoRows {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = DescriptionDisplay(id, description).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// showDescriptionEditorHandler swaps the description field into input mode. Path: /todos/{id}/edit/description
func (a *App) showDescriptionEditorHandler() http.HandlerFunc {
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

		if err = DescriptionEditor(todo.ID, todo.Description).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
