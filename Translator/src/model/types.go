package model

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	"strings"
)

// ** Entries

type Entry struct {
	Original string
	PartOf   string
}

func (entry *Entry) ID() string {
	if entry == nil {
		return ""
	}

	var str = entry.Original
	if entry.PartOf != "" && entry.PartOf != entry.Original {
		str = entry.Original + "  ----  " + entry.PartOf
	}
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

const entryFields = "Original, PartOf"

func parseEntry(rows *sql.Rows) (Result, error) {
	e := Entry{}
	err := rows.Scan(&e.Original, &e.PartOf)
	fmt.Println("Entry ID: " + e.ID() + " (" + string(len(e.ID())) + ")")
	return e, err
}

func makeEntries(results []Result) []*Entry {
	entries := make([]*Entry, len(results))
	for i, result := range results {
		if entry, ok := result.(Entry); ok {
			entries[i] = &entry
		}
	}
	return entries
}

func CountEntries() int {
	return query("select count(*) from Entries").count()
}

func GetEntryByID(id string) *Entry {
	result := query("select "+entryFields+" from Entries where EntryID = ?", id).row(parseEntry)
	if entry, ok := result.(Entry); ok {
		return &entry
	}
	return nil
}

func GetEntries() []*Entry {
	results := query("select " + entryFields + " from Entries").rows(parseEntry)
	return makeEntries(results)
}

func GetEntriesAt(game string, level int, show, search, language string, translator *User) []*Entry {
	if game == "" && level == 0 && show == "" && search == "" {
		return GetEntries()
	}
	args := make([]interface{}, 0, 2)
	sql := "select Original, PartOf from Entries " +
		"inner join EntrySources on Entries.EntryID = EntrySources.EntryID " +
		"inner join Sources on SourcePath = Filepath"
	if show == "conflicts" {
		sql = sql + " inner join Translations Mine on EntryID = Mine.EntryID and Mine.Language = ? and Mine.Translator = ?" +
			"inner join Translations Others on EntryID = Others.EntryID and Others.Language = ? and Others.Translator != ?"
		args = append(args, language)
		args = append(args, translator.Email)
		args = append(args, language)
		args = append(args, translator.Email)
	} else if show == "mine" {
		sql = sql + " inner join Translations Mine on EntryID = Mine.EntryID and Mine.Language = ? and Mine.Translator = ?"
		args = append(args, language)
		args = append(args, translator.Email)
	} else if show == "others" {
		sql = sql + " inner join Translations Others on EntryID = Others.EntryID and Others.Language = ? and Others.Translator = ?"
		args = append(args, language)
		args = append(args, translator.Email)
	} else if show != "" {
		sql = sql + " left join Translations on EntryID = Translations.EntryID and Translations.Language = ?"
		args = append(args, language)
	}
	sql = sql + " where 1 = 1"

	if game != "" {
		sql = sql + " and Game = ?"
		args = append(args, game)
	}
	if level != 0 {
		sql = sql + " and Level = ?"
		args = append(args, level)
	}
	if show == "conflicts" {
		sql = sql + " and Mine.Translation = Others.Translation"
	}
	// if show != "" {
	// 	sql = sql+" and Translations.Language = ?"
	// 	args = append(args, language)
	// }
	if search != "" {
		searchTerms := strings.Split(search, " ")
		fmt.Println("Searching for:", search)
		for _, term := range searchTerms {
			term = strings.ToLower(term)
			sql = sql + " and lower(Original) like ?"
			args = append(args, "%"+term+"%")
		}
	}

	sql = sql + " group by Original, PartOf"
	if show == "translated" {
		sql = sql + " having count(Translations.Translation) > 0"
	} else if show == "untranslated" {
		sql = sql + " having count(Translations.Translation) = 0"
	}
	fmt.Println("Get entries:", sql)
	results := query(sql, args...).rows(parseEntry)
	return makeEntries(results)
}

func (entry *Entry) Save() {
	keyfields := map[string]interface{}{
		"EntryID": entry.ID(),
	}
	fields := map[string]interface{}{
		"Original": entry.Original,
		"PartOf":   entry.PartOf,
	}
	saveRecord("Entries", keyfields, fields)
}

func (entry *Entry) CountTranslations() map[string]int {
	counts := make(map[string]int, len(Languages))
	query("select Language, Count(*) from Translations where EntryID = ? group by Language", entry.ID()).rows(func(rows *sql.Rows) (Result, error) {
		var language string
		var count int
		rows.Scan(&language, &count)
		counts[language] = count
		return nil, nil
	})
	return counts
}

func (entry *Entry) GetParts() []*Entry {
	results := query("select "+entryFields+" from Entries where PartOf = ?", entry.PartOf).rows(parseEntry)
	return makeEntries(results)
}

// ** Sources

type Source struct {
	Filepath string
	Page     string
	Volume   string
	Level    int
	Game     string
}

func parseSource(rows *sql.Rows) (Result, error) {
	s := Source{}
	err := rows.Scan(&s.Filepath, &s.Page, &s.Volume, &s.Level, &s.Game)
	return s, err
}

const sourceFields = "Filepath, Page, Volume, Level, Game"

func GetSources() []*Source {
	results := query("select " + sourceFields + " from Sources").rows(parseSource)

	sources := make([]*Source, len(results))
	for i, result := range results {
		if source, ok := result.(Source); ok {
			sources[i] = &source
		}
	}
	return sources
}
func GetSourcesAt(game string, level int, show string) []*Source {
	if game == "" && level == 0 && show == "" {
		return GetSources()
	}
	args := make([]interface{}, 0, 2)
	sql := "select " + sourceFields + " from Sources "

	sql = sql + " where 1 = 1"

	if game != "" {
		sql = sql + " and Game = ?"
		args = append(args, game)
	}
	if level != 0 {
		sql = sql + " and Level = ?"
		args = append(args, level)
	}

	// sql = sql+" group by Original, PartOf"
	if show == "translated" || show == "untranslated" {
		sql = sql + " and Filepath"
		if show == "untranslated" {
			sql = sql + " not"
		}
		sql = sql + " in (select SourcePath from EntrySources" +
			" inner join Translations on EntrySources.EntryID = Translations.EntryID)"
	}

	fmt.Println("Get entries:", sql)
	results := query(sql, args...).rows(parseSource)

	sources := make([]*Source, 0, len(results))
	for _, result := range results {
		if source, ok := result.(Source); ok {
			sources = append(sources, &source)
		}
	}
	return sources
}

func (source *Source) Save() {
	keyfields := map[string]interface{}{
		"Filepath": source.Filepath,
	}
	fields := map[string]interface{}{
		"Page":   source.Page,
		"Volume": source.Volume,
		"Level":  source.Level,
		"Game":   source.Game,
	}
	saveRecord("Sources", keyfields, fields)
}

func (source *Source) GetLanguageCompletion() map[string]int {
	var completion = make(map[string]int, len(Languages))

	total := query("select count(distinct EntryID) from Entries "+
		"inner join EntrySources on Entries.EntryID = EntrySources.EntryID "+
		"where SourcePath = ?", source.Filepath).count()
	for _, lang := range Languages {
		count := query("select count(distinct Translations.EntryID) from Translations "+
			"inner join EntrySources on Translations.EntryID = EntrySources.EntryID "+
			"where SourcePath = ? and Language = ?", source.Filepath, lang).count()
		completion[lang] = 100 * count / total
	}
	return completion
}

type EntrySource struct {
	Entry  Entry
	Source Source
	Count  int
}

func parseEntrySource(rows *sql.Rows) (Result, error) {
	es := EntrySource{}
	var entryID string
	err := rows.Scan(&entryID, &es.Source.Filepath, &es.Source.Page, &es.Source.Volume, &es.Source.Level, &es.Source.Game, &es.Count)
	if entry := GetEntryByID(entryID); entry == nil {
		return nil, nil
	} else {
		es.Entry = *entry
	}
	return es, err
}

const entrySourceFields = ""

func GetEntrySources() []*EntrySource {
	results := query("select EntryID, SourcePath, Sources.Page, Sources.Volume, Sources.Level, Sources.Game, Count " +
		"from EntrySources inner join Sources on EntrySources.SourcePath = Sources.Filepath").rows(parseEntrySource)

	sources := make([]*EntrySource, len(results))
	for i, result := range results {
		if source, ok := result.(EntrySource); ok {
			sources[i] = &source
		}
	}
	return sources
}

func GetSourcesForEntry(entry *Entry) []*EntrySource {
	results := query("select EntryID, SourcePath, Page, Volume, Level, Game, Count "+
		"from EntrySources inner join Sources on SourcePath = Filepath "+
		"where EntryID = ?", entry.ID()).rows(parseEntrySource)

	sources := make([]*EntrySource, len(results))
	for i, result := range results {
		if source, ok := result.(EntrySource); ok {
			sources[i] = &source
		}
	}
	return sources
}

func (es *EntrySource) Save() {
	keyfields := map[string]interface{}{
		"EntryID":    es.Entry.ID(),
		"SourcePath": es.Source.Filepath,
	}
	fields := map[string]interface{}{
		"Count": es.Count,
	}
	saveRecord("EntrySources", keyfields, fields)
}

// ** Translations

type Translation struct {
	Entry        Entry
	Language     string
	Translation  string
	Translator   string
	IsPreferred  bool
	IsConflicted bool
}

func (translation *Translation) ID() string {
	if translation == nil {
		return ""
	}

	var str = translation.Entry.ID() + "  ---  " + translation.Language + "  ---  " + translation.Translator
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func parseTranslation(rows *sql.Rows) (Result, error) {
	t := Translation{}
	var entryID string
	err := rows.Scan(&entryID, &t.Language, &t.Translation, &t.Translator, &t.IsPreferred, &t.IsConflicted)
	if entry := GetEntryByID(entryID); entry == nil {
		return nil, nil
	} else {
		t.Entry = *entry
	}
	return t, err
}

const translationFields = "EntryID, Language, Translation, Translator, IsPreferred, IsConflicted"

func GetTranslations() []*Translation {
	results := query("select " + translationFields + " from Translations").rows(parseTranslation)
	translations := make([]*Translation, len(results))
	for i, result := range results {
		if translation, ok := result.(Translation); ok {
			translations[i] = &translation
		}
	}
	return translations
}

func GetTranslationByID(id string) *Translation {
	result := query("select "+translationFields+" from Translations where TranslationID = ?", id).row(parseTranslation)
	if translation, ok := result.(Translation); ok {
		return &translation
	}
	return nil
}

func GetTranslationsForLanguage(language string) []*Translation {
	results := query("select "+translationFields+" from Translations where Language = ?", language).rows(parseTranslation)
	translations := make([]*Translation, len(results))
	for i, result := range results {
		if translation, ok := result.(Translation); ok {
			translations[i] = &translation
		}
	}
	return translations
}

func (entry *Entry) GetTranslations(language string) []*Translation {
	results := query("select "+translationFields+" from Translations where EntryID = ? and Language = ?", entry.ID(), language).rows(parseTranslation)
	translations := make([]*Translation, len(results))
	for i, result := range results {
		if translation, ok := result.(Translation); ok {
			translations[i] = &translation
		}
	}
	return translations
}

func (entry *Entry) GetTranslationBy(language, translator string) *Translation {
	result := query("select "+translationFields+" from Translations where EntryID = ? and Language = ? and Translator = ?", entry.ID(), language, translator).row(parseTranslation)
	if translation, ok := result.(Translation); ok {
		return &translation
	}
	return nil
}

func (entry *Entry) GetMatchingTranslation(language, translation string) *Translation {
	result := query("select "+translationFields+" from Translations where EntryID = ? and Language = ? and Translation = ?", entry.ID(), language, translation).row(parseTranslation)
	if translation, ok := result.(Translation); ok {
		return &translation
	}
	return nil
}

func (translation *Translation) Save() {
	keyfields := map[string]interface{}{
		"TranslationID": translation.ID(),
	}
	fields := map[string]interface{}{
		"Language":    translation.Language,
		"Translator":  translation.Translator,
		"Translation": translation.Translation,
		"IsPreferred": translation.IsPreferred,
	}
	saveRecord("Translations", keyfields, fields)
	ClearVotes(translation)
}

// ** Votes

type Vote struct {
	Translation Translation
	Voter       *User
	Vote        bool
}

const voteFields = "TranslationID, Voter, Vote"

func parseVote(rows *sql.Rows) (Result, error) {
	v := Vote{}
	var translationID, voter string
	err := rows.Scan(&translationID, &voter, &v.Vote)
	if err != nil {
		return nil, err
	}

	if translation := GetTranslationByID(translationID); translation == nil {
		return nil, nil
	} else {
		v.Translation = *translation
	}

	v.Voter = GetUserByEmail(voter)
	return v, err
}

func (translation *Translation) GetVote(voter *User) *Vote {
	result := query("select " + voteFields + " from Votes").row(parseVote)
	if vote, ok := result.(Vote); ok {
		vote.Translation = *translation
		vote.Voter = voter
		return &vote
	}
	return nil
}

func (entry *Entry) GetTranslationVotes(language string) []*Vote {
	results := query("select "+voteFields+" from Votes where EntryID = ? and Language = ?", entry.ID(), language).rows(parseVote)
	votes := make([]*Vote, len(results))
	for i, result := range results {
		if vote, ok := result.(Vote); ok {
			votes[i] = &vote
		}
	}
	return votes
}

func (vote *Vote) Save() {
	keyfields := map[string]interface{}{
		"EntryID":    vote.Translation.Entry.ID(),
		"Language":   vote.Translation.Language,
		"Translator": vote.Translation.Translator,
		"Voter":      vote.Voter.Email,
	}
	fields := map[string]interface{}{
		"Vote": vote.Vote,
	}
	saveRecord("Votes", keyfields, fields)
}

func DeleteVote(vote *Vote) {
	keyfields := map[string]interface{}{
		"EntryID":    vote.Translation.Entry.ID(),
		"Language":   vote.Translation.Language,
		"Translator": vote.Translation.Translator,
		"Voter":      vote.Voter.Email,
	}
	deleteRecord("Votes", keyfields)
}

func ClearVotes(translation *Translation) {
	keyfields := map[string]interface{}{
		"EntryID":    translation.Entry.ID(),
		"Language":   translation.Language,
		"Translator": translation.Translator,
	}
	deleteRecord("Votes", keyfields)
}

func ClearOtherVotes(translation *Translation) {
	keyfields := map[string]interface{}{
		"EntryID":  translation.Entry.ID(),
		"Language": translation.Language,
		"Vote":     true,
	}
	deleteRecord("Votes", keyfields)
}

/*
func AddTranslation(entry *Entry, language, translation string, translator *User) {
	keyfields := map[string]interface{}{
		"EntryID":   entry.ID(),
		"Language":      language,
		"Translator":    translator.Email,
	}
	fields := map[string]interface{}{
		"Translation": translation,
	}
	saveRecord("Translations", keyfields, fields)
}*/

// ** Users

type User struct {
	Email          string
	Password       string
	Secret         string
	Name           string
	IsAdmin        bool
	Language       string
	IsLanguageLead bool
}

func parseUser(rows *sql.Rows) (Result, error) {
	u := User{}
	err := rows.Scan(&u.Email, &u.Password, &u.Secret, &u.Name, &u.IsAdmin, &u.Language, &u.IsLanguageLead)
	return u, err
}

const userFields = "Email, Password, Secret, Name, IsAdmin, Language, IsLanguageLead"

func GetUsers() []*User {
	results := query("select " + userFields + " from Users order by IsAdmin desc, Language asc, Name asc").rows(parseUser)
	users := make([]*User, len(results))
	for i, result := range results {
		if user, ok := result.(User); ok {
			users[i] = &user
		}
	}
	return users
}

func GetUserByEmail(email string) *User {
	result := query("select "+userFields+" from Users where Email = ?", email).row(parseUser)
	if user, ok := result.(User); ok {
		return &user
	}
	return nil
}

func GetUsersByLanguage(language string) []*User {
	results := query("select "+userFields+" from Users where Language = ? order by IsLanguageLead desc, Name asc", language).rows(parseUser)
	users := make([]*User, len(results))
	for i, result := range results {
		if user, ok := result.(User); ok {
			users[i] = &user
		}
	}
	return users
}

func GetLanguageLead(language string) *User {
	result := query("select Email, "+userFields+" from Users where Language = ? and IsLanguageLead = 1", language).row(parseUser)
	if result != nil {
		if user, ok := result.(User); ok {
			return &user
		}
	}

	// users := GetUsersByLanguage(language)
	// if len(users) > 0 {
	// 	return users[0]
	// }
	return nil
}

func (user *User) Save() bool {
	keyfields := map[string]interface{}{
		"Email": user.Email,
	}
	fields := map[string]interface{}{
		"Password":       user.Password,
		"Secret":         user.Secret,
		"Name":           user.Name,
		"IsAdmin":        user.IsAdmin,
		"Language":       user.Language,
		"IsLanguageLead": user.IsLanguageLead,
	}
	return saveRecord("Users", keyfields, fields)
}

func (user *User) Delete() {
	keyfields := map[string]interface{}{
		"Email": user.Email,
	}
	deleteRecord("Users", keyfields)
}

func (user *User) CountTranslations() map[string]int {
	counts := make(map[string]int, len(Languages))
	query("select Language, Count(*) from Translations where Translator = ? group by Language", user.Email).rows(func(rows *sql.Rows) (Result, error) {
		var language string
		var count int
		rows.Scan(&language, &count)
		counts[language] = count
		return nil, nil
	})
	return counts
}

// ** Comments

type Comment struct {
	Entry       Entry
	Language    string
	Commenter   string
	Comment     string
	CommentDate mysql.Timestamp
}
