package main

import (
	"Adoutchquizz/server"
)

func main() {
    server := server.NewServer()
    err := server.ListenAndServe()
	if err != nil {
        panic(err)
	}
}
