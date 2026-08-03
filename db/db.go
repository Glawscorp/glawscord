package db

import (
	"database/sql"
	"fmt"
	"github.com/glawscorp/glawscord/structs"
	_ "modernc.org/sqlite"
)

func init() {
	db := getDB()
	defer db.Close()

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, password TEXT)`)
	if err != nil {
		panic(err)
	}

}

func getDB() *sql.DB {

	db, err := sql.Open("sqlite", "glawscord.db")
	if err != nil {
		panic(err)
	}
	return db
}

func GetUsers() ([]string, error) {
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

func GetUserByName(username string) (*structs.User, error) {
	db := getDB()
	q := fmt.Sprintf(`SELECT * FROM users WHERE username = '%s'`, username)
	r, err := db.Query(q)
	var u structs.User

	if err != nil {
		return nil, err
	}

	for r.Next() {
		err = r.Scan(&u)
		if err != nil {
			return nil, err
		}
	}
	return &u, nil
}

func CreateUser(username string, password string) error {

	db := getDB()
	defer db.Close()

	query := fmt.Sprintf(`INSERT INTO users (username, password) VALUES('%s','%s')`, username, password)

	_, err := db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}
