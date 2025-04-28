package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"github.com/hriday111/go-rest-api/internal/db"
)


type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
func main(){

	db.Connect()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, welcome to the Go REST API project!")
	})

	http.HandleFunc("/register", registerHandler)

	fmt.Println("Starting server on :8080...")
	err:= http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost{
		http.Error(w, "Only POST Methodallowed", http.StatusMethodNotAllowed)
		return
	}

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil{
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Printf("Recieved registration: Name=%s, Email=%s\n", user.Name, user.Email)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully!")
}