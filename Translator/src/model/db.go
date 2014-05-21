package model

import (
	"database/sql"
	"fmt"
	"strings"
	// _ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/godrv"
)

const (
	dbname     = "chartrans"
	dbuser     = "chartrans"
	dbpassword = "fiddlesticks"
)

var Debug = 0

//  users

var db, err = sql.Open("mymysql", dbname+"/"+dbuser+"/"+dbpassword)

/*
func WithDB(f func(db *sql.DB)) {
	// connect to database
	db, err := sql.Open("mymysql", dbname+"/"+dbuser+"/"+dbpassword)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	f(db)
	db.Close()
}
*/

func WithDB(f func(tx *sql.Tx)) {
	// connect to database
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	f(tx)
	tx.Commit()
}

type Query struct {
	sql  string
	args []interface{}
}

type Result interface {
}

func query(query string, args ...interface{}) Query {
	return Query{query, args}
}

func (query Query) exists() bool {
	if Debug >= 2 {
		fmt.Println("Exists:", query.sql, query.args)
	}
	exists := false
	WithDB(func(db *sql.Tx) {
		rows, err := db.Query(query.sql, query.args...)
		if err != nil {
			fmt.Println("Exists: error:", err)
			return
		}
		exists = rows.Next()
		rows.Close()
	})
	return exists
}

func (query Query) count() int {
	if Debug >= 2 {
		fmt.Println("Count:", query.sql, query.args)
	}
	count := 0
	WithDB(func(db *sql.Tx) {
		rows, err := db.Query(query.sql, query.args...)
		if err != nil {
			fmt.Println("Count: error:", err)
			return
		}
		if rows.Next() {
			rows.Scan(&count)
		}
		rows.Close()
	})
	return count
}

func (query Query) exec() bool {
	if Debug >= 2 {
		fmt.Println("Exec:", query.sql, query.args)
	}
	success := false
	WithDB(func(db *sql.Tx) {
		_, err := db.Exec(query.sql, query.args...)
		if err != nil {
			fmt.Println("Error executing:", err)
			return
		}
		success = true
	})
	return success
}

func (query Query) rows(f func(*sql.Rows) (Result, error)) []Result {
	if Debug >= 2 {
		fmt.Println("Query:", query.sql, query.args)
	}
	results := make([]Result, 0, 100)
	WithDB(func(db *sql.Tx) {
		rows, err := db.Query(query.sql, query.args...)
		if err != nil {
			fmt.Println("Error querying database:", err)
			return
		}

		for rows.Next() {
			result, err := f(rows)
			if err != nil {
				fmt.Println("Error parsing row:", err)
			} else if result != nil {
				results = append(results, result)
			}
		}
		rows.Close()
	})
	if Debug >= 2 {
		fmt.Println("Found", len(results), "results")
	}
	return results
}

func (query Query) row(f func(*sql.Rows) (Result, error)) Result {
	if Debug >= 2 {
		fmt.Println("Query:", query.sql, query.args)
	}
	var result Result = nil
	WithDB(func(db *sql.Tx) {
		rows, err := db.Query(query.sql, query.args...)
		if err != nil {
			fmt.Println("Error querying database:", err)
			return
		}

		if rows.Next() {
			result, err = f(rows)
			if err != nil {
				fmt.Println("Error decoding row:", err)
			}
		}
		rows.Close()
	})
	return result
}

func recordExists(table string, keyfields map[string]interface{}) bool {
	conditions := make([]string, 0, len(keyfields))
	args := make([]interface{}, 0, len(keyfields))
	for key, value := range keyfields {
		conditions = append(conditions, key+" = ?")
		args = append(args, value)
	}
	sql := "select 1 from " + table + " where " + strings.Join(conditions, " and ")
	if Debug >= 2 {
		fmt.Println("Checking ", table, ":", sql, args)
	}
	return query(sql, args...).exists()
}

func saveRecord(table string, keyfields, fields map[string]interface{}) bool {
	if Debug >= 2 {
		fmt.Println("Saving record")
	}

	update := recordExists(table, keyfields)

	var sql string
	args := make([]interface{}, 0, len(keyfields)+len(fields))
	if update {
		if len(fields) == 0 {
			if Debug >= 2 {
				fmt.Println("Record exists, skipping")
			}
			return true
		}

		if Debug >= 2 {
			fmt.Println("Record exists, updating")
		}
		names := make([]string, 0, len(fields))
		for key, value := range fields {
			names = append(names, key+" = ?")
			args = append(args, value)
		}
		conditions := make([]string, 0, len(keyfields))
		for key, value := range keyfields {
			conditions = append(conditions, key+" = ?")
			args = append(args, value)
		}

		sql = "update " + table + " set " + strings.Join(names, ", ") + " where " + strings.Join(conditions, " and ")
	} else {
		if Debug >= 2 {
			fmt.Println("Record doesn't exist, inserting")
		}
		names := make([]string, 0, len(keyfields)+len(fields))
		qs := make([]string, 0, len(keyfields)+len(fields))
		for key, value := range keyfields {
			names = append(names, key)
			qs = append(qs, "?")
			args = append(args, value)
		}
		for key, value := range fields {
			names = append(names, key)
			qs = append(qs, "?")
			args = append(args, value)
		}
		sql = "insert into " + table + " (" + strings.Join(names, ", ") + ") values (" + strings.Join(qs, ", ") + ")"
	}

	return query(sql, args...).exec()
}

func deleteRecord(table string, keyfields map[string]interface{}) {
	conditions := make([]string, 0, len(keyfields))
	args := make([]interface{}, 0, len(keyfields))
	for key, value := range keyfields {
		conditions = append(conditions, key+" = ?")
		args = append(args, value)
	}
	sql := "delete from " + table + " where " + strings.Join(conditions, " and ")
	if Debug >= 2 {
		fmt.Println("Deleting ", table, ":", sql, args)
	}
	query(sql, args...).exec()
}
