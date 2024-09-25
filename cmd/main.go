package main

import (
    "slices"
    "os"
    
    "Adoutchquizz/database"
	"Adoutchquizz/server"
)

func main() {

    if slices.Contains(os.Args, "--migrate") {
        db := database.New()
        err := db.Migrate()
        if err != nil {
            panic(err)
        }
        return
    }
    if slices.Contains(os.Args, "--drop") {
        db := database.New()
        err := db.Drop()
        if err != nil {
            panic(err)
        }
        return
    }


    server := server.NewServer()
    err := server.ListenAndServe()
    if err != nil {
        panic(err)
	}
}
