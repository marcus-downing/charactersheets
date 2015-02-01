package model

import "fmt"

//  Profile translations

type TranslationProfile struct {
	Language     string
	Level        int
	LevelName    string
	TotalEntries int

	Completed               int
	Uncompleted             int
	ByMe                    int
	ByMeAlone               int
	ByOthers                int
	ByOthersAlone           int
	ByMeAndOthers           int
	ByMeAndOthersNoConflict int
	ByMeAndOthersConflict   int
	ByNobody                int

	CompletedPercent               float64
	UncompletedPercent             float64
	ByMePercent                    float64
	ByMeAlonePercent               float64
	ByOthersPercent                float64
	ByOthersAlonePercent           float64
	ByMeAndOthersPercent           float64
	ByMeAndOthersNoConflictPercent float64
	ByMeAndOthersConflictPercent   float64
	ByNobodyPercent                float64
}

func roundPercent(value float64) float64 {
	return float64(int(value*10000.0)) / 100.0
}

func ProfileTranslations(user *User) [4]*TranslationProfile {
	lang := user.Language

	var profiles [4]*TranslationProfile
	for level, _ := range LevelNames {
		if level == 0 {
			continue
		}

		fmt.Println(LanguageNames[lang], " at level ", LevelNames[level])

		total := query("select count(*) from (select Entries.EntryID from Entries "+
			"inner join EntrySources on Entries.EntryID = EntrySources.EntryID "+
			"inner join Sources on EntrySources.SourceID = Sources.SourceID and Sources.Level = ? "+
			"group by Entries.EntryID) as sq",
			level).count()
		fmt.Println(" -- total = ", total)

		completed := query("select count(*) from (select Translations.EntryID from Translations "+
			"inner join EntrySources on Translations.EntryID = EntrySources.EntryID "+
			"inner join Sources on EntrySources.SourceID = Sources.SourceID and Sources.Level = ? "+
			"where Language = ? group by Translations.EntryID) as sq",
			level, lang).count()
		fmt.Println(" -- completed = ", completed)

		byme := query("select count(*) from (select Translations.EntryID from Translations "+
			"inner join EntrySources on Translations.EntryID = EntrySources.EntryID "+
			"inner join Sources on Sources.SourceID = EntrySources.SourceID and Sources.Level = ? "+
			"where Language = ? and Translator = ? group by Translations.EntryID) as sq",
			level, lang, user.Email).count()
		fmt.Println(" -- by me = ", byme)

		byothers := query("select count(*) from (select Translations.EntryID from Translations "+
			"inner join EntrySources on Translations.EntryID = EntrySources.EntryID "+
			"inner join Sources on Sources.SourceID = EntrySources.SourceID and Sources.Level = ? "+
			"where Language = ? and Translator != ? group by Translations.EntryID) as sq",
			level, lang, user.Email).count()
		fmt.Println(" -- by others =", byothers)

		byboth := query("select count(*) from (select A.EntryID from Translations A "+
			"inner join Translations B on A.EntryID = B.EntryID and A.Language = B.Language and A.Translator != B.Translator "+
			"inner join EntrySources on A.EntryID = EntrySources.EntryID "+
			"inner join Sources on Sources.SourceID = EntrySources.SourceID and Sources.Level = ? "+
			"where A.Language = ? and A.Translator = ? "+
			"group by A.EntryID"+
			") as qs", level, lang, user.Email).count()
		fmt.Println(" -- by both =", byboth)

		conflict := query("select count(*) from (select A.EntryID from Translations A "+
			"inner join Translations B on A.EntryID = B.EntryID and A.Language = B.Language and A.Translator != B.Translator and A.Translation != B.Translation "+
			"inner join EntrySources on A.EntryID = EntrySources.EntryID "+
			"inner join Sources on Sources.SourceID = EntrySources.SourceID and Sources.Level = ? "+
			"where A.Language = ? and A.Translator = ? "+
			"group by A.EntryID"+
			") as qs", level, lang, user.Email).count()
		fmt.Println(" -- conflicting =", conflict)

		total64 := float64(total) + 0.0001
		completed64 := float64(completed) + 0.0001

		profile := TranslationProfile{
			Language:     lang,
			Level:        level,
			LevelName:    LevelNames[level],
			TotalEntries: total,

			Completed:               completed,
			Uncompleted:             total - completed,
			ByMe:                    byme,
			ByMeAlone:               byme - byboth,
			ByOthers:                byothers,
			ByOthersAlone:           byothers - byboth,
			ByMeAndOthers:           byboth,
			ByMeAndOthersNoConflict: byboth - conflict,
			ByMeAndOthersConflict:   conflict,
			ByNobody:                total - (byme + byothers - byboth),

			CompletedPercent:               roundPercent(float64(completed) / total64),
			UncompletedPercent:             roundPercent(float64(total-completed) / total64),
			ByMePercent:                    roundPercent(float64(byme-byboth) / completed64),
			ByOthersPercent:                roundPercent(float64(byothers-byboth) / completed64),
			ByMeAndOthersPercent:           roundPercent(float64(byboth) / completed64),
			ByMeAndOthersNoConflictPercent: roundPercent(float64(byboth-conflict) / completed64),
			ByMeAndOthersConflictPercent:   roundPercent(float64(conflict) / completed64),
			ByNobodyPercent:                roundPercent(float64(total-(byme+byothers-byboth)) / completed64),
		}
		fmt.Println(" -- profile =", profile)
		profiles[level-1] = &profile
	}
	return profiles
}
