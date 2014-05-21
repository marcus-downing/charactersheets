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
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type TemplateData struct {
	BodyClass       string
	CurrentUser     *model.User
	IsAdmin         bool
	CurrentLanguage string
	RecentUsers     []RecentUser

	User               *model.User
	Page               *Pagination
	Languages          []string
	DisplayLanguages   []string
	LanguageNames      map[string]string
	LanguageCompletion map[string][4]int
	Users              []*model.User
	Sources            []*model.Source
	Entries            []*model.StackedEntry
	Translations       []*model.Translation
	CurrentGame        string
	CurrentLevel       string
	CurrentShow        string
}

type Pagination struct {
	Page int
	Size int

	Offset int
	Slice  int

	PrevPage int
	NextPage int
	LastPage int

	Url string
}

func Paginate(r *http.Request, size, datasize int) *Pagination {
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}
	fmt.Println("Paginating: page", page, "size =", size, "data size =", datasize)
	if page < 1 {
		page = 1
	}
	lastPage := int(math.Floor(float64(datasize)/float64(size)) + 1)
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
		Slice:  slice,

		PrevPage: prevPage,
		NextPage: nextPage,
		LastPage: lastPage,

		Url: baseUrl.String(),
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
	recentUsers := GetRecentUsers()

	if currentUser == nil {
		return TemplateData{
			BodyClass:        bodyClass,
			IsAdmin:          false,
			CurrentLanguage:  "gb",
			Languages:        model.Languages,
			DisplayLanguages: model.DisplayLanguages,
			LanguageNames:    model.LanguageNames,
			RecentUsers:      recentUsers,
		}
	}
	return TemplateData{
		BodyClass:        bodyClass,
		CurrentUser:      currentUser,
		IsAdmin:          currentUser.IsAdmin,
		CurrentLanguage:  currentUser.Language,
		Languages:        model.Languages,
		DisplayLanguages: model.DisplayLanguages,
		LanguageNames:    model.LanguageNames,
		RecentUsers:      recentUsers,
	}
}

