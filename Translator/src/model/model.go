package model

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"strconv"
	// "time"
)

// Languages
var Languages []string = []string{
	"gb", "it", "fr", "de", "es", "pt", "cy", "kl", "us",
}

var LanguageNames map[string]string = map[string]string{
	"gb": "English",
	"it": "Italiano",
	"fr": "Français",
	"de": "Deutch",
	"es": "Español",
	"pt": "Português",
	"cy": "Cymraeg",
	"kl": "Klingon",
	"us": "US English",
}


//  completion

func GetLanguageCompletion() map[string][4]int {
	var completion = make(map[string][4]int, len(Languages))
	var totals [4]int
	for i := 1; i <= 4; i++ {
		totals[i-1] = query("select count(distinct Original, PartOf) from Entries "+
			"inner join EntrySources on Original = EntryOriginal and PartOf = EntryPartOf "+
			"inner join Sources on SourcePath = Filepath "+
			"where Level = ?", i).count()
	}

	for _, lang := range Languages {
		if lang == "gb" {
			completion[lang] = [4]int{100, 100, 100, 100}
		} else {
			var values [4]int
			for i := 1; i <= 4; i++ {
				count := query("select count(distinct Translations.EntryOriginal, Translations.EntryPartOf) from Translations "+
					"inner join EntrySources on Translations.EntryOriginal = EntrySources.EntryOriginal and Translations.EntryPartOf = EntrySources.EntryPartOf "+
					"inner join Sources on SourcePath = Filepath "+
					"where Level = ? and Language = ?", i, lang).count()
				if totals[i-1] > 0 {
					values[i-1] = 100 * count / totals[i-1]
				}
				fmt.Println("Completion of", LanguageNames[lang], "@", i, "=", count, "/", totals[i-1])
			}
			completion[lang] = values
		}
	}
	return completion
}

type StackedEntry struct {
	FullText     string
	Entries      []*Entry
	EntrySources []*EntrySource
	Count        int
}

func GetStackedEntries(game, level, show, language string) []*StackedEntry {
	leveln, err := strconv.Atoi(level)
	if err != nil || leveln > 4 || leveln < 1 {
		leveln = 0
	}
	entries := GetEntriesAt(game, leveln, show, language)
	return stackEntries(entries)
}

func (se *StackedEntry) GetTranslations(language string) []*StackedTranslation {
	length := len(se.Entries)
	translations := make(map[string][]*Translation, 30)

	for _, entry := range se.Entries {
		entryTranslations := entry.GetTranslations(language)
		for _, translation := range entryTranslations {
			if _, ok := translations[translation.Translator]; !ok {
				translations[translation.Translator] = make([]*Translation, 0, length)
			}
			translations[translation.Translator] = append(translations[translation.Translator], translation)
		}
	}

	stackedTranslations := make([]*StackedTranslation, 0, len(translations))
	for translator, parts := range translations {
		stacked := StackedTranslation{
			Entry:      se,
			Language:   language,
			Translator: translator,
			Parts:      parts,
		}
		if !stacked.Empty() {
			stackedTranslations = append(stackedTranslations, &stacked)
		}
	}
	return stackedTranslations
}

func (se *StackedEntry) GetTranslationBy(language, translator string) *StackedTranslation {
	parts := make([]*Translation, len(se.Entries))
	for i, entry := range se.Entries {
		parts[i] = entry.GetTranslationBy(language, translator)
		if parts[i] == nil {
			parts[i] = &Translation{
				Entry:       *entry,
				Language:    language,
				Translation: "",
				Translator:  translator,
			}
		}
	}
	stacked := StackedTranslation{
		Entry:      se,
		Language:   language,
		Translator: translator,
		Parts:      parts,
		Count:      len(parts),
	}
	return &stacked
}

func WhereNotMe(in chan string) <- chan string {
	out := make(chan string)
	for s := range in {
		if s != "me" {
			out <- s
		}
	}
	return out
}

type StackedTranslation struct {
	Entry      *StackedEntry
	Language   string
	Translator string
	Parts      []*Translation
	Count      int
}

