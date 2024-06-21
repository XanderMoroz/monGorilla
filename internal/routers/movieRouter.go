package routers

import (
	"github.com/gorilla/mux"

	"github.com/XanderMoroz/mongoMovies/internal/controllers"
)

// Captial means exporting the method
func Router() *mux.Router {

	router := mux.NewRouter()
	router.HandleFunc("/api/movies", controllers.GetAlIMovies).Methods("GET")
	router.HandleFunc("/api/movie", controllers.CreateMovie).Methods("POST")
	router.HandleFunc("/api/movie/{id}", controllers.MarkAsWatched).Methods("PUT")
	router.HandleFunc("/api/movie/{id}", controllers.DeleteMovie).Methods("DELETE")
	router.HandleFunc("/api/delallmovies", controllers.DeleteAllMovies).Methods("DELETE")
	return router
}
