package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/XanderMoroz/mongoMovies/internal/controllers"
)

// Captial means exporting the method
func CommonRouter() *mux.Router {

	router := mux.NewRouter()

	// Swagger routes
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		//The url pointing to API definition
		// httpSwagger.URL("http://localhost:4000/swagger/docs.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	// User routes
	router.HandleFunc("/api/users", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", controllers.GetUserByID).Methods("GET")

	router.HandleFunc("/api/movies", controllers.GetAlIMovies).Methods("GET")
	router.HandleFunc("/api/movie/{id}", controllers.MarkAsWatched).Methods("PUT")
	router.HandleFunc("/api/movie/{id}", controllers.DeleteMovie).Methods("DELETE")
	router.HandleFunc("/api/deleteallmovies", controllers.DeleteAllMovies).Methods("DELETE")

	return router
}
