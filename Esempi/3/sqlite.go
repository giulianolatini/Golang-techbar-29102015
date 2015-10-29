package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

/*
NullorValue restituisce la stringa "NULL" se v è vuoto
oppure il valore v in tutti gli altri casi
*/
func NullorValue(v sql.RawBytes) (s string) {

	var result string
	if v != nil {
		result = string(v)
	} else {
		result = "NULL"
	}
	return result
}

func main() {

	// Init Section

	// Operation Modules

	db, err := sql.Open("sqlite3", os.TempDir()+"/localdb.db")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	result, err := db.Exec("CREATE TABLE IF NOT EXISTS DownloadObj (" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL," +
		"FS_Path TEXT," +
		"FS_Name TEXT," +
		"FS_Ext TEXT," +
		"URL_Path TEXT" +
		"); CREATE INDEX IF NOT EXISTS PK_WEBID on downloadobj (id ASC);")
	fmt.Println(result)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	result, err = db.Exec("CREATE TABLE IF NOT EXISTS ostree (" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL," +
		"path TEXT," +
		"filename TEXT," +
		"date TEXT " +
		"); CREATE INDEX IF NOT EXISTS PK_OSID on ostree (id ASC);")
	fmt.Println(result)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

}
