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
		fmt.Fprintln(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)

	if err != nil {
		fmt.Fprintln(w, err)
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

func updateUsername(w http.ResponseWriter, r *http.Request) {
	var up UpdateUsername
	err := json.NewDecoder(r.Body).Decode(&up)

	if err != nil {
		fmt.Fprintln(w, err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	username := up.Username
	new_name := up.NewName

	if !validUsername(new_name) {
		fmt.Fprintln(w, "new username: "+new_name+" is invalid")
		return
	}

	exists, err := db.GetUserByName(new_name)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	if exists != nil {
		fmt.Fprintln(w, "user with username: "+new_name+" already exists")
		w.WriteHeader(http.StatusConflict)
		return
	}

	err = db.UpdateUsername(username, new_name)

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprint(w, "username updated")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {

	user_id, err := strconv.Atoi(chi.URLParam(r, "ID"))

	if err != nil {
		fmt.Fprintln(w, "couldn't get user id from url")
		return
	}

	err = db.DeleteUser(user_id)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	fmt.Fprintln(w, "user deleted")
	return

}
