package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func todoInputDueDateForDueDateEditor(dueDate DueDate) string {
	if dueDate == nil {
		return ""
	}
	return dueDate.Format("2006-01-02")
}

// saveDueDateInlineHandler handles inline updates to the due date field. Path: /todos/{id}/update/due-date
func (a *App) saveDueDateInlineHandler() http.HandlerFunc {
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

		dateStr := strings.TrimSpace(r.FormValue("due_date"))
		var dueDate *time.Time
		if dateStr != "" {
			dt, parseErr := time.Parse("2006-01-02", dateStr)
			if parseErr != nil {
				http.Error(w, "due_date must be in YYYY-MM-DD format", http.StatusBadRequest)
				return
			}
			dueDate = &dt
		}

		if err = updateTodoPartial(r.Context(), a.DB, id, TodoPartialUpdate{
			DueDate:    dueDate,
			DueDateSet: true,
		}); err != nil {
			if err == pgx.ErrNoRows {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		display := "N/A"
		if dueDate != nil {
			display = dueDate.Format("2006-01-02")
		}
		if err = DueDateDisplay(id, display).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// showDueDateEditorHandler swaps the due date field into input mode. Path: /todos/{id}/edit/due-date
func (a *App) showDueDateEditorHandler() http.HandlerFunc {
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

		if err = DueDateEditor(todo.ID, todoInputDueDateForDueDateEditor(todo.DueDate)).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
