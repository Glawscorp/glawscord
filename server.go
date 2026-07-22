package server

import (
    "net/http"
    "github.com/glawscorp/glawscord/db"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
    users, err := db.getUsers()
    if err != nil {
        fmt.Fprintln(w, err.Error())
        return
    }
    content := strings.Join(users, "\n")
    fmt.Fprintln(w, content)
}

func initServer() *chi.Mux {
    r := chi.NewRouter()

	r.Get("/api/v1/users", usersRouter)
    r.post("/api/v1/users"

    log.Println("listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func  (w http.ResponseWriter, r *http.Request) {
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

}


