package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	os.Remove("./foo.db")

	var err error
	db, err = sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	router := chi.NewRouter()
	router.Post("/user", addUser)
	router.Get("/user", getUsers)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Printf("Starting server in 8080...")

	log.Fatal(server.ListenAndServe())
}

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error in decoding: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `insert into foo (name) values (?)`
	_, err = db.Exec(query, user.Name)
	if err != nil {
		log.Printf("Error in inserting: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	query := `select id, name from foo`
	rows, err := db.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
