package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

/* ---------- MODEL ---------- */
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

/* ---------- DATABASE ---------- */
var db *sql.DB

func connectDB() {
	connStr := "postgres://postgres:postgres@localhost:5432/crud_db?sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("PostgreSQL Connected")
}

/* ---------- CREATE USER ---------- */
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User

	json.NewDecoder(r.Body).Decode(&user)

	_, err := db.Exec(
		"INSERT INTO users (name, age) VALUES ($1, $2)",
		user.Name,
		user.Age,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created",
	})
}

/* ---------- GET USERS ---------- */
func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Age)
		users = append(users, u)
	}

	json.NewEncoder(w).Encode(users)
}

/* ---------- ROUTER ---------- */
func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
		createUser(w, r)
		return
	}

	if r.Method == http.MethodGet {
		getUsers(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}


/* ---------- MAIN ---------- */
func main() {
	connectDB()

	http.HandleFunc("/users", usersHandler)

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
