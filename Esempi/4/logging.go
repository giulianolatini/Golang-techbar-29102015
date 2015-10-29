package main

import (
	"database/sql"
	"io"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

// Global Section
var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

/*
Init : Inizialize Logs Handlers
*/
func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

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

	var err error // Generic variables error

	// Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)
	// Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stdout)
	// Init(ioutil.Discard, ioutil.Discard, os.Stdout, os.Stdout)
	// Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stdout)
	Init(os.Stdout, os.Stdout, os.Stdout, os.Stdout)

	// Operation Modules

	localdb, err := sql.Open("sqlite3", os.TempDir()+"/localdb.db")
	Info.Println("Open local SQLite")
	if err != nil {
		Error.Println("Can't open local DB")
		//panic(localerr.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	defer localdb.Close()

	if err = localdb.Ping(); err != nil {
		Error.Println("Failed to keep connection alive")
		os.Exit(4)
		//panic(localerr.Error()) // proper error handling instead of panic in your app
	}

	result, localerr := localdb.Exec("CREATE TABLE IF NOT EXISTS DownloadObj (" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL," +
		"FS_Path TEXT," +
		"FS_Name TEXT," +
		"FS_Ext TEXT," +
		"URL_Path TEXT" +
		"); CREATE INDEX IF NOT EXISTS PK_WEBID on downloadobj (id ASC);")
	Info.Println(result)
	if localerr != nil {
		Error.Println("Failed to create DownloadObj Table and Index")
		os.Exit(6) // proper error handling instead of panic in your app
	} else {
		Info.Println("Creation DonloadObj OK:", result)
	}
	Trace.Println(result)
	result, localerr = localdb.Exec("CREATE TABLE IF NOT EXISTS ostree (" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL," +
		"path TEXT," +
		"filename TEXT," +
		"date TEXT " +
		"); CREATE INDEX IF NOT EXISTS PK_OSID on ostree (id ASC);")
	Trace.Println(result)
	if localerr != nil {
		Error.Println("Failed to create OsTree Table and Index")
		os.Exit(6) // proper error handling instead of panic in your app
	} else {
		Info.Println("Creation OsTree OK:", result)
	}

}
