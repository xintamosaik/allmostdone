package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"strconv"

	"time"

	"github.com/jackc/pgx/v5"
)

type Todo struct {
	ID          int
	Short       string
	Description string
	DueDate     *time.Time
	CostOfDelay int16
	Effort      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

const html_test = `
<h1>Todo List</h1>
`

// helper to parse common form fields from a request
func parseTodoForm(r *http.Request) (short string, description string, dueDate *time.Time, costOfDelay int16, effort string, err error) {
	if err = r.ParseForm(); err != nil {
		return
	}
	short = r.FormValue("short")
	description = r.FormValue("description")
	dateStr := r.FormValue("due_date")
	if dateStr != "" {
		var dt time.Time
		dt, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return
		}
		dueDate = &dt
	}
	if codStr := r.FormValue("cost_of_delay"); codStr != "" {
		var tmp int
		tmp, err = strconv.Atoi(codStr)
		if err != nil {
			return
		}
		if tmp < -2 || tmp > 2 {
			err = fmt.Errorf("cost_of_delay must be between -2 and 2")
			return
		}
		costOfDelay = int16(tmp)
	}
	effort = r.FormValue("effort")
	switch effort {
	case "mins", "hours", "days", "weeks", "months":
		// ok
	case "":
		effort = "hours"
	default:
		err = fmt.Errorf("invalid effort")
		return
	}
	return
}

func main() {

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var db string
	err = conn.QueryRow(context.Background(), "select current_database()").Scan(&db)
	if err != nil {
		panic(err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("GET /todos/list", listHandler(conn))
	http.HandleFunc("GET /todos/new", newHandler(conn))
	http.HandleFunc("POST /todos/create", createHandler(conn))
	http.HandleFunc("GET /todos/{id}/edit", editHandler(conn))
	http.HandleFunc("POST /todos/{id}/update", updateHandler(conn))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
