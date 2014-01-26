package control

import (
	"../model"
	"fmt"
	// "code.google.com/p/go.crypto/bcrypt"
	// "crypto/md5"
	// "encoding/hex"
	// "html/template"
	// "math/rand"
	"net/http"
)

func APIMasterHandler(w http.ResponseWriter, r *http.Request) {

}

func APITranslateHandler(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r)
	if user == nil {
		fmt.Println("Unknown user")
		return
	}

	entry := model.Entry{
		Original: r.FormValue("original"),
		PartOf:   r.FormValue("partOf"),
	}
	if entry.Original == "" {
		fmt.Println("Unknown string")
		return
	}
	language := r.FormValue("language")
	translation := r.FormValue("translation")

	if language == "" {
		fmt.Println("Unknown language:", language)
		return
	}
	if translation == "" {
		fmt.Println("Blank translation:", entry.Original)
		return
	}
	fmt.Println("Adding", language, "translation for:", entry.Original)

	t := &model.Translation{ entry, language, translation, user.Email }
	t.Save()

	fmt.Fprint(w, "OK")
}

func APISetLeadHandler(w http.ResponseWriter, r *http.Request) {
	me := GetCurrentUser(r)
	if !me.IsAdmin {
		fmt.Println("Hah!")
		return
	}

	email := r.FormValue("user")
	user := model.GetUserByEmail(email)

	if user == nil {
		fmt.Print(w, "Unknown user")
		return
	}
	fmt.Println("Setting", user.Name, "as language lead for", model.LanguageNames[user.Language])
	user.SetLanguageLead()
}

func APIClearLeadHandler(w http.ResponseWriter, r *http.Request) {
	me := GetCurrentUser(r)
	if !me.IsAdmin {
		fmt.Println("Hah!")
		return
	}

	email := r.FormValue("user")
	user := model.GetUserByEmail(email)

	if user == nil {
		fmt.Print(w, "Unknown user")
		return
	}
	fmt.Println("Removing", user.Name, "as language lead for", model.LanguageNames[user.Language])
	user.ClearLanguageLead()
}