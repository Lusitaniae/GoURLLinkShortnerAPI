package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewLinkShortenerRouter(routes Routes) *mux.Router {
	// When StrictSlash is set to true, if the route path is "/path/", accessing "/path" will redirect
	// to the former and vice versa.
	router := mux.NewRouter().StrictSlash(true)

	// Files pointing to /static/ will be served by current dir
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("."))))

	//Feed the router the necessary information for the web service to function properly
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}
