package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

type App struct {
	DB *pgxpool.Pool
}

func main() {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		if _, writeErr := fmt.Fprintln(os.Stderr, "Unable to connect to database: DATABASE_URL is not set"); writeErr != nil {
			log.Printf("unable to write database connection error to stderr: %v", writeErr)
		}
		os.Exit(1)
	}

	dbPool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		if _, writeErr := fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err); writeErr != nil {
			log.Printf("unable to write database connection error to stderr: %v", writeErr)
		}
		os.Exit(1)
	}
	defer dbPool.Close()

	app := &App{DB: dbPool}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("GET /todos/list", app.listHandler())
	http.HandleFunc("GET /todos/new", app.newHandler())
	http.HandleFunc("POST /todos/create", app.createHandler())
	http.HandleFunc("GET /todos/{id}/edit", app.editHandler())
	http.HandleFunc("POST /todos/{id}/update", app.updateHandler())

	fmt.Println("Listening on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("http server failed: %v", err)
	}
}
