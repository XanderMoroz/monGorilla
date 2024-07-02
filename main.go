package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/XanderMoroz/mongoMovies/docs"
	"github.com/XanderMoroz/mongoMovies/internal/routers"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {
	fmt.Println("MongoDB setup for Golang")
	r := routers.Router()
	fmt.Println("Server Is Getting Started...")
	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("Listening 4000 port...")
}
