package server

import (
	"encoding/json"
	"fmt"
	"github.com/glawscorp/glawscord/db"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func userRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", createUser)
		r.Get("/", getUsers)
		r.Patch("/", updateUsername)
		r.Delete("/{ID}", deleteUser)
	})

}

func getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetUsers()
	if err != nil {
		fmt.Printf("failed to retrieve users: %v\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)

	if err != nil {
		fmt.Printf("failed to send payload to client: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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
		http.Error(w, "username is invalid", http.StatusBadRequest)
		return

	}

	exists, err := db.GetUserByName(u.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	if exists != nil {
		http.Error(w, "username already exists", http.StatusConflict)
		return
	}

	if !validPassword(u.Password) {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	err = db.CreateUser(u.Username, u.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	success := "successfully created user: " + u.Username
	_, _ = fmt.Fprintln(w, success)

}

func updateUsername(w http.ResponseWriter, r *http.Request) {
	var up UpdateUsername
	err := json.NewDecoder(r.Body).Decode(&up)

	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	username := up.Username
	new_name := up.NewName

	if !validUsername(new_name) {
		http.Error(w, "invalid username", http.StatusBadRequest)
		return
	}

	exists, err := db.GetUserByName(new_name)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	if exists != nil {
		http.Error(w, "user with that username already exists", http.StatusConflict)
		return
	}

	err = db.UpdateUsername(username, new_name)
	if err != nil {
		http.Error(w, "unable to update username", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	_, err = fmt.Fprint(w, "username updated")
	if err != nil {
		fmt.Printf("failed to send response payload: %v\n", err)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {

	user_id, err := strconv.Atoi(chi.URLParam(r, "ID"))

	if err != nil {
		http.Error(w, "couldn't get user id from url", http.StatusBadRequest)
		return
	}

	err = db.DeleteUser(user_id)
	if err != nil {
		http.Error(w, "issue deleting user", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	_, _ = fmt.Fprintln(w, "user deleted")

}
