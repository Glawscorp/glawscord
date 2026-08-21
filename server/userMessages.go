package server

import (
	"encoding/json"
	"fmt"
	"github.com/glawscorp/glawscord/db"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func userMessageRoutes(r chi.Router) {

	r.Route("/users/{ID}/messages", func(r chi.Router) {
		r.Post("/", sendUserMessage)
		r.Get("/", getUserMessages)

	})

}

func sendUserMessage(w http.ResponseWriter, r *http.Request) {
	var m db.UserMessage

	err := json.NewDecoder(r.Body).Decode(&m)

	if err != nil {
		fmt.Fprintln(w, err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if m.Content == "" {
		http.Error(w, "cannot send empty message", http.StatusBadRequest)
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
