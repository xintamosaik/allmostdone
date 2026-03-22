package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

// saveCostOfDelayInlineHandler handles inline updates to the cost of delay field. Path: /todos/{id}/update/cost-of-delay
func (a *App) saveCostOfDelayInlineHandler() http.HandlerFunc {
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

		costOfDelayStr := strings.TrimSpace(r.FormValue("cost_of_delay"))
		if costOfDelayStr == "" {
			http.Error(w, "cost_of_delay is required", http.StatusBadRequest)
			return
		}

		tmp, parseErr := strconv.Atoi(costOfDelayStr)
		if parseErr != nil {
			http.Error(w, "cost_of_delay must be a number", http.StatusBadRequest)
			return
		}
		if tmp < -2 || tmp > 2 {
			http.Error(w, "cost_of_delay must be between -2 and 2", http.StatusBadRequest)
			return
		}

		costOfDelay := int16(tmp)
		if err = updateTodoPartial(r.Context(), a.DB, id, TodoPartialUpdate{
			CostOfDelay:    costOfDelay,
			CostOfDelaySet: true,
		}); err != nil {
			if err == pgx.ErrNoRows {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = CostOfDelayDisplay(id, costOfDelay).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// showCostOfDelayEditorHandler swaps the cost of delay field into input mode. Path: /todos/{id}/edit/cost-of-delay
func (a *App) showCostOfDelayEditorHandler() http.HandlerFunc {
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

		if err = CostOfDelayEditor(todo.ID, todo.CostOfDelay).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
