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
	"strconv"
	"strings"
	"time"
	"bufio"
	"mime/multipart"
)

const (
	PageSize = 20
)

func SourcesHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("sources", w, r, func(data TemplateData) TemplateData {
		data.Sources = nil
		return data
	})
}

func EntriesHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("entries", w, r, func(data TemplateData) TemplateData {
		data.CurrentGame = r.FormValue("game")
		data.CurrentLevel = r.FormValue("level")
		data.CurrentShow = r.FormValue("show")

		data.Entries = model.GetStackedEntries(data.CurrentGame, data.CurrentLevel, data.CurrentShow, "gb")
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

		data.CurrentGame = r.FormValue("game")
		data.CurrentLevel = r.FormValue("level")
		data.CurrentShow = r.FormValue("show")

		data.Entries = model.GetStackedEntries(data.CurrentGame, data.CurrentLevel, data.CurrentShow, data.CurrentLanguage)
		data.Page = Paginate(r, PageSize, len(data.Entries))
		data.Entries = data.Entries[data.Page.Offset:data.Page.Slice]
		return data
	})
}

func importMasterData(data []map[string]string) {
	sleepTime, _ := time.ParseDuration("5ms")
	fmt.Println("Importing", len(data), "master records")
	for _, record := range data {
		// fmt.Println("Inserting translation:", record["Original"], ";", record["Part of"])
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
			Page:     name,
			Volume:   record["Volume"],
			Level:    level,
			Game:     record["Game"],
		}
		source.Save()

		count, _ := strconv.Atoi(record["Count"])
		entrySource := &model.EntrySource{
			Entry:  *entry,
			Source: *source,
			Count:  count,
		}
		entrySource.Save()
		time.Sleep(sleepTime)
	}
	fmt.Println("Import complete")
}

func importTranslationData(data []map[string]string, language string, translator *model.User) {
	sleepTime, _ := time.ParseDuration("5ms")
	fmt.Println("Importing", len(data), "translation records as", translator.Name)
	num := 0
	for _, record := range data {
		t := record["Translation"]
		if t == "" {
			continue
		}
		translation := &model.Translation{
			Entry: model.Entry{
				Original: record["Original"],
				PartOf:   record["Part of"],
			},
			Language:    language,
			Translation: t,
			Translator:  translator.Email,
		}
		translation.Save()
		time.Sleep(sleepTime)
		num++
	}
	fmt.Println("Import complete:", num, "of", len(data))
}

func ImportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("POST import")
		// clean := r.FormValue("clean-import") == "on"
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

		file2 := stripBOM(file)
		lines, err := csv.NewReader(file2).ReadAll()
		if err != nil {
			fmt.Println("Error reading CSV:", err)
			http.Redirect(w, r, "/import", 303)
			return
		}
		file.Close()
		data := associateData(lines)
		fmt.Println("Found", len(data), "lines")

		if importType == "master" {
			go importMasterData(data)
		} else {
			language := r.FormValue("language")
			translator := model.GetUserByEmail(r.FormValue("translator"))
			go importTranslationData(data, language, translator)
		}

		http.Redirect(w, r, "/import/done", 303)
	} else {
		renderTemplate("import", w, r, func(data TemplateData) TemplateData {
			data.Users = model.GetUsers()
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
	if r.Method == "POST" {
		language := r.FormValue("language")
		fmt.Println("Exporting in", language)
		translations := model.GetPreferredTranslations(language)

		w.Header().Set("Content-Encoding", "UTF-8")
		w.Header().Set("Content-Type", "application/csv; charset=UTF-8")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+model.LanguageNames[language]+".csv\"")

		out := csv.NewWriter(w)
		out.Write([]string{
			"Original",
			"Part of",
			"Translation",
		})
		for _, translation := range translations {
			out.Write([]string{
				translation.Entry.Original,
				translation.Entry.PartOf,
				translation.Translation,
			})
		}
		out.Flush()
		return
	} else {
		renderTemplate("export", w, r, func(data TemplateData) TemplateData {
			return data
		})
	}
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

func stripBOM(file multipart.File) *bufio.Reader {
	br := bufio.NewReader(file)
	rune, _, _ := br.ReadRune()
	if rune != '\uFEFF' {
        br.UnreadRune() // Not a BOM -- put the rune back
    }
    return br
}