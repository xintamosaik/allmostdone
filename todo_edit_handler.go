package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoPartialUpdate struct {
	Short          string
	ShortSet       bool
	Description    string
	DescriptionSet bool
	DueDate        *time.Time
	DueDateSet     bool
}

func parseInlineForm(r *http.Request) error {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return r.ParseMultipartForm(1 << 20)
	}
	return r.ParseForm()
}

func todoIDFromRequest(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func updateTodoPartial(ctx context.Context, db *pgxpool.Pool, id int, patch TodoPartialUpdate) error {
	tag, err := db.Exec(
		ctx,
		`UPDATE todos
         SET short = CASE WHEN $1 THEN $2 ELSE short END,
             description = CASE WHEN $3 THEN $4 ELSE description END,
             due_date = CASE WHEN $5 THEN $6 ELSE due_date END,
             updated_at=now()
         WHERE id=$7`,
		patch.ShortSet,
		patch.Short,
		patch.DescriptionSet,
		patch.Description,
		patch.DueDateSet,
		patch.DueDate,
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

func updateTodo(ctx context.Context, db *pgxpool.Pool, id int, in TodoInput) error {
	tag, err := db.Exec(
		ctx,
		`UPDATE todos
         SET short=$1,
             description=$2,
             due_date=$3,
             cost_of_delay=$4,
             effort=$5,
             updated_at=now()
         WHERE id=$6`,
		in.Short,
		in.Description,
		in.DueDate,
		in.CostOfDelay,
		in.Effort,
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

// updateHandler handles updates to an existing item. Path: /todos/{id}/update
func (a *App) updateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := todoIDFromRequest(r)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		in, err := parseTodoForm(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			a.renderEditTodoForm(w, r, id, TodoFormData{
				Input:          in,
				Error:          err.Error(),
				DueDateRaw:     r.FormValue("due_date"),
				CostOfDelayRaw: r.FormValue("cost_of_delay"),
			})
			return
		}

		if err = updateTodo(r.Context(), a.DB, id, in); err != nil {
			if err == pgx.ErrNoRows {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		a.renderTodoList(w, r)
	}
}

// updateShortHandler handles inline updates to the short field. Path: /todos/{id}/update/short
func (a *App) updateShortHandler() http.HandlerFunc {
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

		if err = EditShort(id, short).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// updateDescriptionHandler handles inline updates to the description field. Path: /todos/{id}/update/description
func (a *App) updateDescriptionHandler() http.HandlerFunc {
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

		if err = EditDescription(id, description).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// updateDueDateHandler handles inline updates to the due date field. Path: /todos/{id}/update/due-date
func (a *App) updateDueDateHandler() http.HandlerFunc {
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
		if err = EditDueDate(id, display).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// editShortHandler swaps the short field into input mode. Path: /todos/{id}/edit/short
func (a *App) editShortHandler() http.HandlerFunc {
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

		if err = InputShort(todo.ID, todo.Short).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// editDescriptionHandler swaps the description field into input mode. Path: /todos/{id}/edit/description
func (a *App) editDescriptionHandler() http.HandlerFunc {
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

		if err = InputDescription(todo.ID, todo.Description).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// editDueDateHandler swaps the due date field into input mode. Path: /todos/{id}/edit/due-date
func (a *App) editDueDateHandler() http.HandlerFunc {
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

		if err = InputDueDate(todo.ID, todoInputDueDateForCell(todo.DueDate)).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getTodo(ctx context.Context, db *pgxpool.Pool, id int) (Todo, error) {
	var t Todo

	err := db.QueryRow(
		ctx,
		`SELECT id, short, description, due_date, cost_of_delay, effort, created_at, updated_at
         FROM todos
         WHERE id=$1`,
		id,
	).Scan(&t.ID, &t.Short, &t.Description, &t.DueDate, &t.CostOfDelay, &t.Effort, &t.CreatedAt, &t.UpdatedAt)

	return t, err
}

// editHandler serves the form to edit an existing todo item. Path: /todos/{id}/edit
func (a *App) editHandler() http.HandlerFunc {
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

		a.renderEditTodoForm(w, r, todo.ID, TodoFormData{
			Input: TodoInput{
				Short:       todo.Short,
				Description: todo.Description,
				DueDate:     todo.DueDate,
				CostOfDelay: todo.CostOfDelay,
				Effort:      todo.Effort,
			},
			DueDateRaw:     todoInputDueDateValue(TodoInput{DueDate: todo.DueDate}),
			CostOfDelayRaw: todoInputCostOfDelayValue(TodoInput{CostOfDelay: todo.CostOfDelay}),
		})
	}
}

func (a *App) renderEditTodoForm(w http.ResponseWriter, r *http.Request, todoID int, data TodoFormData) {
	if err := EditTodoForm(todoID, data).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