func DurString(dur time.Duration) string {
	minutes := int(dur.Minutes())
	hours := int(dur.Hours())
	days := int(hours / 24)

	if days > 0 {
		return fmt.Sprintf("%d days ago", days)
	}
	if hours > 0 {
		return fmt.Sprintf("%d hours ago", hours)
	}
	if minutes > 0 {
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	return "just now"
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

type TranslationSet struct {
	Entry        *model.StackedEntry
	Others       []*model.StackedTranslation
	Mine         *model.StackedTranslation
	Untranslated bool
}

func getTranslations(entry *model.StackedEntry, language string, me *model.User) *TranslationSet {
	translations := entry.GetTranslations(language)

	others := make([]*model.StackedTranslation, 0, len(translations))
	var mine *model.StackedTranslation = nil

	for _, translation := range translations {
		if translation != nil {
			if translation.Translator == me.Email {
				mine = translation
			} else {
				others = append(others, translation)
			}
		}
	}

	return &TranslationSet{
		Entry:        entry,
		Others:       others,
		Mine:         mine,
		Untranslated: mine == nil,
	}
}

func myTranslation(set *TranslationSet) *model.StackedTranslation {
	if set == nil || set.Mine == nil {
		return &model.StackedTranslation{}
	}
	return set.Mine
}

func otherTranslations(set *TranslationSet) []*model.StackedTranslation {
	return set.Others
}

func countUserTranslations(user *model.User) map[string]int {
	return user.CountTranslations()
}

func countEntryTranslations(entry *model.StackedEntry) map[string]int {
	return entry.CountTranslations()
}

func profileTranslations(user *model.User) map[string]*model.TranslationProfile {
	return model.ProfileTranslations(user)
}

func entryClass(entry *model.StackedEntry, language string, me *model.User) string {
	classes := make([]string, 0, 20)

	translations := entry.GetTranslations(language)
	if len(translations) == 0 {
		classes = append(classes, "untranslated")
	}
	for _, translation := range translations {
		if translation.Translator == me.Email {
			classes = append(classes, "with-translation")
		}
	}

	//  TODO more classes
	return strings.Join(classes, " ")
}

func paginateTemplate(page *Pagination) template.HTML {
	url := page.Url
	if strings.Index(url, "?") != -1 {
		url = url + "&"
	} else {
		url = url + "?"
	}

	format := "<a href='%spage=%d' class='btn btn-default'>%s</a>"
	disabled := "<span class='btn btn-default' disabled='disabled'>%s</span>"

	first := "<span class='glyphicon glyphicon-chevron-left'></span> First"
	back := "<span class='glyphicon glyphicon-chevron-left'></span> Back"
	if page.Page > 1 {
		first = fmt.Sprintf(format, url, 1, first)
		back = fmt.Sprintf(format, url, page.PrevPage, back)
	} else {
		first = fmt.Sprintf(disabled, first)
		back = fmt.Sprintf(disabled, back)
	}

	next := "Next <span class='glyphicon glyphicon-chevron-right'></span>"
	last := "Last <span class='glyphicon glyphicon-chevron-right'></span>"
	if page.Page < page.LastPage {
		next = fmt.Sprintf(format, url, page.NextPage, next)
		last = fmt.Sprintf(format, url, page.LastPage, last)
	} else {
		next = fmt.Sprintf(disabled, next)
		last = fmt.Sprintf(disabled, last)
	}

	return template.HTML("<span class='pagination btn-group'>" + first + back + next + last + "</span>")
}

func sourcePath(source *model.Source) template.HTML {
	ext := path.Ext(source.Filepath)
	path := strings.TrimSuffix(source.Filepath, ext)
	parts := strings.Split(path, "/")
	lis := strings.Join(parts, "</li><li>")
	return template.HTML("<ol class='breadcrumb'><li>" + lis + "</li></ol>")
}

func sourceURL(source *model.Source) template.HTML {
	path := source.Filepath
	path = strings.Replace(path, "3.5", "dnd35", 1)
	path = strings.Replace(path, "Pathfinder", "pathfinder", 1)
	return template.HTML("/pdf/" + path + ".pdf")
}

func previewURL(language string, source *model.Source) template.HTML {
	languagePath := model.LanguagePaths[language]
	if languagePath != "" {
		languagePath = "languages/" + languagePath
	}
	path := source.Filepath
	path = strings.Replace(path, "3.5", "dnd35", 1)
	path = strings.Replace(path, "Pathfinder", "pathfinder", 1)
	return template.HTML("/pdf/" + languagePath + "/" + path + ".pdf")
}

func previewExists(language string, source *model.Source) bool {
	languagePath := model.LanguagePaths[language]
	if languagePath != "" {
		languagePath = "languages/" + languagePath
	}
	path := source.Filepath
	path = strings.Replace(path, "3.5", "dnd35", 1)
	path = strings.Replace(path, "Pathfinder", "pathfinder", 1)
	fullPath := "../../Composer 2.1.3/public/pdf/" + languagePath + "/" + path + ".pdf"

	fi, err := os.Stat(fullPath)
	return err == nil && !fi.IsDir()
}

func sourceCompletion(source *model.Source) map[string]int {
	return source.GetLanguageCompletion()
}

func isVotedUp(translation *model.StackedTranslation, voter *model.User) bool {
	votes := translation.GetVotes()
	for _, vote := range votes {
		fmt.Println("Vote by", vote.Voter.Email, "=", vote.Vote)
		if vote.Voter.Email == voter.Email {
			return vote.Vote
		}
	}
	return false
}

func isVotedDown(translation *model.StackedTranslation, voter *model.User) bool {
	votes := translation.GetVotes()
	for _, vote := range votes {
		fmt.Println("Vote by", vote.Voter.Email, "=", vote.Vote)
		if vote.Voter.Email == voter.Email {
			return !vote.Vote
		}
	}
	return false
}

var templateFuncs = template.FuncMap{
	"percentColour":          percentColour,
	"md5":                    md5sum,
	"getTranslations":        getTranslations,
	"otherTranslations":      otherTranslations,
	"myTranslation":          myTranslation,
	"countUserTranslations":  countUserTranslations,
	"countEntryTranslations": countEntryTranslations,
	"profileTranslations":    profileTranslations,
	"entryClass":             entryClass,
	"pagination":             paginateTemplate,
	"sourcePath":             sourcePath,
	"sourceURL":              sourceURL,
	"sourceCompletion":       sourceCompletion,
	"previewURL":             previewURL,
	"previewExists":          previewExists,
	"isVotedUp":              isVotedUp,
	"isVotedDown":            isVotedDown,
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

//  Recent users

type RecentUser struct {
	User        *model.User
	LoggedInFor string
}

var recentUsers map[string]time.Time = make(map[string]time.Time, 100)

func PingUser(email string) {
	recentUsers[email] = time.Now()
}

func PingCurrentUser(r *http.Request) {
	session := seshcookie.Session.Get(r)
	if session == nil {
		return
	}
	if id, ok := session["user"].(string); ok && id != "" {
		recentUsers[id] = time.Now()
	}
}

func GetRecentUsers() []RecentUser {
	threshold, _ := time.ParseDuration("168h") // 7 days
	recent := make([]RecentUser, 0, len(recentUsers))
	for email, t := range recentUsers {
		user := model.GetUserByEmail(email)
		if user == nil {
			continue
		}
		dur := time.Since(t)
		if dur.Hours() > threshold.Hours() {
			continue
		}

		durstr := DurString(dur)
		recent = append(recent, RecentUser{user, durstr})
	}
	return recent
}
