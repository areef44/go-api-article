package routes

import (
	"go-api-article/controllers/categorycontroller"

	"github.com/gorilla/mux"
)

func CategoryRoutes(r *mux.Router) {
	router := r.PathPrefix("/categories").Subrouter()

	router.HandleFunc("", categorycontroller.Index).Methods("GET")
	router.HandleFunc("", categorycontroller.Create).Methods("POST")
	router.HandleFunc("/{id}/detail", categorycontroller.Detail).Methods("GET")
	router.HandleFunc("/{id}/update", categorycontroller.Update).Methods("PUT")
	router.HandleFunc("/{id}/delete", categorycontroller.Delete).Methods("DELETE")

}
