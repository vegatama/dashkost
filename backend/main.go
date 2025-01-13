package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	email string `json:"email"`
}

func main() {
	// connect database
	db err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, email TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// create router
	r := mux.NewRouter()
	r.HandleFunc("/api/users", getUsers(db)).Methods("GET")
	r.HandleFunc("/api/users", createUser(db)).Methods("POST")
	r.HandleFunc("/api/users/{id}", getUser(db)).Methods("GET")
	r.HandleFunc("/api/users/{id}", updateUser(db)).Methods("PUT")
	r.HandleFunc("/api/users/{id}", deleteUser(db)).Methods("DELETE")

	// wrap router
	enchancedRouter := enableCORS(jsonMiddleware(r))

	// start server
	log.Fatal(http.ListenAndServe(":"+os.Getenv("8000"), enchancedRouter))

	func enableCORS (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	func jsonMiddleware(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}

	


