package main

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"net/http"
)

func dbHelloWorld() {
	db, err := sql.Open("sqlite", "glawscord.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT, password TEXT)`)
	if err != nil {
		panic(err)
	}
}

func serverHelloWorld() {
	http.HandleFunc("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})
	http.ListenAndServe(":8080", nil)
}

func main() {
	serverHelloWorld()
}
