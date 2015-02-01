package model

var liveStrings []string = []string{
	"Acrobatics", "Appraise", "Bluff", "Climb", "Diplomacy", "Disable Device", "Disguise",
	"Escape Artist", "Fly", "Handle Animal", "Heal", "Intimidate", "Linguistics", "Perception",
	"Ride", "Sense Motive", "Sleight of Hand", "Spellcraft", "Stealth", "Survival", "Swim",
	"Use Magic Device",
}

func liveEntries() []*Entry {
	entries := make([]*Entry, 0, len(liveStrings))
	for _, str := range liveStrings {
		entry := Entry{str, str}
		entries = append(entries, &entry)
	}
	return entries
}

func GetLiveTranslations() []*StackedTranslation {
	entries := liveEntries()

	translations := make([]*StackedTranslation, 0, len(entries)*len(Languages))
	return translations
}
