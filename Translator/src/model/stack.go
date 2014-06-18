package model

import (
	// "code.google.com/p/go.crypto/bcrypt"
	"fmt"
	// "math/rand"
	"sort"
	"strconv"
	"strings"
	// "time"
)

type StackedEntry struct {
	FullText     string
	Entries      []*Entry
	EntrySources []*EntrySource
	Count        int
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

	// put entries in order
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

func GetStackedEntries(game, level, show, language string, user *User) []*StackedEntry {
	leveln, err := strconv.Atoi(level)
	if err != nil || leveln > 4 || leveln < 1 {
		leveln = 0
	}
	entries := GetEntriesAt(game, leveln, show, language, user)
	return stackEntries(entries)
}

/* Stacked Translations */

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

type StackedTranslation struct {
	Entry       *StackedEntry
	Language    string
	Translator  string
	Parts       []*Translation
	Count       int
	IsPreferred bool
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

func (se *StackedEntry) CountTranslations() map[string]int {
	entryCounts := make([]map[string]int, len(se.Entries))
	for i, entry := range se.Entries {
		entryCounts[i] = entry.CountTranslations()
	}

	langCounts := make(map[string]int, len(Languages))
	for _, lang := range Languages {
		min := 0
		for _, counts := range entryCounts {
			count := counts[lang]
			if count < min || min == 0 {
				min = count
			}
		}
		if min > 0 {
			langCounts[lang] = min
		}
	}
	return langCounts
}

func (st *StackedTranslation) GetVotes() []*Vote {
	entry := st.Entry.Entries[0]
	var results []Result
	if entry.PartOf == "" {
		results = query("select "+voteFields+" from Votes where EntryOriginal = ? and EntryPartOf = '' and Language = ? and Translator = ?", entry.Original, st.Language, st.Translator).rows(parseVote)
	} else {
		results = query("select "+voteFields+" from Votes where EntryPartOf = ? and Language = ? and Translator = ?", entry.PartOf, st.Language, st.Translator).rows(parseVote)
	}

	votes := make([]*Vote, len(results))
	for i, result := range results {
		if vote, ok := result.(Vote); ok {
			votes[i] = &vote
		}
	}
	return votes
}

func GetPreferredTranslations(language string) []*StackedTranslation {
	lead := GetLanguageLead(language)
	var leadEmail string = ""
	if lead != nil {
		leadEmail = lead.Email
	}

	entries := stackEntries(GetEntries())
	pref := make([]*StackedTranslation, 0, len(entries))
	for _, entry := range entries {
		translations := entry.GetTranslations(language)
		selected := SelectPreferredTranslation(entry, language, translations, leadEmail)
		if selected != nil {
			pref = append(pref, selected)
		}
	}

	return pref
}

func SelectPreferredTranslation(entry *StackedEntry, language string, translations []*StackedTranslation, lead string) *StackedTranslation {
	if len(translations) == 0 {
		return nil
	}
	if len(translations) == 1 {
		return translations[0]
	}

	//  count scores for the text of a translation, so duplicates are merged
	//  votes are worth two; language lead is worth one (so it's a tie-breaker)
	scores := make(map[string]int, len(translations))

	for _, st := range translations {
		text := st.FullText()
		scores[text] = 0
		votes := st.GetVotes()
		for _, vote := range votes {
			if vote.Vote {
				scores[text] += 2
			} else {
				scores[text] -= 2
			}
		}
		if st.Translator == lead {
			scores[text]++
		}
	}

	//  get translations from people who haven't voted

	//  pick the highest score
	highestText := ""
	highestScore := 0
	for text, score := range scores {
		if score > highestScore {
			highestScore = score
			highestText = text
		}
	}
	for _, st := range translations {
		text := st.FullText()
		if text == highestText {
			return st
		}
	}
	return translations[0]
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
