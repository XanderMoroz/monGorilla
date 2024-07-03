package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/XanderMoroz/mongoMovies/configs"
	"github.com/XanderMoroz/mongoMovies/internal/routers"

	_ "github.com/XanderMoroz/mongoMovies/docs"
)

// @title MonGorilla Project
// @version 1.0
// @description This is a sample server on Gorrilla Mux + MongoDB.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:4000/
// // @BasePath /v1
func main() {

	configs.ConnectDB()

	r := routers.CommonRouter()

	fmt.Println("=========== MonGorilla Start ============")
	fmt.Println("== URL: http://127.0.0.1:4000/swagger/ ==")
	log.Fatal(http.ListenAndServe(":4000", r))
}
