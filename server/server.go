package server

import (
	"fmt"
	"github.com/glawscorp/glawscord/db"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

var validUsernameChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

func validUsername(username string) bool {
	for _, c := range username {
		if !strings.ContainsRune(validUsernameChars, c) {
			return false
		}
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
		fmt.Fprintln(w, username+"already exists")
		w.WriteHeader(http.StatusConflict)
		return
	}

	password := q.Get("password")
	err = db.CreateUser(username, password)

	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	success := "successfully created user: " + username
	fmt.Fprintln(w, success)

}
