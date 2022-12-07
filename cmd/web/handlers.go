package main

import (
	"fmt"
	"net/http"
	"teeth_datastructures/internal/validator"

	"github.com/julienschmidt/httprouter"
)

type userForm struct {
	Name     string
	Username string
	Password string
	Email    string
	ID       string
	validator.Validator
}

type loginForm struct {
	Username string
	Password string
	validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	form := &userForm{}
	data.Form = form
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.infoLog.Print("Error parsing user form")
		app.serverError(w, http.StatusInternalServerError)
		return
	}

	form := &userForm{
		Name:     r.PostForm.Get("name"),
		Username: r.PostForm.Get("username"),
		Password: r.PostForm.Get("password"),
		Email:    r.PostForm.Get("email"),
	}
	form.IsBlank(form.Name, "Name")
	form.IsBlank(form.Password, "Password")
	form.IsBlank(form.Username, "username")
	form.IsBlank(form.Email, "email")
	form.ValidEmail(form.Email)
	if valid := form.Valid(); !valid {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	userID, err := app.users.CreateUser(form.Name, form.Username, form.Password, form.Email)
	if err != nil {
		app.errorLog.Print("Error creating user.")
		return
	}

	isAdmin, err := app.users.GetUser(userID)
	if err != nil {
		app.serverError(w, http.StatusNotFound)
		return
	}
	if isAdmin.Admin {
		http.Redirect(w, r, "/admin/patients", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/user/view/%v", userID), http.StatusSeeOther)
	}

}

// LOGIN
func (app *application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	form := &loginForm{}
	data.Form = form
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Print("Failed to parse form")
		app.serverError(w, http.StatusUnprocessableEntity)
		return
	}

	form := &loginForm{
		Username: r.PostForm.Get("username"),
		Password: r.PostForm.Get("password"),
	}

	form.IsBlank(form.Username, "username")
	form.IsBlank(form.Password, "password")
	if valid := form.Valid(); !valid {
		app.infoLog.Print("HERE")
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	// Authenticate user
	valid, authenticatedUser := app.users.Authenticate(form.Username, form.Password)
	if !valid {
		data := app.newTemplateData(r)
		form.AddNonFieldErrors("Invalid password")
		data.Form = form
		app.render(w, http.StatusUnauthorized, "login.html", data)
		return
	}
	if authenticatedUser.Admin {
		http.Redirect(w, r, "/admin/view", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/user/view/%v", authenticatedUser.ID), http.StatusSeeOther)

	}

}

// All the views
func (app *application) userView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	data := app.newTemplateData(r)
	data.ID = id
	data.Admin = false
	app.render(w, http.StatusOK, "view.html", data)
}

func (app *application) adminView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Users = app.users
	data.Admin = true
	app.render(w, http.StatusOK, "patients.html", data)
}

func (app *application) userUpdate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = &loginForm{}
	data.Admin = false
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")
	data.ID = id
	app.render(w, http.StatusOK, "userUpdate.html", data)
}

func (app *application) userUpdatePost(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userID := params.ByName("id")

	app.infoLog.Print("HITS")

	err := r.ParseForm()
	if err != nil {
		app.errorLog.Print("error parsing update user form")
		data := app.newTemplateData(r)
		app.render(w, http.StatusUnprocessableEntity, fmt.Sprintf("/user/view/%v", userID), data)
		return
	}

	form := &userForm{
		Name:     r.PostForm.Get("name"),
		Username: r.PostForm.Get("username"),
		Password: r.PostForm.Get("password"),
		Email:    r.PostForm.Get("email"),
	}
	form.IsBlank(form.Name, "Name")
	form.IsBlank(form.Password, "Password")
	form.IsBlank(form.Username, "username")
	form.IsBlank(form.Email, "email")
	form.ValidEmail(form.Email)
	if valid := form.Valid(); !valid {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "userUpdate.html", data)
		return
	}

	err = app.users.UpdateUser(userID, form.Name, form.Username, form.Email)
	if err != nil {
		app.errorLog.Print("Failed to update")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/user/view/%v", userID), http.StatusSeeOther)

}

// CRUD admin
func (app *application) adminDelete(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userID := params.ByName("id")

	app.users.DeleteUser(userID)
	http.Redirect(w, r, "/admin/view", http.StatusSeeOther)
}

func (app *application) adminUpdate(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")

	data := app.newTemplateData(r)
	data.Form = &loginForm{}
	data.ID = id
	app.render(w, http.StatusOK, "adminUpdateUser.html", data)
}

func (app *application) adminUpdatePost(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	userID := params.ByName("id")

	err := r.ParseForm()
	if err != nil {
		app.errorLog.Print("error parsing update user form")
		data := app.newTemplateData(r)
		app.render(w, http.StatusUnprocessableEntity, "/admin/view", data)
		return
	}

	form := &userForm{
		Name:     r.PostForm.Get("name"),
		Username: r.PostForm.Get("username"),
		Password: r.PostForm.Get("password"),
		Email:    r.PostForm.Get("email"),
		ID:       userID,
	}
	form.IsBlank(form.Name, "Name")
	form.IsBlank(form.Password, "Password")
	form.IsBlank(form.Username, "username")
	form.IsBlank(form.Email, "email")
	form.ValidEmail(form.Email)
	if valid := form.Valid(); !valid {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "adminUpdateUser.html", data)
		return
	}

	err = app.users.UpdateUser(userID, form.Name, form.Username, form.Email)
	if err != nil {
		app.errorLog.Print("Failed to update")
		return
	}

	http.Redirect(w, r, "/admin/view", http.StatusSeeOther)

}
