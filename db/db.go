package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

var createUsersTableQuery = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY, 
	username TEXT,
	password TEXT,
	joined_at DATETIME DEFAULT CURRENT_TIMESTAMP
)`

var createUserMessagesTableQuery = `
CREATE TABLE IF NOT EXISTS user_messages (
	id INTEGER PRIMARY KEY,
	sender INT,
	receiver INT,
	sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	content TEXT
)`

var createUserMessage = `
INSERT INTO user_messages (
	sender,
	receiver,
	content
) VALUES (
	?,
	?,
	?
)`

var getUserMessages = `
SELECT * FROM user_messages WHERE sender = ? AND receiver = ? ORDER BY sent_at LIMIT ? OFFSET ? 
`
var updateUsername = `
UPDATE users SET username = ? WHERE username = ?
`

var deleteUser = `
DELETE FROM users WHERE id = ?
`

var getUserByID = `
SELECT * FROM users WHERE id = ?
`
var dbPath string

func InitDB(path string) error {
	dbPath = path
	db := GetDB()
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Printf("error closing database: %v\n", err)
		}
	}()

	_, err := db.Exec(createUsersTableQuery)

	if err != nil {
		fmt.Printf("issue creating users table: %v\n", err)
	}

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(createUserMessagesTableQuery)

	if err != nil {
		panic(err)
	}

	return err
}

func GetDB() *sql.DB {

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic(err)
	}
	return db
}

func CreateUserMessage(sender int, receiver int, content string) error {
	db := GetDB()
	_, err := db.Exec(createUserMessage, sender, receiver, content)
	if err != nil {
		return err
	}
	return nil
}

func GetUserMessages(sender, receiver, limit, offset int) ([]*UserMessage, error) {
	db := GetDB()

	rows, err := db.Query(getUserMessages, sender, receiver, limit, offset)

	if err != nil {
		return nil, err
	}

	result := []*UserMessage{}
	for rows.Next() {
		var m UserMessage
		err := rows.Scan(&m.ID, &m.Sender, &m.Receiver, &m.SentAt, &m.Content)
		if err != nil {
			return nil, err
		}
		result = append(result, &m)
	}
	//add this check everywhere we loop over queried rows
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func GetUsers() ([]string, error) {
	db := GetDB()
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

func GetUserByName(username string) (*User, error) {
	var u User
	db := GetDB()
	q := fmt.Sprintf(`SELECT * FROM users WHERE username = '%s'`, username)
	if err := db.QueryRow(q).Scan(&u.ID, &u.Username, &u.Password, &u.JoinedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &u, nil
}

func CreateUser(username string, password string) error {

	db := GetDB()
	//nolint:errcheck
	defer db.Close()

	query := fmt.Sprintf(`INSERT INTO users (username, password) VALUES('%s','%s')`, username, password)

	_, err := db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}

func UpdateUsername(username string, new_name string) error {
	db := GetDB()

	_, err := GetUserByName(username)

	if err != nil {
		fmt.Println("no user found with the given username")
		return err
	}

	_, err = db.Exec(updateUsername, new_name, username)

	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(id int) error {

	db := GetDB()

	_, err := GetUserByID(id)

	if err != nil {
		fmt.Println("no user found with the given id")
		return err
	}

	_, err = db.Exec(deleteUser, id)

	if err != nil {
		return err
	}

	return nil

}

func GetUserByID(id int) (*User, error) {
	db := GetDB()

	var u User
	if err := db.QueryRow(getUserByID, id).Scan(&u.ID, &u.Username, &u.Password, &u.JoinedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &u, nil

}
