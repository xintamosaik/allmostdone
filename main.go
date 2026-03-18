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

// helper to parse common form fields from a request
func parseTodoForm(r *http.Request) (short string, description string, dueDate *time.Time, costOfDelay int16, effort string, err error) {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		err = r.ParseMultipartForm(1 << 20)
	} else {
		err = r.ParseForm()
	}
	if err != nil {
		return
	}

	short = r.FormValue("short")
	description = r.FormValue("description")

	dateStr := r.FormValue("due_date")
	if dateStr != "" {
		var dt time.Time
		dt, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			err = fmt.Errorf("due_date must be in YYYY-MM-DD format")
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
		if _, writeErr := fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err); writeErr != nil {
			log.Printf("unable to write database connection error to stderr: %v", writeErr)
		}
		os.Exit(1)
	}
	defer func() {
		if closeErr := conn.Close(context.Background()); closeErr != nil {
			log.Printf("unable to close database connection: %v", closeErr)
		}
	}()

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
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("http server failed: %v", err)
	}
}
