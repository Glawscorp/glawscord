package main

import (
	"database/sql"
	"fmt"
    "log"
	_ "modernc.org/sqlite"
	"net/http"
	"strings"
    // "time"
)



func main() {
    log.SetFlags(log.LstdFlags|log.Lshortfile)
	initDB()
	serverHelloWorld()
}
