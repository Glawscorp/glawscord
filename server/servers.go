package server

import (
	"encoding/json"
	"fmt"
	"github.com/glawscorp/glawscord/db"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func serverRoutes(r chi.Router) {
	r.Route("/servers", func(r chi.Router) {
		//user ID
		r.Post("/", createServer)
		r.Post("/{ID}", sendServerMessage)
	})

}

func createServer(w http.ResponseWriter, r *http.Request) {
	//still needs validation for server names

	var s db.Server

	err := json.NewDecoder(r.Body).Decode(&s)
	fmt.Printf("attempting to create server: %v", s)

	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err = db.CreateServer(s.Name, s.OwnerID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	success := "successfully created server: " + s.Name
	_, _ = fmt.Fprintln(w, success)

}

func sendServerMessage(w http.ResponseWriter, r *http.Request) {

	var s db.ServerMessage

	err := json.NewDecoder(r.Body).Decode(&s)

	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err = db.CreateServerMessage(s.ServerID, s.Sender, s.Content)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
}
