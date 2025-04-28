package db
import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {

	dbHost := getEnv("DB_HOST", "localhost")
	dbPost := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "admin")
	dbPassword := getEnv("DB_PASSWORD", "adminpassword")
	dbName := getEnv("DB_NAME", "userdb")

	connStr := fmt.Sprintf("Host=%s Port=%s User=%s Password=%s DBName=%s sslmode=disable",	dbHost, dbPost, dbUser, dbPassword, dbName)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	fmt.Println("Connected to the PostrgreSQL database successfully!")
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}