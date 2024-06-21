package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/XanderMoroz/mongoMovies/internal/routers"
)

func main() {
	fmt.Println("MongoDB setup for Golang")
	r := routers.Router()
	fmt.Println("Server Is Getting Started...")
	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("Listening 4000 port...")
}
