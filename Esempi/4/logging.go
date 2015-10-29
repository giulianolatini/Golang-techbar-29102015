package main

import (
	"database/sql"
	"fmt"
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

func main() {

	// Init Section

	var err error        // Generic variables error
	var ptyColumns *int  // Terminal Columns to Output
	var debug *int       // Debug level
	var rootPath *string // Root Path to start filesystem visiting

	// Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)
	// Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stdout)
	// Init(ioutil.Discard, ioutil.Discard, os.Stdout, os.Stdout)
	// Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stdout)
	Init(os.Stdout, os.Stdout, os.Stdout, os.Stdout)

	// Operation Modules

	entradb, entraerr := sql.Open("mysql", "Entra:w78KIJ10R@/EntraNIA")
	Info.Println("Open E-ntra MySQL")
	if entraerr != nil {
		Error.Println("Can't connect MySQL")
		//panic(entraerr.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	localdb, err := sql.Open("sqlite3", os.TempDir()+"/localdb.db")
	Info.Println("Open local SQLite")
	if err != nil {
		Error.Println("Can't open local DB")
		//panic(localerr.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	defer entradb.Close()
	defer localdb.Close()

	if entraerr = entradb.Ping(); entraerr != nil {
		Error.Println("Failed to keep connection alive")
		os.Exit(2)
		//panic(entraerr.Error()) // proper error handling instead of panic in your app
	}

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
		os.Exit(6) // proper error handling instead of panic in your app
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

	// Execute the query
	//rows, err := db.Query("SELECT T13_Id_Obj FROM EntraNIA.Oggetti " +
	//                      "WHERE (T13_Id_SottoSito='001' " +
	//                      "AND T13_StatoPagina='0' AND T13_Id_Oggetto='13');")
	entrarows, entraerr := entradb.Query("SELECT T25_Dir, T25_Name, T25_Ext, T25_Path " +
		"FROM EntraNIA.Download " +
		"INNER JOIN EntraNIA.Download_Obj " +
		"ON EntraNIA.Download.T25_Id_Obj=EntraNIA.Download_Obj.T13_id_Obj;")

	if entraerr != nil {
		Error.Println("Failed to query for Download_Obj")
		os.Exit(8)
	}

	// Get column names
	columns, err := entrarows.Columns()
	if err != nil {
		Error.Println("Failed to get columns name from Download_Obj")
		os.Exit(10) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	iRow := 0
	for entrarows.Next() {
		// get RawBytes from data
		err = entrarows.Scan(scanArgs...)
		if err != nil {
			Error.Println("Failed to parsing columsData from Download_Obj")
			os.Exit(12) // proper error handling instead of panic in your app
		}

		result, localerr = localdb.Exec("Insert into DownloadObj "+
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

			Info.Printf("Creation row: %d in DownloadObj OK", iRow)
		}
	}
	if err = entrarows.Err(); err != nil {
		Error.Printf("Failed on row: %d for corruptions data", iRow)
		os.Exit(16) // proper error handling instead of panic in your app
	}
}
