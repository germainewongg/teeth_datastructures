package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	fileserver := http.FileServer(http.Dir("ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static/", fileserver))

	router.HandlerFunc(http.MethodGet, "/", app.home)

	return router

}
