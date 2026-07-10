package main

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"net/http"
	"strings"
)

func getDB() *sql.DB {

	db, err := sql.Open("sqlite", "glawscord.db")
	if err != nil {
		panic(err)
	}
	return db
}
func initDB() {
	db := getDB()
	defer db.Close()

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, password TEXT)`)
	if err != nil {
		panic(err)
	}
}

func getUsers() ([]string, error) {
	db := getDB()
	r, err := db.Query(`SELECT username FROM users`)
	if err != nil {
		return nil, err
	}
	result := []string{}
	for r.Next() {
		var u string
		err = r.Scan(&u)
		if err != nil {
			return nil, err
		}
		result = append(result, u)
	}
	return result, nil
}

func serverHelloWorld() {
	http.HandleFunc("/api/v1/hello", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			users, err := getUsers()
			if err != nil {
				fmt.Fprintln(w, err.Error())
				return
			}
			content := strings.Join(users, "\n")
			fmt.Fprintln(w, content)
		} else if r.Method == http.MethodPost {
			q := r.URL.Query()
			username := q.Get("username")
			password := q.Get("password")
			err := createUser(username, password)

			if err != nil {
				fmt.Fprintln(w, err.Error())

			}
			success := "successfully created user: " + username
			fmt.Fprintln(w, success)
		}

	})
	http.ListenAndServe(":8080", nil)
}

func createUser(username string, password string) error {

	db := getDB()
	defer db.Close()

	query := fmt.Sprintf(`INSERT INTO users (username, password) VALUES('%s','%s')`, username, password)
	_, err := db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func main() {
	initDB()
	serverHelloWorld()
}
