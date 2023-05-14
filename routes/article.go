package routes

import (
	"go-api-article/controllers/articlecontroller"

	"github.com/gorilla/mux"
)

func ArticleRoutes(r *mux.Router) {
	router := r.PathPrefix("/articles").Subrouter()

	router.HandleFunc("", articlecontroller.Index).Methods("GET")
	router.HandleFunc("", articlecontroller.Create).Methods("POST")
	router.HandleFunc("/{id}/detail", articlecontroller.Detail).Methods("GET")
	router.HandleFunc("/{id}/update", articlecontroller.Update).Methods("PUT")
	router.HandleFunc("/{id}/delete", articlecontroller.Delete).Methods("DELETE")

}
