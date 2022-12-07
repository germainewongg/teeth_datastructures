package main

import (
	"net/http"
)

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		app.errorLog.Print("Template not found")
		app.notFound(w)
		return
	}

	w.WriteHeader(status)

	// Execute template and write to body
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.clientError(w, err)
	}
}

func (app *application) clientError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (app *application) serverError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.serverError(w, http.StatusNotFound)
}
