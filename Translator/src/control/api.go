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
	language := r.FormValue("language")
	name := r.FormValue("name")
	translation := r.FormValue("translation")

	if language == "" {
		fmt.Println("Unknown language:", language)
		return
	}
	if translation == "" {
		fmt.Println("Blank translation:", name)
		return
	}
	fmt.Println("Adding", language, "translation for:", name)
	model.AddTranslation(name, language, translation, user);
}