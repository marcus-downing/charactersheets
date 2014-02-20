package model

import (
	// "crypto/md5"
	"database/sql"
	// "encoding/hex"
	"github.com/ziutek/mymysql/mysql"
	"fmt"
)

// ** Entries

type Entry struct {
	Original string
	PartOf   string
}

func parseEntry(rows *sql.Rows) (Result, error) {
	e := Entry{}
	err := rows.Scan(&e.Original, &e.PartOf)
	return e, err
}

func CountEntries() int {
	return query("select count(*) from Entries").count()
}

func GetEntries() []*Entry {
	results := query("select Original, PartOf from Entries").rows(parseEntry)

	entries := make([]*Entry, len(results))
	for i, result := range results {
		if entry, ok := result.(Entry); ok {
			entries[i] = &entry
		}
	}
	return entries
}

func GetEntriesAt(game string, level int, show, language string) []*Entry {
	if game == "" && level == 0 && show == "" {
		return GetEntries()
	}
	args := make([]interface{}, 0, 2)
	sql := "select Original, PartOf from Entries "+
		"inner join EntrySources on Original = EntrySources.EntryOriginal and PartOf = EntrySources.EntryPartOf "+
		"inner join Sources on SourcePath = Filepath"
	if show != "" {
		sql = sql+" left join Translations on Original = Translations.EntryOriginal and PartOf = Translations.EntryPartOf and Translations.Language = ?"
		args = append(args, language)
	}
	sql = sql+" where 1 = 1"

	if game != "" {
		sql = sql+" and Game = ?"
		args = append(args, game)
	}
	if level != 0 {
		sql = sql+" and Level = ?"
		args = append(args, level)
	}
	// if show != "" {
	// 	sql = sql+" and Translations.Language = ?"
	// 	args = append(args, language)
	// }

	sql = sql+" group by Original, PartOf"
	if show == "translated" {
		sql = sql+" having count(Translations.Translation) > 0"
	} else if show == "untranslated" {
		sql = sql+" having count(Translations.Translation) = 0"
	}
	fmt.Println("Get entries:", sql)
	results := query(sql, args...).rows(parseEntry)

	entries := make([]*Entry, 0, len(results))
	for _, result := range results {
		if entry, ok := result.(Entry); ok {
			entries = append(entries, &entry)
		}
	}
	return entries
}

func (entry *Entry) Save() {
	keyfields := map[string]interface{}{
		"Original": entry.Original,
		"PartOf":   entry.PartOf,
	}
	fields := map[string]interface{}{}
	saveRecord("Entries", keyfields, fields)
}

