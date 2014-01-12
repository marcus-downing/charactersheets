package control

import (
	"../model"
	// "code.google.com/p/go.crypto/bcrypt"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/bpowers/seshcookie"
	"html/template"
	// "math/rand"
	"net/http"
)

type TemplateData struct {
	BodyClass       string
	CurrentUser     *model.User
	IsAdmin         bool
	CurrentLanguage string

	Languages          []string
	LanguageNames      map[string]string
	LanguageCompletion map[string][4]int
	Users              []*model.User
	Sources            []*model.Source
	Entries            []*model.StackedEntry
	Translations       []*model.Translation
}

func SetCurrentUser(user *model.User, r *http.Request) {
	session := seshcookie.Session.Get(r)
	if session == nil {
		return
	}
	if user == nil {
		session["id"] = nil
	} else {
		session["id"] = user.Email
	}
}

func GetCurrentUser(r *http.Request) *model.User {
	session := seshcookie.Session.Get(r)
	if session == nil {
		return nil
	}
	if id, ok := session["user"].(string); ok {
		if id == "" {
			return nil
		}
		return model.GetUserByEmail(id)
	}
	return nil
}

func GetTemplateData(r *http.Request, bodyClass string) TemplateData {
	currentUser := GetCurrentUser(r)
	return TemplateData{
		BodyClass:     bodyClass,
		CurrentUser:   currentUser,
		IsAdmin:       currentUser.IsAdmin,
		CurrentLanguage: currentUser.Language,
		Languages:     model.Languages,
		LanguageNames: model.LanguageNames,
	}
}

func percentColour(pc int) string {
	if pc >= 90 {
		return "success"
	} else if pc >= 60 {
		return "info"
	} else if pc >= 30 {
		return "warning"
	} else {
		return "danger"
	}
}

func md5sum(email string) string {
	hasher := md5.New()
	hasher.Write([]byte(email))
	return hex.EncodeToString(hasher.Sum(nil))
}

// func optionalTemplate(name string, data interface{}) bool {
//     if t := template.Lookup(name) {
//     	t.Execute(w, data)
//     }
// }

func translateEntry(original, partOf, language string) string {
	translations := model.GetPartTranslations(original, partOf, language)
	fmt.Println("Found", len(translations), "translations of:", original)
	if len(translations) > 0 {
		return translations[len(translations)-1].Translation
	}
	return ""
}

var templateFuncs = template.FuncMap{
	"percentColour": percentColour,
	"md5":           md5sum,
	// "optionalTemplate":       optionalTemplate,
	"translate": translateEntry,
}

func renderTemplate(name string, w http.ResponseWriter, r *http.Request, dataproc func(data TemplateData) TemplateData) {
	var data = GetTemplateData(r, name)
	if dataproc != nil {
		data = dataproc(data)
	}
	fmt.Println("Rendering page:", name)

	t, err := template.New("_base.html").Funcs(templateFuncs).ParseFiles("../view/_base.html", "../view/"+name+".html")
	if err != nil {
		fmt.Fprint(w, "Error:", err)
		fmt.Println("Error:", err)
	}
	err = t.Execute(w, data)
	if err != nil {
		fmt.Fprint(w, "Error:", err)
		fmt.Println("Error:", err)
	}
}
