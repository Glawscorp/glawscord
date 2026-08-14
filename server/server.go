package server

import (
	"encoding/json"
	"fmt"
	"github.com/glawscorp/glawscord/db"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
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
	if (utf8.RuneCountInString(username)) < 4 || (utf8.RuneCountInString(username)) > 25 {
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

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
}

func InitServer() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Post("/", createUser)
		r.Get("/", getUsers)

	})
	r.Route("/users/{ID}/messages", func(r chi.Router) {
		r.Post("/", sendUserMessage)
		r.Get("/", getUserMessages)

	})

	return r
}

func sendUserMessage(w http.ResponseWriter, r *http.Request) {
	var m db.UserMessage

	err := json.NewDecoder(r.Body).Decode(&m)

	if err != nil {
		fmt.Fprintln(w, err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err = db.CreateUserMessage(m.Sender, m.Receiver, m.Content)
	fmt.Println(err)

}

func getUserMessages(w http.ResponseWriter, r *http.Request) {
	var msgs []*db.UserMessage
	q := r.URL.Query()

	sender := q.Get("sender_id")
	receiver := q.Get("receiver_id")
	limit := q.Get("limit")
	offset := q.Get("offset")

	s, err := strconv.Atoi(sender)
	if err != nil {
		http.Error(w, fmt.Sprintf("invlaid sender_id '%s'", sender), http.StatusBadRequest)
		return
	}

	rec, err := strconv.Atoi(receiver)
	if err != nil {
		http.Error(w, fmt.Sprintf("invlaid receiver_id '%s'", receiver), http.StatusBadRequest)
		return
	}

	if limit == "" || limit == "0" {
		limit = "10"
	}
	l, err := strconv.Atoi(limit)

	if err != nil {
		http.Error(w, fmt.Sprintf("invlaid limit '%s'", limit), http.StatusBadRequest)
		return
	}
	if offset == "" {
		offset = "0"
	}

	o, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("invlaid offset '%s'", offset), http.StatusBadRequest)
		return
	}

	msgs, err = db.GetUserMessages(s, rec, l, o)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "made it to db func call", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(msgs)
	if err != nil {
		fmt.Println(err)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {

	var u db.User

	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if !validUsername(u.Username) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "username "+u.Username+" contains invalid chars")
		return

	}

	exists, err := db.GetUserByName(u.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	if exists != nil {
		fmt.Fprintln(w, u.Username+" already exists")
		w.WriteHeader(http.StatusConflict)
		return
	}

	if !validPassword(u.Password) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "password is invalid")
		return
	}

	err = db.CreateUser(u.Username, u.Password)

	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	success := "successfully created user: " + u.Username
	fmt.Fprintln(w, success)

}
