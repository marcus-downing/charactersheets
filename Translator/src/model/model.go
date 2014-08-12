package model

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"math/rand"
	// "sort"
	// "strconv"
	"strings"
	// "time"
)

// Languages
var Languages []string = []string{
	"gb", "it", "fr", "de", "es", "pl", "pt", "ru", "cy", "kl", "us",
}

var DisplayLanguages []string = []string{
	"it", "de", "es", "fr", "pl", "pt", "ru",
}

var LanguageNames map[string]string = map[string]string{
	"gb": "English",
	"it": "Italiano",
	"fr": "Français",
	"de": "Deutch",
	"es": "Español",
	"pl": "Polski",
	"pt": "Português",
	"ru": "Ру́сский",
	"cy": "Cymraeg",
	"kl": "Klingon",
	"us": "US English",
}

var LanguagePaths map[string]string = map[string]string{
	"it": "italian",
	"fr": "french",
	"de": "german",
	"es": "spanish",
	"pl": "polish",
	"pt": "portuguese",
	"ru": "russian",
	"cy": "welsh",
	"kl": "klingon",
	"us": "american",
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

func WhereNotMe(in chan string) <-chan string {
	out := make(chan string)
	for s := range in {
		if s != "me" {
			out <- s
		}
	}
	return out
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

//  Profile translations

type TranslationProfile struct {
	Language     string
	TotalEntries int

	ByMe                    int
	ByMeAlone               int
	ByOthers                int
	ByOthersAlone           int
	ByMeAndOthers           int
	ByMeAndOthersNoConflict int
	ByMeAndOthersConflict   int
	ByNobody                int

	ByMePercent                    int
	ByMeAlonePercent               int
	ByOthersPercent                int
	ByOthersAlonePercent           int
	ByMeAndOthersPercent           int
	ByMeAndOthersNoConflictPercent int
	ByMeAndOthersConflictPercent   int
	ByNobodyPercent                int
}

func ProfileTranslations(user *User) map[string]*TranslationProfile {
	total := CountEntries()
	profiles := make(map[string]*TranslationProfile, len(DisplayLanguages))

	for _, lang := range DisplayLanguages {
		if lang == "gb" {
			continue
		}
		byme := query("select count(*) from (select count(*) from Translations where Language = ? and Translator = ? group by EntryOriginal, EntryPartOf) as sq", lang, user.Email).count()
		if byme > 0 || user.IsAdmin {
			byothers := query("select count(*) from (select count(*) from Translations where Language = ? and Translator != ? group by EntryOriginal, EntryPartOf) as sq", lang, user.Email).count()
			byboth := query("select count(*) from Translations A "+
				"inner join Translations B on A.EntryOriginal = B.EntryOriginal and A.EntryPartOf = B.EntryPartOf and A.Language = B.Language "+
				"where A.Language = ? and A.Translator = ? and B.Translator != ? "+
				"group by A.EntryOriginal, A.EntryPartOf"+
				"", lang, user.Email, user.Email).count()
			conflict := query("select count(*) from Translations A "+
				"inner join Translations B on A.EntryOriginal = B.EntryOriginal and A.EntryPartOf = B.EntryPartOf and A.Language = B.Language "+
				"where A.Language = ? and A.Translator = ? and B.Translator != ? and A.Translation != B.Translation "+
				"group by A.EntryOriginal, A.EntryPartOf"+
				"", lang, user.Email, user.Email).count()

			fmt.Println(LanguageNames[lang], "-- by me = ", byme, "; by others =", byothers, "; by both =", byboth)
			profile := TranslationProfile{
				Language:                lang,
				TotalEntries:            total,
				ByMe:                    byme,
				ByMeAlone:               byme - byboth,
				ByOthers:                byothers,
				ByOthersAlone:           byothers - byboth,
				ByMeAndOthers:           byboth,
				ByMeAndOthersNoConflict: conflict,
				ByMeAndOthersConflict:   byboth - conflict,
				ByNobody:                total - (byme + byothers - byboth),

				ByMePercent:                    (100 * (byme - byboth)) / total,
				ByOthersPercent:                (100 * (byothers - byboth)) / total,
				ByMeAndOthersPercent:           (100 * byboth) / total,
				ByMeAndOthersNoConflictPercent: (100 * (byboth - conflict)) / total,
				ByMeAndOthersConflictPercent:   (100 * conflict) / total,
				ByNobodyPercent:                (100 * (total - (byme + byothers - byboth))) / total,
			}
			profiles[lang] = &profile
		}
	}
	return profiles
}
