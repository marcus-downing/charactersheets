package model

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	// "time"
)

// Languages
var Languages []string = []string{
	"gb", "it", "fr", "de", "es", "pt", "us",
}

var LanguageNames map[string]string = map[string]string{
	"gb": "English",
	"it": "Italiano",
	"fr": "Français",
	"de": "Deutch",
	"es": "Español",
	"pt": "Português",
	"us": "US English",
}

// data

var CurrentLanguage = "en"

//  DB access

func LoadEntries() {

}

//  import / export

func ImportCSV() {

}

func ExportCSV() {

}

//  completion

func GetLanguageCompletion() map[string][4]int {
	var completion = make(map[string][4]int, len(Languages))
	for _, lang := range Languages {
		if lang == "gb" || lang == "us" {
			completion[lang] = [4]int{100, 100, 100, 100}
		} else {
			completion[lang] = [4]int{90, 60, 40, 20}
		}
	}
	return completion
}

type StackedEntry struct {
	FullText string
	Entries  []*Entry
	EntrySources  []*EntrySource
	Count int
}

func GetStackedEntries() []*StackedEntry {
	entries := GetEntries()
	return stackEntries(entries)
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
		entrySources := make(map[string]*EntrySource, len(se.Entries) * 10)
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
