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

type ID = int
type Short = string
type Description = string
type DueDate = *time.Time
type CostOfDelay = int16
type Effort = string
type CreatedAt = time.Time
type UpdatedAt = time.Time
type Todo struct {
	ID          ID
	Short       Short
	Description Description
	DueDate     DueDate
	CostOfDelay CostOfDelay
	Effort      Effort
	CreatedAt   CreatedAt
	UpdatedAt   UpdatedAt
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
	http.HandleFunc("POST /todos/{id}/edit/short", app.showShortEditorHandler())
	http.HandleFunc("POST /todos/{id}/update/short", app.saveShortInlineHandler())
	http.HandleFunc("POST /todos/{id}/edit/description", app.showDescriptionEditorHandler())
	http.HandleFunc("POST /todos/{id}/update/description", app.saveDescriptionInlineHandler())
	http.HandleFunc("POST /todos/{id}/edit/due-date", app.showDueDateEditorHandler())
	http.HandleFunc("POST /todos/{id}/update/due-date", app.saveDueDateInlineHandler())
	http.HandleFunc("POST /todos/{id}/edit/cost-of-delay", app.showCostOfDelayEditorHandler())
	http.HandleFunc("POST /todos/{id}/update/cost-of-delay", app.saveCostOfDelayInlineHandler())
	http.HandleFunc("POST /todos/{id}/edit/effort", app.showEffortEditorHandler())
	http.HandleFunc("POST /todos/{id}/update/effort", app.saveEffortInlineHandler())
	http.HandleFunc("POST /todos/{id}/delete", app.deleteHandler())

	fmt.Println("Listening on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("http server failed: %v", err)
	}
}
