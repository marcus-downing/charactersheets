package control

import (
	"../model"
	// "code.google.com/p/go.crypto/bcrypt"
	// "crypto/md5"
	// "encoding/hex"
	// "html/template"
	// "math/rand"
	"net/http"
)

func SourcesHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("sources", w, r, func(data TemplateData) TemplateData {
		data.Sources = nil
		return data
	})
}

func MasterHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("master", w, r, func(data TemplateData) TemplateData {
		data.CurrentLanguage = r.FormValue("language")
		if data.CurrentLanguage == "" {
			data.CurrentLanguage = "gb"
		}
		data.Entries = model.GetStackedEntries()
		return data
	})
}

func TranslationHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("translate", w, r, func(data TemplateData) TemplateData {
		rlang := r.FormValue("language")
		if rlang != "" {
			data.CurrentLanguage = rlang
		}
		data.Entries = model.GetStackedEntries()
		data.Translations = model.GetTranslations()
		return data
	})
}

func ImportHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("import", w, r, func(data TemplateData) TemplateData {
		return data
	})
}

func ExportHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("export", w, r, func(data TemplateData) TemplateData {
		return data
	})
}