func (st *StackedTranslation) Empty() bool {
	for _, part := range st.Parts {
		if part != nil && strings.TrimSpace(part.Translation) != "" {
			return false
		}
	}
	return true
}

func (st *StackedTranslation) FullText() string {
	text := make([]string, len(st.Parts))
	for i, part := range st.Parts {
		text[i] = part.Translation
	}
	return strings.Join(text, "")
}

// sort entries by index
type entriesByIndex []*Entry

func (this entriesByIndex) Len() int {
	return len(this)
}
func (this entriesByIndex) Less(i, j int) bool {
	return strings.Index(this[i].PartOf, this[i].Original) < strings.Index(this[j].PartOf, this[j].Original)
}
func (this entriesByIndex) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// sort stacked entries by name
type stacksByName []*StackedEntry

func (this stacksByName) Len() int {
	return len(this)
}

func (this stacksByName) Less(i, j int) bool {
	return this[i].FullText < this[j].FullText
}

func (this stacksByName) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// sort stacked entries by number of uses
type stacksByCount []*StackedEntry

func (this stacksByCount) Len() int {
	return len(this)
}

func (this stacksByCount) Less(i, j int) bool {
	return this[i].Count > this[j].Count
}

func (this stacksByCount) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func stackEntries(entries []*Entry) []*StackedEntry {
	fmt.Println("Stacking", len(entries), "entries")
	stacks := make(map[string][]*Entry, len(entries))
	unstacked := make([]*Entry, 0, len(entries))
	for _, entry := range entries {
		if entry.PartOf != "" {
			if stacks[entry.PartOf] == nil {
				stacks[entry.PartOf] = make([]*Entry, 0, 10)
			}
			stacks[entry.PartOf] = append(stacks[entry.PartOf], entry)
		} else {
			unstacked = append(unstacked, entry)
		}
	}

	//
	values := make([]*StackedEntry, 0, len(stacks)+len(unstacked))
	for _, stack := range stacks {
		sort.Sort(entriesByIndex(stack))
		values = append(values, &StackedEntry{
			FullText: stack[0].PartOf,
			Entries:  stack,
		})
	}
	for _, entry := range unstacked {
		values = append(values, &StackedEntry{
			FullText: entry.Original,
			Entries:  []*Entry{entry},
		})
	}

	// calculate totals
	for _, se := range values {
		entrySources := make(map[string]*EntrySource, len(se.Entries)*10)
		for _, entry := range se.Entries {
			for _, es := range GetSourcesForEntry(entry) {
				entrySources[es.Source.Filepath] = es
			}
		}
		count := 0
		esv := make([]*EntrySource, 0, len(entrySources))
		for _, es := range entrySources {
			esv = append(esv, es)
			count += es.Count
		}
		se.EntrySources = esv
		se.Count = count
	}
	sort.Sort(stacksByName(values))
	sort.Sort(stacksByCount(values))
	return values
}

//

func (user *User) GenerateSecret() string {
	base := make([]byte, 256)
	for i, _ := range base {
		base[i] = byte(rand.Int())
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(base), bcrypt.MinCost)
	if err != nil {
		fmt.Println("Error generating secret hash:", err)
		return ""
	}
	secret := string(bytes)
	// secret = secret[7:]
	fmt.Println("Generate secret:", secret)

	hash, err := bcrypt.GenerateFromPassword([]byte(secret), 12)
	if err != nil {
		fmt.Println("Error generating secret hash:", err)
	}
	fmt.Println("Generate secret: hash:", string(hash))
	user.Secret = string(hash)
	user.Save()

	return secret
}

func (user *User) SetLanguageLead() {
	query("update Users set IsLanguageLead = 0 where Language = ?", user.Language).exec()
	query("update Users set IsLanguageLead = 1 where Email = ?", user.Email).exec()
}

func (user *User) ClearLanguageLead() {
	query("update Users set IsLanguageLead = 0 where Email = ?", user.Email).exec()
}
