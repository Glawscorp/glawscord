package server

import (
	"fmt"
	"github.com/glawscorp/glawscord/db"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"
)

var validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

func validUsername(username string) bool {
	for _, c := range username {
		if !strings.ContainsRune(validChars, c) {
			return false
		}
	}
	if (utf8.RuneCountInString(username)) < 4 || (utf8.RuneCountInString(username)) > 17 {
		return false
	}

	return true

}

func validPassword(password string) bool {
	containsUpper := false
	containsLower := false

	for _, c := range password {
		if !strings.ContainsRune(validChars, c) {

			return false
		}

	}
	if (utf8.RuneCountInString(password)) < 8 || (utf8.RuneCountInString(password)) > 20 {
		return false
	}

	for _, c := range password {
		if unicode.IsUpper(c) {
			containsUpper = true

		}
	}

	for _, c := range password {
		if unicode.IsLower(c) {
			containsLower = true

		}
	}

	if containsUpper == false || containsLower == false {
		return false
	}
	return true

}

func getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetUsers()
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	content := strings.Join(users, "\n")
	fmt.Fprintln(w, content)
}

func InitServer() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Post("/", createUser)
		r.Get("/", getUsers)

	})

	return r
}

func createUser(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	username := q.Get("username")
	if !validUsername(username) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "username"+username+" contains invalid chars")
		return

	}

	u, err := db.GetUserByName(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	if u != "" {
		fmt.Fprintln(w, username+" already exists")
		w.WriteHeader(http.StatusConflict)
		return
	}

	password := q.Get("password")

	if !validPassword(password) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "password is invalid")
		return
	}

	err = db.CreateUser(username, password)

	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	success := "successfully created user: " + username
	fmt.Fprintln(w, success)

}
