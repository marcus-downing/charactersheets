package model

import (
	// "crypto/md5"
	"database/sql"
	// "encoding/hex"
	"github.com/ziutek/mymysql/mysql"
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

func (entry *Entry) Save() {
	keyfields := map[string]interface{}{
		"Original": entry.Original,
		"PartOf":   entry.PartOf,
	}
	fields := map[string]interface{}{}
	saveRecord("Entries", keyfields, fields, nil)
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

func (source *Source) Save() {
	keyfields := map[string]interface{}{
		"Filepath": source.Filepath,
	}
	fields := map[string]interface{}{
		"Page": source.Page,
		"Volume": source.Volume,
		"Level": source.Level,
		"Game": source.Game,
	}
	saveRecord("Sources", keyfields, fields, nil)
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
	results := query("select EntryOriginal, EntryPartOf, SourcePath, Sources.Page, Sources.Volume, Sources.Level, Sources.Game, Count "+
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
		"SourcePath": es.Source.Filepath,
	}
	fields := map[string]interface{}{
		"Count": es.Count,
	}
	saveRecord("EntrySources", keyfields, fields, nil)
}

// ** Translations

type Translation struct {
	Entry       Entry
	Language    string
	Translation string
	Translator  string
}

func parseTranslation(rows *sql.Rows) (Result, error) {
	t := Translation{}
	err := rows.Scan(&t.Entry.Original, &t.Entry.PartOf, &t.Language, &t.Translation, &t.Translator)
	return t, err
}

func GetTranslations() []*Translation {
	results := query("select EntryOriginal, EntryPartOf, Language, Translation, Translator from Translations").rows(parseTranslation)
	translations := make([]*Translation, len(results))
	for i, result := range results {
		if translation, ok := result.(Translation); ok {
			translations[i] = &translation
		}
	}
	return translations
}

func GetPartTranslations(original, partOf, language string) []*Translation {
	results := query("select EntryOriginal, EntryPartOf, Language, Translation, Translator from Translations where EntryOriginal = ? and EntryPartOf = ? and Language = ?", original, partOf, language).rows(parseTranslation)
	translations := make([]*Translation, len(results))
	for i, result := range results {
		if translation, ok := result.(Translation); ok {
			translations[i] = &translation
		}
	}
	return translations
}

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
	saveRecord("Translations", keyfields, fields, nil)
}

// ** Users

type User struct {
	Email    string
	Password string
	Secret   string
	Name     string
	IsAdmin  bool
	Language string
}

func parseUser(rows *sql.Rows) (Result, error) {
	u := User{}
	err := rows.Scan(&u.Email, &u.Password, &u.Secret, &u.Name, &u.IsAdmin, &u.Language)
	return u, err
}

func GetUsers() []*User {
	results := query("select Email, Password, Secret, Name, IsAdmin, Language from Users order by IsAdmin desc, Language asc, Name asc").rows(parseUser)
	users := make([]*User, len(results))
	for i, result := range results {
		if user, ok := result.(User); ok {
			users[i] = &user
		}
	}
	return users
}

func GetUserByEmail(email string) *User {
	result := query("select Email, Password, Secret, Name, IsAdmin, Language from Users where Email = ?", email).row(parseUser)
	if user, ok := result.(User); ok {
		return &user
	}
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
	}
	return saveRecord("Users", keyfields, fields, nil)
}

// ** Comments

type Comment struct {
	Entry       Entry
	Language    string
	Commenter   string
	Comment     string
	CommentDate mysql.Timestamp
}
