package routers

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// Captial means exporting the method
func SwaggerRouter() *mux.Router {

	router := mux.NewRouter()

	// Swagger routes
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		//The url pointing to API definition
		// httpSwagger.URL("http://localhost:4000/swagger/docs.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	return router
}
