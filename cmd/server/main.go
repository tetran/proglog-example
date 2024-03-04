package main

import (
	"log"

	"github.com/tetran/proglog-example/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
