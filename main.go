package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"html/template"
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
const todoList = `
<h1>Todo List</h1>
`
func check (err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createTodo(conn *pgx.Conn, short string, description string, dueDate *time.Time, costOfDelay int16, effort string) (Todo, error) {
	var t Todo

	err := conn.QueryRow(
		context.Background(),
		`INSERT INTO todos (short, description, due_date, cost_of_delay, effort)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, short, description, due_date, cost_of_delay, effort, created_at, updated_at`,
		short,
		description,
		dueDate,
		costOfDelay,
		effort,
	).Scan(&t.ID, &t.Short, &t.Description, &t.DueDate, &t.CostOfDelay, &t.Effort, &t.CreatedAt, &t.UpdatedAt)

	return t, err
}

func getTodos(conn *pgx.Conn) ([]Todo, error) {
	rows, err := conn.Query(context.Background(),
		`SELECT id, short, description, due_date, cost_of_delay, effort, created_at, updated_at
         FROM todos
         ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var t Todo
		err := rows.Scan(&t.ID, &t.Short, &t.Description, &t.DueDate, &t.CostOfDelay, &t.Effort, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	return todos, nil
}

func getTodo(conn *pgx.Conn, id int) (Todo, error) {
	var t Todo

	err := conn.QueryRow(
		context.Background(),
		`SELECT id, short, description, due_date, cost_of_delay, effort, created_at, updated_at
         FROM todos
         WHERE id=$1`,
		id,
	).Scan(&t.ID, &t.Short, &t.Description, &t.DueDate, &t.CostOfDelay, &t.Effort, &t.CreatedAt, &t.UpdatedAt)

	return t, err
}

func updateTodo(conn *pgx.Conn, id int, short string, description string, dueDate *time.Time, costOfDelay int16, effort string) error {
	_, err := conn.Exec(
		context.Background(),
		`UPDATE todos
         SET short=$1,
             description=$2,
             due_date=$3,
             cost_of_delay=$4,
             effort=$5,
             updated_at=now()
         WHERE id=$6`,
		short,
		description,
		dueDate,
		costOfDelay,
		effort,
		id,
	)

	return err
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

// handlers ------------------------------------------------------------------

func listHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := getTodos(conn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t, err := template.New("webpage").Parse(todoList)
		check(err)
		t.Execute(w, todos)
	}
}

func newHandler(_ *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// create an empty Todo so the template never receives a nil pointer
		// (the template itself will also guard against nil values).
		empty := &Todo{}
		t, err := template.New("webpage").Parse(todoList)
		check(err)
		t.Execute(w, empty)
	}
}

func createHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		short, description, dueDate, costOfDelay, effort, err := parseTodoForm(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = createTodo(conn, short, description, dueDate, costOfDelay, effort)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// after create show updated list
		listHandler(conn)(w, r)
	}
}

func editHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// path: /todos/{id}/edit
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 3 {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			http.NotFound(w, r)
			return
		}
		todo, err := getTodo(conn, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		t, err := template.New("webpage").Parse(todoList)
		check(err)
		t.Execute(w, todo)
	}
}

func updateHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// path: /todos/{id}
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 2 {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			http.NotFound(w, r)
			return
		}
		short, description, dueDate, costOfDelay, effort, err := parseTodoForm(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = updateTodo(conn, id, short, description, dueDate, costOfDelay, effort); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		listHandler(conn)(w, r)
	}
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

	// quick verify that the database is reachable
	list, _ := getTodos(conn)
	fmt.Println(list)

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
