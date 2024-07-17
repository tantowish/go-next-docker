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
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// main function
func main() {
	// connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create table if not exists
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, email TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// create router
	router := mux.NewRouter()
	router.HandleFunc("/api/go/users", getUsers(db)).Methods("GET")
	router.HandleFunc("/api/go/users", createUser(db)).Methods("POST")
	router.HandleFunc("/api/go/users/{id}", getUser(db)).Methods("GET")
	router.HandleFunc("/api/go/users/{id}", updateUser(db)).Methods("PUT")
	router.HandleFunc("/api/go/users/{id}", deleteUser(db)).Methods("DELETE")

	// wrap router with CORS and JSON content type middleware
	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))

	// start server
	log.Fatal(http.ListenAndServe(":8000", enhancedRouter))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Check if the request is for CORS preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set JSON Content-Type
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// get all users
func getUsers(db *sql.DB) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			http.Error(response, "Failed to Fetch", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(response).Encode(map[string]string{"message": err.Error()})
				return
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(map[string]string{"message": err.Error()})
			return
		}

		resp := map[string]interface{}{
			"message": "Success Get List Users",
			"data": users,
		}
		json.NewEncoder(response).Encode(resp)
	}
}

// get user by id
func getUser(db *sql.DB) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			json.NewEncoder(response).Encode(map[string]string{"message": "User not found"})
			return
		}

		resp := map[string]interface{}{
			"message": "Success Get User",
			"data":    u,
		}
		json.NewEncoder(response).Encode(resp)
	}
}

// create user
func createUser(db *sql.DB) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var u User
		if err := json.NewDecoder(request.Body).Decode(&u); err != nil {
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(map[string]string{"message": err.Error()})
			return
		}

		err := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(map[string]string{"message": err.Error()})
			return
		}

		resp := map[string]interface{}{
			"message": "Success Create User",
			"data":    u,
		}

		json.NewEncoder(response).Encode(resp)
	}
}

// update user
func updateUser(db *sql.DB) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var u User
		if err := json.NewDecoder(request.Body).Decode(&u); err != nil {
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(map[string]string{"message": err.Error()})
			return
		}

		vars := mux.Vars(request)
		id := vars["id"]

		_, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", u.Name, u.Email, id)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(map[string]string{"message": err.Error()})
			return
		}

		var updatedUser User
		err = db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(map[string]string{"message": err.Error()})
			return
		}

		resp := map[string]interface{}{
			"message": "Success Update User",
			"data":    updatedUser,
		}

		json.NewEncoder(response).Encode(resp)	
	}
}

// delete user
func deleteUser(db *sql.DB) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			response.WriteHeader(http.StatusNotFound)
			json.NewEncoder(response).Encode(map[string]string{"message": "User not found"})
			return
		}

		_, err = db.Exec("DELETE FROM users WHERE id = $1", id)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]interface{}{
			"message": "Success Delete User",
			"user":    u,
		}
		json.NewEncoder(response).Encode(resp)
	}
}
