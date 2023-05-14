package routes

import "github.com/gorilla/mux"

func RouteIndex(r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()

	CategoryRoutes(api)
	ArticleRoutes(api)
}
