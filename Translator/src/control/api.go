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

	model.AddTranslation(&entry, language, translation, user)

	fmt.Fprint(w, "OK")
}
