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
	router.HandlerFunc(http.MethodGet, "/user/signup", app.signup)
	router.HandlerFunc(http.MethodPost, "/user/signup", app.signupPost)
	router.HandlerFunc(http.MethodGet, "/user/login", app.login)
	router.HandlerFunc(http.MethodPost, "/user/login", app.loginPost)

	router.HandlerFunc(http.MethodGet, "/user/view/:id", app.userView)
	router.HandlerFunc(http.MethodGet, "/admin/view", app.adminView)

	router.HandlerFunc(http.MethodGet, "/user/updateUser/:id", app.userUpdate)
	router.HandlerFunc(http.MethodPost, "/user/updateUser/:id", app.userUpdatePost)

	router.HandlerFunc(http.MethodGet, "/admin/updateUser/:id", app.adminUpdate)
	router.HandlerFunc(http.MethodPost, "/admin/updateUser/:id", app.adminUpdatePost)
	router.HandlerFunc(http.MethodGet, "/admin/deleteUser/:id", app.adminDelete)

	return router

}
