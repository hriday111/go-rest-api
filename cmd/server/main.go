package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"database/sql"

	"github.com/hriday111/go-rest-api/internal/db"
	"github.com/hriday111/go-rest-api/internal/models"
	"github.com/hriday111/go-rest-api/internal/utils"
)

func main() {

	db.Connect()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, welcome to the Go REST API project!")
	})

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	fmt.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST Methodallowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var existingUserID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE email = $1", user.Email).Scan(&existingUserID)
	if err != sql.ErrNoRows {
		// If err == nil (found user) or another error happened
		if err == nil {
			// User exists, now check the password
			var storedPassword string
			err = db.DB.QueryRow("SELECT password FROM users WHERE email = $1", user.Email).Scan(&storedPassword)
			if err != nil {
				http.Error(w, "Error retrieving user password: "+err.Error(), http.StatusInternalServerError)
				return
			}
			// Compare the hashed password with the provided password
			if !utils.CheckPasswordHash(user.Password, storedPassword) {
				http.Error(w, "Invalid password", http.StatusUnauthorized)
				return
			}

			// Password is correct, login successful
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Login successful! User ID:", existingUserID)
			return

		}

		http.Error(w, "User does not exist "+err.Error(), http.StatusInternalServerError)

		return
	}

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST Methodallowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var existingUserID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE email = $1", user.Email).Scan(&existingUserID)
	if err != sql.ErrNoRows {
		// If err == nil (found user) or another error happened
		if err == nil {
			http.Error(w, "Email already registered", http.StatusBadRequest)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword
	// Save the user to DB

	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	err = db.DB.QueryRow(query, user.Name, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Error inserting user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Recieved registration: Name=%s, Email=%s\n", user.Name, user.Email)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully! ID:", user.ID)
}
