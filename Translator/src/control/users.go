package control

import (
	"../model"
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"net/http"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Dashboard")

	renderTemplate("home", w, r, func(data TemplateData) TemplateData {
		data.LanguageCompletion = model.GetLanguageCompletion()
		return data
	})
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		name := r.FormValue("name")
		language := r.FormValue("language")
		user := &model.User{
			Email:    email,
			Name:     name,
			Language: language,
			Password: "",
			Secret:   "",
		}
		user.Save()

		http.Redirect(w, r, "/users", 303)
	} else {
		renderTemplate("users", w, r, func(data TemplateData) TemplateData {
			data.Users = model.GetUsers()
			return data
		})
	}
}

func UsersAddHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("users_add", w, r, nil)
}

func UsersDelHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := GetCurrentUser(r)
	if !currentUser.IsAdmin {
		http.Redirect(w, r, "/users", 303)
		return
	}

	email := r.FormValue("user")
	user := model.GetUserByEmail(email)
	if user == nil {
		http.Redirect(w, r, "/users", 303)
		return
	}

	gonow := r.FormValue("go")
	if r.Method == "POST" && gonow == "yes" {
		user.Delete()
		http.Redirect(w, r, "/users", 303)
		return
	} else {
		renderTemplate("users_del", w, r, func(data TemplateData) TemplateData {
			data.User = user
			return data
		})
	}
}

func AccountHandler(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r)

	if r.Method == "POST" {
		user.Name = r.FormValue("name")
		language := r.FormValue("language")
		if language != "" {
			user.Language = language
		}
		user.Save()

		http.Redirect(w, r, "/home", 303)
	} else {
		renderTemplate("account", w, r, func(data TemplateData) TemplateData {
			return data
		})
	}
}

func AccountReclaimHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("account_reclaim", w, r, nil)
}

func SetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r)

	if r.Method == "POST" {
		password := r.FormValue("password")
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err == nil {
			user.Password = string(hash)
			user.Save()
		}
		http.Redirect(w, r, "/account", 303)
	} else {
		renderTemplate("account_set_password", w, r, func(data TemplateData) TemplateData {
			return data
		})
	}
}
