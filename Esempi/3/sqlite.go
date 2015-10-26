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

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	result, err = db.Exec("CREATE TABLE IF NOT EXISTS ostree (" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL," +
		"path TEXT," +
		"filename TEXT," +
		"date TEXT " +
		"); CREATE INDEX IF NOT EXISTS PK_OSID on ostree (id ASC);")

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for entrarows.Next() {
		// get RawBytes from data
		err = entrarows.Scan(scanArgs...)
		if err != nil {
			Error.Println("Failed to parsing columsData from Download_Obj")
			os.Exit(12) // proper error handling instead of panic in your app
		}

		result, err = db.Exec("Insert into DownloadObj "+
			"(id, FS_Path, FS_Name, FS_Ext, URL_Path) values"+
			" (?, ?, ?, ?, ?)",
			nil,
			NullorValue(values[0]), NullorValue(values[1]),
			NullorValue(values[2]), NullorValue(values[3]))
		if localerr != nil {
			Error.Println("Failed to create DownloadObj Table and Index")
			os.Exit(14) // proper error handling instead of panic in your app
		} else {
			iRow++
			switch *debug {
			case 0, 2:
				{ // Print . for Insert
					fmt.Printf("%s", ".")
					if (iRow % *ptyColumns) == 0 {
						fmt.Println(".")
					}
				}
			}
		}
	}
}
