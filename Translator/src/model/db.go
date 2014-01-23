package model

import (
	"database/sql"
	"fmt"
	"strings"
	// _ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/godrv"
)

var dbname = "chartrans"
var dbuser = "chartrans"
var dbpassword = "fiddlesticks"

//  users

func withDB(f func(db *sql.DB) interface{}) interface{} {
	// connect to database
	db, err := sql.Open("mymysql", dbname+"/"+dbuser+"/"+dbpassword)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil
	}
	defer db.Close()

	return f(db)
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
	db, err := sql.Open("mymysql", dbname+"/"+dbuser+"/"+dbpassword)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return false
	}
	defer db.Close()

	fmt.Println("Exists: ", query.sql, query.args)
	rows, err := db.Query(query.sql, query.args...)
	if err != nil {
		fmt.Println("Exists: error:", err)
		return false
	}
	return rows.Next()
}

func (query Query) exec(success func(lastInsertId int64)) bool {
	db, err := sql.Open("mymysql", dbname+"/"+dbuser+"/"+dbpassword)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return false
	}
	defer db.Close()

	result, err := db.Exec(query.sql, query.args...)
	if err == nil {
		fmt.Println("Executed successfully", result)
		if success != nil {
			lastInsertId, _ := result.LastInsertId()
			success(lastInsertId)
		}
		return true
	} else {
		fmt.Println("Error executing:", err)
	}
	return false
}

func (query Query) rows(f func(*sql.Rows) (Result, error)) []Result {
	// connect to database
	db, err := sql.Open("mymysql", dbname+"/"+dbuser+"/"+dbpassword)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil
	}
	defer db.Close()

	rows, err := db.Query(query.sql, query.args...)
	if err != nil {
		fmt.Println("Error querying database:", err)
		return nil
	}

	var results []Result = make([]Result, 0, 100)
	for rows.Next() {
		result, err := f(rows)
		if err != nil {
			fmt.Println("Error parsing row:", err)
		} else if result != nil {
			results = append(results, result)
		}
	}
	// fmt.Println("Found", len(results), "results")
	return results
}

func (query Query) row(f func(*sql.Rows) (Result, error)) Result {
	// connect to database
	db, err := sql.Open("mymysql", dbname+"/"+dbuser+"/"+dbpassword)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil
	}
	defer db.Close()

	rows, err := db.Query(query.sql, query.args...)
	if err != nil {
		fmt.Println("Error querying database:", err)
		return nil
	}

	if rows.Next() {
		result, err := f(rows)
		if err != nil {
			fmt.Println("Error decoding row:", err)
		} else if result != nil {
			return result
		}
	}
	return nil
}

func recordExists(table string, keyfields map[string]interface{}) bool {
	conditions := make([]string, 0, len(keyfields))
	args := make([]interface{}, 0, len(keyfields))
	for key, value := range keyfields {
		conditions = append(conditions, key+" = ?")
		args = append(args, value)
	}
	sql := "select 1 from " + table + " where " + strings.Join(conditions, " and ")
	fmt.Println("Checking ", table, ":", sql, args)
	return query(sql, args...).exists()
}

func saveRecord(table string, keyfields, fields map[string]interface{}, success func(lastInsertId int64)) bool {
	fmt.Println("Saving record")

	update := recordExists(table, keyfields)

	var sql string
	args := make([]interface{}, 0, len(keyfields)+len(fields))
	if update {
		fmt.Println("Record exists, updating")
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
		fmt.Println("Record doesn't exist, inserting")
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

	return query(sql, args...).exec(success)
}
