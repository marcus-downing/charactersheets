package main

import (
	"./model"
	"database/sql"
	"fmt"
	_ "github.com/ziutek/mymysql/godrv"
)

func main() {
	dbname1 := "chartrans"
	dbname2 := "chartrans2"
	dbuser := "chartrans"
	dbpassword := "fiddlesticks"

	db1, err := sql.Open("mymysql", dbname1+"/"+dbuser+"/"+dbpassword)
	if err != nil {
		fmt.Println("Error opening database 1:", err)
	}
	db2, err := sql.Open("mymysql", dbname2+"/"+dbuser+"/"+dbpassword)
	if err != nil {
		fmt.Println("Error opening database 2:", err)
	}

	// clear out
	fmt.Print("Clearing out old data... ")
	_, _ = db2.Exec("delete from Entries")
	_, _ = db2.Exec("delete from Sources")
	_, _ = db2.Exec("delete from EntrySources")
	_, _ = db2.Exec("delete from Translations")
	_, _ = db2.Exec("delete from Users")
	_, _ = db2.Exec("delete from Votes")
	fmt.Println("done.")

	// entries
	fmt.Print("Entries... ")
	rows, err := db1.Query("select Original, PartOf from Entries")
	if err != nil {
		fmt.Println("Error reading entries:", err)
	} else {
		n_entries := 0
		for rows.Next() {
			entry := model.Entry{}
			err = rows.Scan(&entry.Original, &entry.PartOf)
			if err != nil {
				fmt.Println("Error reading entry:", err)
				continue
			}
			_, err = db2.Exec("insert into Entries(EntryID, Original, PartOf) values (?,?,?)", entry.ID(), entry.Original, entry.PartOf)
			if err != nil {
				fmt.Println("Error writing entry:", err)
				continue
			}
			n_entries++
		}
		fmt.Println(n_entries)
	}
	rows.Close()

	// sources
	fmt.Print("Sources... ")
	result, err := db2.Exec("insert into Sources (Filepath, Page, Volume, Level, Game) select Filepath, Page, Volume, Level, Game from " + dbname1 + ".Sources")
	if err != nil {
		fmt.Println("Error transferring sources:", err)
	} else {
		n_sources, _ := result.RowsAffected()
		fmt.Println(n_sources)
	}

	fmt.Print("Source lines... ")
	rows, err = db1.Query("select EntryOriginal, EntryPartOf, SourcePath, Count from EntrySources")
	if err != nil {
		fmt.Println("Error reading source lines:", err)
	} else {
		n_es := 0
		for rows.Next() {
			es := model.EntrySource{}
			err = rows.Scan(&es.Entry.Original, &es.Entry.PartOf, &es.Source.Filepath, &es.Count)
			if err != nil {
				fmt.Println("Error reading source line:", err)
			}
			_, err = db2.Exec("insert into EntrySources(EntryID, SourcePath, Count) values (?,?,?)", es.Entry.ID(), es.Source.Filepath, es.Count)
			if err != nil {
				fmt.Println("Error writing source line:", err)
				continue
			}
			n_es++
		}
		fmt.Println(n_es)
	}
	rows.Close()

	// translations
	fmt.Print("Translations... ")
	rows, err = db1.Query("select EntryOriginal, EntryPartOf, Language, Translator, Translation, IsPreferred from Translations")
	if err != nil {
		fmt.Println("Error reading translations:", err)
	} else {
		n_translations := 0
		for rows.Next() {
			translation := model.Translation{}
			err = rows.Scan(&translation.Entry.Original, &translation.Entry.PartOf, &translation.Language, &translation.Translator, &translation.Translation, &translation.IsPreferred)
			if err != nil {
				fmt.Println("Error reading translation:", err)
				continue
			}
			_, err = db2.Exec("insert into Translations(TranslationID, EntryID, Language, Translator, Translation, IsPreferred, IsConflicted) values (?,?,?,?,?,?,?)",
				translation.ID(), translation.Entry.ID(), translation.Language, translation.Translator, translation.Translation, translation.IsPreferred, false)
			n_translations++
		}
		fmt.Println(n_translations)
	}
	rows.Close()

	// users
	fmt.Print("Users... ")
	result, err = db2.Exec("insert into Users (Email, Password, Secret, Name, IsAdmin, Language, IsLanguageLead) select Email, Password, Secret, Name, IsAdmin, Language, IsLanguageLead from " + dbname1 + ".Users")
	if err != nil {
		fmt.Println("Error transferring users:", err)
	} else {
		n_users, _ := result.RowsAffected()
		fmt.Println(n_users)
	}

	// votes
	fmt.Print("Votes...")
	rows, err = db1.Query("select EntryOriginal, EntryPartOf, Language, Translator, Voter, Vote from Votes")
	if err != nil {
		fmt.Println("Error reading votes:", err)
	} else {
		n_votes := 0
		for rows.Next() {
			vote := model.Vote{}
			var voter string
			err = rows.Scan(&vote.Translation.Entry.Original, &vote.Translation.Entry.PartOf, &vote.Translation.Language, &vote.Translation.Translator, &voter, &vote.Vote)
			if err != nil {
				fmt.Println("Error reading vote:", err)
				continue
			}
			_, err = db2.Exec("insert into Votes (TranslationID, Voter, Vote) values (?,?,?)", vote.Translation.ID(), voter, vote.Vote)
			if err != nil {
				fmt.Println("Error writing vote:", err)
				continue
			}
			n_votes++
		}
		fmt.Println(n_votes)
	}
}
