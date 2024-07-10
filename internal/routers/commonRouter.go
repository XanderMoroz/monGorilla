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
	router.HandleFunc("/api/users/register", controllers.Register).Methods("POST")
	router.HandleFunc("/api/users/login", controllers.Login).Methods("POST")
	router.HandleFunc("/api/users/current_user", controllers.CurrentUser).Methods("GET")

	router.HandleFunc("/api/recipes", controllers.CreateRecipe).Methods("POST")
	router.HandleFunc("/api/recipes/{id}", controllers.GetRecipeByID).Methods("GET")
	router.HandleFunc("/api/recipes", controllers.GetAllMyRecipes).Methods("GET")
	// router.HandleFunc("/api/users/{id}", controllers.UpdateUserByID).Methods("PUT")
	// router.HandleFunc("/api/users/{id}", controllers.DeleteBook).Methods("DELETE")

	// router.HandleFunc("/api/movies", controllers.GetAlIMovies).Methods("GET")
	// router.HandleFunc("/api/movie/{id}", controllers.MarkAsWatched).Methods("PUT")
	// router.HandleFunc("/api/movie/{id}", controllers.DeleteMovie).Methods("DELETE")
	// router.HandleFunc("/api/deleteallmovies", controllers.DeleteAllMovies).Methods("DELETE")

	return router
}
