package main

import (
	"fmt"
	"github.com/glawscorp/glawscord/server"
	"net/http"
)

func main() {
	r := server.InitServer()
	fmt.Println("listening on port 8080")
	http.ListenAndServe(":8080", r)

}
