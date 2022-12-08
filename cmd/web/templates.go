package main

import (
	"net/http"
	"path/filepath"
	"teeth_datastructures/internal/model"
	"text/template"
)

// There might be situations where we want to output multiple pieces of data.
type templateData struct {
	Users            *model.Users
	Form             any
	Session          string
	ID               string
	AppointmentID    string
	Admin            bool
	Appointments     []*model.Appointment
	AppointmentModel *model.Appointments
	Sessions         *model.Sessions
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.ParseFiles("ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts

	}
	return cache, nil

}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{}
}
