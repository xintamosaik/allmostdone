package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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

func main() {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		if _, writeErr := fmt.Fprintln(os.Stderr, "Unable to connect to database: DATABASE_URL is not set"); writeErr != nil {
			log.Printf("unable to write database connection error to stderr: %v", writeErr)
		}
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), databaseURL)
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
