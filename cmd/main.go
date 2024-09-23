package main

import (
    "fmt"

	"Adoutchquizz/server"
)

func main() {
    server := server.NewServer()
    fmt.Printf("Server started : %s ", server.Addr);
	err := server.ListenAndServe()
	if err != nil {
        panic(err)
	}
}
