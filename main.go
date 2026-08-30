package main

import (
	"flag"
	"fmt"
	"github.com/glawscorp/glawscord/db"
	"github.com/glawscorp/glawscord/server"
	"net/http"
)

func main() {
	dbPath := flag.String("db-path", "./glawscord.db", "the path to the database file")

	flag.Parse()

	err := db.InitDB(*dbPath)
	if err != nil {
		fmt.Println(err)
	}

	r := server.InitServer()
	fmt.Println("listening on port 8080")
	http.ListenAndServe(":8080", r)

}
