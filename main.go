package main

import (
	"context"
	"fmt"
	"net/http"

	"os"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5"
)

type Todo struct {
	ID          int
	Short       string
	Description string
}

func createTodo(conn *pgx.Conn, short string, description string) (Todo, error) {
	var t Todo

	err := conn.QueryRow(
		context.Background(),
		`INSERT INTO todos (short, description)
		 VALUES ($1, $2)
		 RETURNING id, short, description`,
		short,
		description,
	).Scan(&t.ID, &t.Short, &t.Description)

	return t, err
}
func getTodos(conn *pgx.Conn) ([]Todo, error) {
	rows, err := conn.Query(context.Background(),
		`SELECT id, short, description
		 FROM todos
		 ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var t Todo
		err := rows.Scan(&t.ID, &t.Short, &t.Description)
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
		`SELECT id, short, description
		 FROM todos
		 WHERE id=$1`,
		id,
	).Scan(&t.ID, &t.Short, &t.Description)

	return t, err
}
func updateTodo(conn *pgx.Conn, id int, short string, description string) error {
	_, err := conn.Exec(
		context.Background(),
		`UPDATE todos
		 SET short=$1,
		     description=$2
		 WHERE id=$3`,
		short,
		description,
		id,
	)

	return err
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
	
	list, _ := getTodos(conn)
	fmt.Println(list)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	component := hello("John")
	http.Handle("/test", templ.Handler(component))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