func (entry *Entry) CountTranslations() map[string]int {
	counts := make(map[string]int, len(Languages))
	query("select Language, Count(*) from Translations where EntryOriginal = ? and EntryPartOf = ? group by Language", entry.Original, entry.PartOf).rows(func (rows *sql.Rows) (Result, error) {
		var language string
		var count int
		rows.Scan(&language, &count)
		counts[language] = count
		return nil, nil
	})
	return counts
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

func GetSources() []*Source {
	results := query("select Filepath, Page, Volume, Level, Game from Sources").rows(parseSource)

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
	sql := "select Filepath, Page, Volume, Level, Game from Sources "
		// "inner join EntrySources on Original = EntrySources.EntryOriginal and PartOf = EntrySources.EntryPartOf "+
		// "inner join Sources on SourcePath = Filepath"
	// if show != "" {
	// 	sql = sql+" left join Translations on Original = Translations.EntryOriginal and PartOf = Translations.EntryPartOf"
	// }
	sql = sql+" where 1 = 1"

	if game != "" {
		sql = sql+" and Game = ?"
		args = append(args, game)
	}
	if level != 0 {
		sql = sql+" and Level = ?"
		args = append(args, level)
	}
	// if show != "" {
	// 	sql = sql+" and Translations.Language = ?"
	// 	args = append(args, language)
	// }

	sql = sql+" group by Original, PartOf"
	// if show == "translated" {
	// 	sql = sql+" having count(Translations.Translation) > 0"
	// } else if show == "untranslated" {
	// 	sql = sql+" having count(Translations.Translation) = 0"
	// }
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

type EntrySource struct {
	Entry  Entry
	Source Source
	Count  int
}

func parseEntrySource(rows *sql.Rows) (Result, error) {
	es := EntrySource{}
	err := rows.Scan(&es.Entry.Original, &es.Entry.PartOf, &es.Source.Filepath, &es.Source.Page, &es.Source.Volume, &es.Source.Level, &es.Source.Game, &es.Count)
	return es, err
}

func GetEntrySources() []*EntrySource {
	results := query("select EntryOriginal, EntryPartOf, SourcePath, Sources.Page, Sources.Volume, Sources.Level, Sources.Game, Count " +
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
	results := query("select EntryOriginal, EntryPartOf, SourcePath, Page, Volume, Level, Game, Count "+
		"from EntrySources inner join Sources on SourcePath = Filepath "+
		"where EntryOriginal = ? and EntryPartOf = ?", entry.Original, entry.PartOf).rows(parseEntrySource)

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
		"EntryOriginal": es.Entry.Original,
		"EntryPartOf":   es.Entry.PartOf,
		"SourcePath":    es.Source.Filepath,
	}
	fields := map[string]interface{}{
		"Count": es.Count,
	}
	saveRecord("EntrySources", keyfields, fields)
}


// ** Translations

type Translation struct {
	Entry       Entry
	Language    string
	Translation string
	Translator  string
	IsPreferred bool
}

func parseTranslation(rows *sql.Rows) (Result, error) {
	t := Translation{}
	err := rows.Scan(&t.Entry.Original, &t.Entry.PartOf, &t.Language, &t.Translation, &t.Translator, &t.IsPreferred)
	return t, err
}

func GetTranslations() []*Translation {
	results := query("select EntryOriginal, EntryPartOf, Language, Translation, Translator, IsPreferred from Translations").rows(parseTranslation)
	translations := make([]*Translation, len(results))
	for i, result := range results {
		if translation, ok := result.(Translation); ok {
			translations[i] = &translation
		}
	}
	return translations
}

func GetTranslationsForLanguage(language string) []*Translation {
	results := query("select EntryOriginal, EntryPartOf, Language, Translation, Translator, IsPreferred from Translations where Language = ?", language).rows(parseTranslation)
	translations := make([]*Translation, len(results))
	for i, result := range results {
		if translation, ok := result.(Translation); ok {
			translations[i] = &translation
		}
	}
	return translations
}

func (entry *Entry) GetTranslations(language string) []*Translation {
	results := query("select EntryOriginal, EntryPartOf, Language, Translation, Translator, IsPreferred from Translations where EntryOriginal = ? and EntryPartOf = ? and Language = ?", entry.Original, entry.PartOf, language).rows(parseTranslation)
	translations := make([]*Translation, len(results))
	for i, result := range results {
		if translation, ok := result.(Translation); ok {
			translations[i] = &translation
		}
	}
	return translations
}

func (entry *Entry) GetTranslationBy(language, translator string) *Translation {
	result := query("select EntryOriginal, EntryPartOf, Language, Translation, Translator, IsPreferred from Translations where EntryOriginal = ? and EntryPartOf = ? and Language = ? and Translator = ?", entry.Original, entry.PartOf, language, translator).row(parseTranslation)
	if translation, ok := result.(Translation); ok {
		return &translation
	}
	return nil
}

func (translation *Translation) Save() {
	keyfields := map[string]interface{}{
		"EntryOriginal": translation.Entry.Original,
		"EntryPartOf":   translation.Entry.PartOf,
		"Language":      translation.Language,
		"Translator":    translation.Translator,
	}
	fields := map[string]interface{}{
		"Translation": translation.Translation,
	}
	saveRecord("Translations", keyfields, fields)
}

/*
func AddTranslation(entry *Entry, language, translation string, translator *User) {
	keyfields := map[string]interface{}{
		"EntryOriginal": entry.Original,
		"EntryPartOf":   entry.PartOf,
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

func GetUsers() []*User {
	results := query("select Email, Password, Secret, Name, IsAdmin, Language, IsLanguageLead from Users order by IsAdmin desc, Language asc, Name asc").rows(parseUser)
	users := make([]*User, len(results))
	for i, result := range results {
		if user, ok := result.(User); ok {
			users[i] = &user
		}
	}
	return users
}

func GetUserByEmail(email string) *User {
	result := query("select Email, Password, Secret, Name, IsAdmin, Language, IsLanguageLead from Users where Email = ?", email).row(parseUser)
	if user, ok := result.(User); ok {
		return &user
	}
	return nil
}

func GetUsersByLanguage(language string) []*User {
	results := query("select Email, Password, Secret, Name, IsAdmin, Language, IsLanguageLead from Users where Language = ? order by IsLanguageLead desc, Name asc", language).rows(parseUser)
	users := make([]*User, len(results))
	for i, result := range results {
		if user, ok := result.(User); ok {
			users[i] = &user
		}
	}
	return users
}

func GetLanguageLead(language string) *User {
	result := query("select Email, Password, Secret, Name, IsAdmin, Language, IsLanguageLead from Users where Language = ? and IsLanguageLead = 1", language).row(parseUser)
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
		"Password": user.Password,
		"Secret":   user.Secret,
		"Name":     user.Name,
		"IsAdmin":  user.IsAdmin,
		"Language": user.Language,
		"IsLanguageLead": user.IsLanguageLead,
	}
	return saveRecord("Users", keyfields, fields)
}

func (user *User) CountTranslations() map[string]int {
	counts := make(map[string]int, len(Languages))
	query("select Language, Count(*) from Translations where Translator = ? group by Language", user.Email).rows(func (rows *sql.Rows) (Result, error) {
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
