package main

import (
	"fmt"
	"github.com/glawscorp/glawscord/db"
	"github.com/glawscorp/glawscord/server"
	"net/http"
)

func main() {
	err := db.InitDB("~/Workspace/glawscord/glawscord.db")
	if err != nil {
		fmt.Println(err)
	}

	r := server.InitServer()
	fmt.Println("listening on port 8080")
	http.ListenAndServe(":8080", r)

}
