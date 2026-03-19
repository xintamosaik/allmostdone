package main

import (
	"fmt"
	"time"
)

func todoUpdatePath(id int) string {
	return fmt.Sprintf("/todos/%d/update", id)
}

func todoEditPath(id int) string {
	return fmt.Sprintf("/todos/%d/edit", id)
}

func todoDisplayDueDate(todo *Todo) string {
	if todo == nil || todo.DueDate == nil {
		return "N/A"
	}
	return todo.DueDate.Format("2006-01-02")
}

func todoDisplayTimestamp(ts time.Time) string {
	return ts.Format("2006-01-02 15:04:05")
}

const backButton = `
{{define "backButton"}}
<button
  type="button"
  fx-action="/todos/list"
  fx-target="#output"
  fx-swap="innerHTML">
  Back
</button>
{{end}}
`
