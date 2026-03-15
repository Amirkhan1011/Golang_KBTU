package main

import (
	"log"
	"net/http"

	"practice5/db"
	"practice5/handler"
	"practice5/repository"
)

func main() {
	database, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	userRepo := repository.NewUserRepository(database)
	userHandler := handler.NewUserHandler(userRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", userHandler.GetUsers)
	mux.HandleFunc("/common-friends", userHandler.GetCommonFriends)

	addr := ":8080"
	log.Printf("Server is running on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("http server failed: %v", err)
	}
}

