package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"teeth_datastructures/internal/model"
	"text/template"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
	users         *model.Users
	appointments  *model.Appointments
	sessions      *model.Sessions
}

func main() {

	addr := flag.String("addr", ":8080", "HTTP network address")
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Print(err)
		return
	}

	users := &model.Users{
		ErrorLog: errorLog,
	}

	patients := []*model.User{}
	users.Patients = patients
	users.LoadUsers()

	// Create the admin
	if created := users.FindAdmin(); !created {
		_, err = users.CreateAdmin()
		if err != nil {
			errorLog.Print("admin creation failed")
			return
		}
		infoLog.Print("ADMIN CREATED")
	}

	appointments := &model.Appointments{}
	bookings := []*model.Appointment{}
	appointments.Bookings = bookings
	appointments.LoadAppointments()

	sessions := &model.Sessions{}
	stores := []*model.Session{}
	sessions.Sessions = stores

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		templateCache: templateCache,
		users:         users,
		appointments:  appointments,
		sessions:      sessions,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	srv.ListenAndServe()

}
