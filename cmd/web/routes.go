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

	router.HandlerFunc(http.MethodGet, "/user/createAppointment/:id", app.createAppointmentUser)
	router.HandlerFunc(http.MethodPost, "/user/createAppointment/:id", app.createAppointmentUserPost)
	router.HandlerFunc(http.MethodGet, "/user/updateAppointment/:id/:appointmentID", app.userUpdateAppointment)
	router.HandlerFunc(http.MethodPost, "/user/updateAppointment/:id/:appointmentID", app.userUpdateAppointmentPost)
	router.HandlerFunc(http.MethodGet, "/user/deleteAppointment/:id/:appointmentID", app.userDeleteAppointment)

	router.HandlerFunc(http.MethodGet, "/admin/updateAppointment/:appointmentID", app.adminUpdateAppointment)
	router.HandlerFunc(http.MethodPost, "/admin/updateAppointment/:appointmentID", app.adminUpdateAppointmentPost)
	router.HandlerFunc(http.MethodGet, "/admin/deleteAppointment/:appointmentID", app.AdminDeleteAppointment)

	router.HandlerFunc(http.MethodGet, "/admin/sessions", app.adminSessionView)
	router.HandlerFunc(http.MethodGet, "/admin/sessions/delete/:sessionID", app.adminSessionDelete)
	return router

}
