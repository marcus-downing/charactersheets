package model

import (
	"github.com/ziutek/mymysql/mysql"
	"database/sql"
	"crypto/md5"
	"encoding/hex"
)

type DBRow struct {
	Loaded bool
}


// ** Master Entries

type Entry struct {
	DBRow
	MD5 string
	Original string
	PartOf   string
}

// 32-hex MD5 hash of the Original and PartOf fields
func entryMD5(original, partOf string) string {
	hasher := md5.New()
	hasher.Write([]byte(original))
	hasher.Write([]byte(partOf))
	return hex.EncodeToString(hasher.Sum(nil))
}

func parseEntry(rows *sql.Rows) (Result, error) {
	e := Entry{}
	err := rows.Scan(&e.Original, &e.PartOf)
	e.MD5 = entryMD5(e.Original, e.PartOf)
	return e, err
}

func GetEntries() []*Entry {
	results := query("select Original, PartOf from entries").rows(parseEntry)

	entries := make([]*Entry, len(results))
	for i, result := range results {
		if entry, ok := result.(Entry); ok {
			entries[i] = &entry
		}
	}
	return entries
}


// ** Sources

type EntrySource struct {
	DBRow
	Entry  int
	Source int
}

type Source struct {
	DBRow
	ID     int
	Page   string
	Volume string
	Group  string
	Game   string
}

func parseSource(rows *sql.Rows) (Result, error) {
	s := Source{}
	err := rows.Scan(&s.ID, &s.Page, &s.Volume, &s.Group, &s.Game)
	return s, err
}
/*
func GetSources() []*Source {

}*/


// ** Translations

type Translation struct {
	DBRow
	Entry       string
	Language    string
	Translation string
	Translator  string
}

func parseTranslation(rows *sql.Rows) (Result, error) {
	t := Translation{}
	err := rows.Scan(&t.Entry, &t.Language, &t.Translation, &t.Translator)
	return t, err
}

func GetTranslations() []*Translation {
	results := query("select Entry, Language, Translation, Translator from Translations").rows(parseTranslation)
	translations := make([]*Translation, len(results))
	for i, result := range results {
		if translation, ok := result.(Translation); ok {
			translations[i] = &translation
		}
	}
	return translations
}

func GetPartTranslations(original, partOf, language string) []*Translation {
	md5 := entryMD5(original, partOf)
	results := query("select Entry, Language, Translation, Translator from Translations where Entry = ? and Language = ?", md5, language).rows(parseTranslation)
	translations := make([]*Translation, len(results))
	for i, result := range results {
		if translation, ok := result.(Translation); ok {
			translations[i] = &translation
		}
	}
	return translations
}

func AddTranslation(md5, language, translation string, translator *User) {
	keyfields := map[string]interface{}{
		"Entry": md5,
		"Language": language,
		"Translator": translator.Email,
	}
	fields := map[string]interface{}{
		"Translation": translation,
	}
	loaded := false
	saveRecord("translations", keyfields, fields, loaded, nil);
}

type Comment struct {
	DBRow
	Entry       string
	Language    string
	Commenter   string
	Comment     string
	CommentDate mysql.Timestamp
}


// ** Users

type User struct {
	DBRow
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
	results := query("select Email, Password, Secret, Name, IsAdmin, Language from users order by IsAdmin desc, Language asc, Name asc").rows(parseUser)
	users := make([]*User, len(results))
	for i, result := range results {
		if user, ok := result.(User); ok {
			users[i] = &user
		}
	}
	return users
}

func GetUserByEmail(email string) *User {
	result := query("select Email, Password, Secret, Name, IsAdmin, Language from users order by IsAdmin desc, Language asc, Name asc").row(parseUser)
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
		"Secret": user.Secret,
		"Name": user.Name,
		"IsAdmin": user.IsAdmin,
		"Language": user.Language,
	}
	return saveRecord("users", keyfields, fields, user.Loaded, nil)
}