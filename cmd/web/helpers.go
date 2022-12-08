package main

import (
	"net/http"
	"teeth_datastructures/internal/model"
)

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		app.errorLog.Print("Template not found")
		app.notFound(w)
		return
	}

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

func (app *application) sessionManagement(w http.ResponseWriter, r *http.Request) (*model.Session, bool) {
	// Authenticate through sessions
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			app.serverError(w, http.StatusUnauthorized)
			return nil, false
		}
		// For any other type of error, return a bad request status
		app.serverError(w, http.StatusBadRequest)
		return nil, false
	}
	sessionToken := c.Value
	valid, validSession := app.sessions.Validate(sessionToken)
	if !valid {
		app.serverError(w, http.StatusUnauthorized)
		return nil, false
	}

	// Check expiry
	if expired := validSession.Expired(); expired {
		app.sessions.RemoveSession(validSession.Username)
		app.serverError(w, http.StatusUnauthorized)
		return nil, false
	}

	return validSession, true
}
