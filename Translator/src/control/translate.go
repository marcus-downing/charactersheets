package control

import (
	"../model"
	// "code.google.com/p/go.crypto/bcrypt"
	// "crypto/md5"
	// "encoding/hex"
	// "html/template"
	// "math/rand"
	"encoding/csv"
	// "io"
	"fmt"
	"net/http"
	"path"
	"strings"
	"strconv"
)

const(
	PageSize = 10
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
		data.Page = Paginate(r, PageSize, len(data.Entries))
		data.Entries = data.Entries[data.Page.Offset:data.Page.Slice]
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
		data.Page = Paginate(r, PageSize, len(data.Entries))
		data.Entries = data.Entries[data.Page.Offset:data.Page.Slice]

		data.Translations = model.GetTranslations()
		return data
	})
}

func importMasterData(data []map[string]string, clean bool) {
	for _, record := range data {
		entry := &model.Entry{
			Original: record["Original"],
			PartOf:   record["Part of"],
		}
		entry.Save()

		filepath := record["File"]
		filename := path.Base(filepath)
		ext := path.Ext(filepath)
		name := strings.TrimSuffix(filename, ext)
		level, _ := strconv.Atoi(record["Level"])
		source := &model.Source{
			Filepath: filepath,
			Page: name,
			Volume: record["Volume"],
			Level: level,
			Game: record["Game"],
		}
		source.Save()

		count, _ := strconv.Atoi(record["Count"])
		entrySource := &model.EntrySource{
			Entry: *entry,
			Source: *source,
			Count: count,
		}
		entrySource.Save()
	}
}

func importTranslationData(data []map[string]string, clean bool) {

}

func ImportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("POST import")
		clean := r.FormValue("clean-import") == "on"
		importType := r.FormValue("type")
		if importType != "master" && importType != "translations" {
			fmt.Println("Missing type")
			http.Redirect(w, r, "/import", 303)
			return
		}
		file, _, err := r.FormFile("import-file")
		if err != nil {
			fmt.Println("Error reading file:", err)
			http.Redirect(w, r, "/import", 303)
			return
		}
		if file == nil {
			fmt.Println("Missing file")
			http.Redirect(w, r, "/import", 303)
			return
		}

		lines, err := csv.NewReader(file).ReadAll()
		if err != nil {
			fmt.Println("Error reading CSV:", err)
			http.Redirect(w, r, "/import", 303)
			return
		}
		data := associateData(lines)
		fmt.Println("Found", len(data), "lines")

		if importType == "master" {
			go importMasterData(data, clean)
		} else {
			go importTranslationData(data, clean)
		}

		http.Redirect(w, r, "/import/done", 303)
	} else {
		renderTemplate("import", w, r, func(data TemplateData) TemplateData {
			return data
		})
	}
}

func ImportDoneHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("import_done", w, r, func(data TemplateData) TemplateData {
		return data
	})
}

func ExportHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("export", w, r, func(data TemplateData) TemplateData {
		return data
	})
}

func associateData(in [][]string) []map[string]string {
	out := make([]map[string]string, 0, len(in)-1)
	fields := in[0]
	linelen := len(fields)
	for i, line := range in {
		if i == 0 {
			continue
		}
		linedata := make(map[string]string, linelen)
		for j, value := range line {
			if value != "" {
				linedata[fields[j]] = value
			}
		}
		out = append(out, linedata)
	}
	return out
}
