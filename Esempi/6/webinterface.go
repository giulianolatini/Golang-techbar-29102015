package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql" //Import for side effect
	_ "github.com/mattn/go-sqlite3"    //Importazione statica anche in mancato uso
)

// Global Section
var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

var localdb *sql.DB
var rootPath *string // Root Path to start filesystem visiting

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

/* walkpath funzione di callback per
la stampa di files e paths ispezionati
*/
func walkpath(path string, f os.FileInfo, err error) error {
	Info.Printf("%s with %d bytes\n", path[len(*rootPath):], f.Size())
	localdb.Exec("Insert into ostree "+
		"(id, path, filename, date) values"+
		" (?, ?, ?, ?)",
		nil, path[len(*rootPath):], f.Name(), nil)
	return nil
}

/* UX Interface to check files to dalete
 */
func hello(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	io.WriteString(
		res,
		`<doctype html>
			<html>
			<head>
				<title>Hello World</title>
			</head>
			<body>
				Hello World!
			</body>
			</html>`,
	)
	localrows, err := localdb.Query("select id,path from filesnodownload;")
	if err != nil {
		Error.Println("Failed to query for filesnodownload")
		os.Exit(30)
	}

	// Get column names
	columns, err := localrows.Columns()
	if err != nil {
		Error.Println("Failed to get columns name from filesnodownload")
		os.Exit(31) // proper error handling instead of panic in your app
	}

	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for localrows.Next() {
		// get RawBytes from data
		err = localrows.Scan(scanArgs...)
		if err != nil {
			Error.Println("Failed to parsing columsData from Download_Obj")
			os.Exit(32) // proper error handling instead of panic in your app
		}
		io.WriteString(res,
			fmt.Sprintf(
				"<div>"+
					"<button onclick=myFunction(\"https://localhost/api/delete/%s\")>DELETE</button>"+
					"<a href=http://www.univpm.it/Entra/Engine/RAServeFile.php/f%s> %s </a>"+
					"</div>",
				NullorValue(values[0]),
				NullorValue(values[1]),
				NullorValue(values[1][1:])))
	}

}

func deletefile(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err == nil {
			Info.Printf("Id %d mark to delete", id)
			localdb.Exec("Insert into todelete "+
				"(id, deleteid) values"+
				" (?, ?)",
				nil, id)
		} else {
			Error.Println("Failed to parsing Id from Validation Check")
			return // proper error handling instead of panic in your app
		}
	} else {
		return
	}
}

func main() {

	// Init Section

	type ObjDwnLoadInfo struct {
		idkey   sql.NullInt64
		fsPath  sql.NullString
		fsName  sql.NullString
		fsExt   sql.NullString
		urlPath sql.NullString
	}
	var err error       // Generic variables error
	var ptyColumns *int // Terminal Columns to Output
	var debug *int      // Debug level

	// Command Line flags definition
	rootPath = flag.String("rootPath", ".", "the root from start the filesystem's visit")
	debug = flag.Int("d", 0, "debugging level: 0 = Null, 8 = INFO")
	ptyColumns = flag.Int("columns", 80, "large columns for format output")

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
	}

	//	for _, e := range os.Environ() {
	//		pair := strings.Split(e, "=")
	//		fmt.Printf("%s = %s \n", pair[0], pair[1])
	//	}

	// Operation Modules

	entradb, entraerr := sql.Open("mysql", "Entra:w78KIJ10R@/EntraNIA")
	Info.Println("Open E-ntra MySQL")
	if entraerr != nil {
		Error.Println("Can't connect MySQL")
		//panic(entraerr.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	localdb, err = sql.Open("sqlite3", os.TempDir()+"/localdb.db")
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

		// Now do something with the data.
		// Here we just print each column as a string.
		//value := make([]string, 1)
		//for i, col := range values {
		//	// Here we can check if the value is nil (NULL value)
		//	if col == nil {
		//		value = append(value, "NULL")
		//	} else {
		//		value = append(value, string(col))
		//	}
		//	fmt.Println(columns[i], ": ", value)
		//}
		//fmt.Println("-----------------------------------")
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

	//readerDir, err := ioutil.ReadDir(*rootPath)

	//for _, fileInfo := range readerDir {

	//	fmt.Printf("\n File Name : %s \n", fileInfo.Name())
	//	fmt.Printf("\n File Is Directory? : %v \n", fileInfo.IsDir())
	//	fmt.Printf("\n File Size : %d \n", fileInfo.Size())
	//	fmt.Printf("\n File Last Modified Time : %s \n", fileInfo.ModTime())
	//	fmt.Printf("\n File Permission : %s \n", fileInfo.Mode())
	//	fmt.Println("----------------")
	//}

	//fmt.Println("Error : ", err)

	//root := flag.Arg(0) // 1st argument is the directory location
	filepath.Walk(*rootPath, walkpath)

	result, localerr = localdb.Exec("CREATE VIEW IF NOT EXISTS filesnodownload AS " +
		"SELECT ostree.id,path,filename FROM ostree " +
		"LEFT OUTER JOIN downloadobj " +
		"ON ostree.filename = downloadobj.FS_Name " +
		"WHERE downloadobj.id IS null;")

	if localerr != nil {
		Error.Println("Failed to create filesnodownload View")
		os.Exit(6) // proper error handling instead of panic in your app
	} else {
		Info.Println("Creation filesnodownload View OK:", result)
	}

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/api/delete/{id}", deletefile)
	http.ListenAndServe(":9000", nil)

}
