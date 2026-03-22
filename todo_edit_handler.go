package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func parseInlineForm(r *http.Request) error {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return r.ParseMultipartForm(1 << 20)
	}
	return r.ParseForm()
}

func todoIDFromRequest(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
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

func updateTodoShort(ctx context.Context, db *pgxpool.Pool, id int, short string) error {
	tag, err := db.Exec(
		ctx,
		`UPDATE todos
	         SET short=$1,
	             updated_at=now()
	         WHERE id=$2`,
		short,
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

func updateTodoDescription(ctx context.Context, db *pgxpool.Pool, id int, description string) error {
	tag, err := db.Exec(
		ctx,
		`UPDATE todos
	         SET description=$1,
	             updated_at=now()
	         WHERE id=$2`,
		description,
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

		if err = updateTodoShort(r.Context(), a.DB, id, short); err != nil {
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

		if err = updateTodoDescription(r.Context(), a.DB, id, description); err != nil {
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
