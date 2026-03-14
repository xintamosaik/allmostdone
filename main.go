package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
 
	"strconv"
	"strings"
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
func check (err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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
		costOfDelay = int16(tmp)
	}
	effort = r.FormValue("effort")
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

	fmt.Println("Connected to database:", db)

 
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// todo endpoints
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listHandler(conn)(w, r)
		case http.MethodPost:
			createHandler(conn)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/todos/new", newHandler(conn))

	http.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/edit") {
			editHandler(conn)(w, r)
		} else if r.Method == http.MethodPost {
			updateHandler(conn)(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
