package main

import (
	"log"
	"net/http"

	"1/internal/handlers"
	"1/internal/middleware"
)

func main() {
	taskHandler := handlers.NewTaskHandler()

	var handler http.Handler = taskHandler
	handler = middleware.Logging(handler)
	handler = middleware.Auth(handler)

	http.Handle("/tasks", handler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
