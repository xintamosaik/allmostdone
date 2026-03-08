package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	component := hello("John")
	
	http.Handle("/test", templ.Handler(component))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}