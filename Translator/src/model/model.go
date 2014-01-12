package model

import (
	"math/rand"
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
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
	Entries []*Entry
	Sources []*Source
}

func GetStackedEntries() []*StackedEntry {
	entries := GetEntries()
	return stackEntries(entries)
}

func stackEntries(entries []*Entry) []*StackedEntry {
	fmt.Println("Stacking", len(entries), "entries")
	stacks := make(map[string][]*Entry, len(entries))
	unstacked := make([]*Entry, 0, len(entries))
	for _, entry := range entries {
		if entry.PartOf != "" {
			if stacks[entry.PartOf] == nil {
				fmt.Println("Creating stack:", entry.PartOf)
				stacks[entry.PartOf] = make([]*Entry, 0, 10)
			}
			fmt.Println("Adding to stack:", entry.PartOf)
			stacks[entry.PartOf] = append(stacks[entry.PartOf], entry)
		} else {
			unstacked = append(unstacked, entry)
		}
	}

	fmt.Println("Making", len(stacks), "and", len(unstacked), "stacks")
	values := make([]*StackedEntry, 0, len(stacks)+len(unstacked))
	i := 0
	for _, stack := range stacks {
		values[i] = &StackedEntry{
			FullText: stack[0].PartOf,
			Entries: stack,
		}
		i++
	}
	for _, entry := range unstacked {
		values = append(values, &StackedEntry{
			FullText: entry.Original,
			Entries: []*Entry{entry},
		})
	}
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
	}
	secret := string(bytes)
	secret = secret[7:]
	fmt.Println("Generate secret:", secret)

	hash, err := bcrypt.GenerateFromPassword([]byte(secret), 12)
	if err != nil {
		fmt.Println("Error generating secret hash:", err)
	}
	fmt.Println("Generate secret: hash:", string(hash))
	user.Secret = string(hash)
	user.Save()

	// time.Sleep(5000 * time.Millisecond)

	return secret
}