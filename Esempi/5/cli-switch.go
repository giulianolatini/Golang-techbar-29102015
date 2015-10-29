package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

var rootPath *string // Root Path to start filesystem visiting
var debug *int       // Debug level

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

func main() {

	// Init Section

	// Command Line flags definition
	rootPath = flag.String("rootPath", ".", "the root from start the filesystem's visit")
	debug = flag.Int("d", 0, "debugging level: 0 = Null, 8 = INFO")
	ptyColumns := flag.Int("columns", 80, "large columns for format output")
	button := flag.Bool("button", false, "Print Enviroment Variabile: on(true)/off(false)")
	flag.Parse()

	switch *debug {
	case 0: //Null Output
		{
			Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)
		}
	case 1: //Error channel enabled
		{
			Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stdout)
		}
	case 2: //Error and Warning channel enabled
		{
			Init(ioutil.Discard, ioutil.Discard, os.Stdout, os.Stdout)
		}
	case 4: //Error and Warning channel enabled
		{
			Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stdout)
		}
	case 8: //Error and Warning channel enabled
		{
			Init(os.Stdout, os.Stdout, os.Stdout, os.Stdout)
		}
	default: //Null Output
		{
			Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)
		}
	}

	Trace.Printf("Value of rootPath: %s", *rootPath)
	Trace.Printf("Value of ptyColumns: %d", *ptyColumns)
	Trace.Printf("Value of debug: %d", *debug)
	if *button {
		Trace.Println("Value of button: TRUE")
	} else {
		Trace.Println("Value of button: FALSE")
	}

	if *button {
		for _, e := range os.Environ() {
			pair := strings.Split(e, "=")
			fmt.Printf("%s = %s \n", pair[0], pair[1])
		}
	}

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

	if localerr != nil {
		Error.Println("Failed to create DownloadObj Table and Index")
		os.Exit(5) // proper error handling instead of panic in your app
	} else {
		Info.Println("Creation DonloadObj OK:", result)
	}

	result, localerr = localdb.Exec("CREATE TABLE IF NOT EXISTS ostree (" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL," +
		"path TEXT," +
		"filename TEXT," +
		"date TEXT " +
		"); CREATE INDEX IF NOT EXISTS PK_OSID on ostree (id ASC);")

	if localerr != nil {
		Error.Println("Failed to create OsTree Table and Index")
		os.Exit(6) // proper error handling instead of panic in your app
	} else {
		Info.Println("Creation OsTree OK:", result)
	}

	result, localerr = localdb.Exec("CREATE TABLE IF NOT EXISTS todelete (" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL," +
		"deleteid INTEGER); CREATE INDEX IF NOT EXISTS PK_TDID on todelete (id ASC);")

	if localerr != nil {
		Error.Println("Failed to create ToDelete Table and Index")
		os.Exit(7) // proper error handling instead of panic in your app
	} else {
		Info.Println("Creation ToDelete OK:", result)
	}
}
