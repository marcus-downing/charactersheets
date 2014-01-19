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
	// "net/url"
	"math"
	"strings"
	"strconv"
)

type TemplateData struct {
	BodyClass       string
	CurrentUser     *model.User
	IsAdmin         bool
	CurrentLanguage string

	Page               *Pagination
	Languages          []string
	LanguageNames      map[string]string
	LanguageCompletion map[string][4]int
	Users              []*model.User
	Sources            []*model.Source
	Entries            []*model.StackedEntry
	Translations       []*model.Translation
}

type Pagination struct {
	Page      int
	Size      int

	Offset    int
	Slice     int

	PrevPage  int
	NextPage  int
	LastPage  int

	Url       string
}

func Paginate(r *http.Request, size, datasize int) *Pagination {
	page, err := strconv.Atoi(r.FormValue("page"))
	if err !=  nil {
		page = 1
	}
	fmt.Println("Paginating: page", page, "size =", size, "data size =", datasize)
	if page < 1 {
		page = 1
	}
	lastPage := int(math.Floor(float64(datasize) / float64(size)) + 1)
	if page > lastPage {
		page = lastPage
	}

	offset := (page - 1) * size
	slice := offset + size
	if slice >= datasize {
		slice = datasize
	}

	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}
	nextPage := page + 1
	if nextPage > lastPage {
		nextPage = lastPage
	}

	baseUrl := r.URL
	query := baseUrl.Query()
	query.Del("page")
	baseUrl.RawQuery = query.Encode()

	fmt.Println("Pagination: page", page, "of", lastPage, "; offset =", offset, "slice =", slice)
	return &Pagination{
		Page: page,
		Size: size,

		Offset: offset,
		Slice: slice,

		PrevPage: prevPage,
		NextPage: nextPage,
		LastPage: lastPage,

		Url:   baseUrl.String(),
	}
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
		BodyClass:       bodyClass,
		CurrentUser:     currentUser,
		IsAdmin:         currentUser.IsAdmin,
		CurrentLanguage: currentUser.Language,
		Languages:       model.Languages,
		LanguageNames:   model.LanguageNames,
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

func translateEntry(original, partOf, language string) string {
	translations := model.GetPartTranslations(original, partOf, language)
	// fmt.Println("Found", len(translations), "translations of:", original)
	if len(translations) > 0 {
		return translations[len(translations)-1].Translation
	}
	return ""
}

func paginateTemplate(page *Pagination) template.HTML {
	url := page.Url
	if strings.Index(url, "?") != -1 {
		url = url+"&"
	} else {
		url = url+"?"
	}

	format := "<a href='%spage=%d' class='btn btn-default'>%s</a>"

	first := ""
	back := ""
	if page.Page > 1 {
		first = fmt.Sprintf(format, url, 1, "<span class='glyphicon glyphicon-arrow-left'></span> First")
		back = fmt.Sprintf(format, url, page.PrevPage, "<span class='glyphicon glyphicon-arrow-left'></span> Back")
	}

	next := ""
	last := ""
	if page.Page < page.LastPage {
		next = fmt.Sprintf(format, url, page.NextPage, "Next <span class='glyphicon glyphicon-arrow-right'></span>")
		last = fmt.Sprintf(format, url, page.LastPage, "Last <span class='glyphicon glyphicon-arrow-right'></span>")
	}

	return template.HTML("<span class='pagination'>"+first+back+next+last+"</span>")
}

var templateFuncs = template.FuncMap{
	"percentColour": percentColour,
	"md5":           md5sum,
	"translate": translateEntry,
	"pagination": paginateTemplate,
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
