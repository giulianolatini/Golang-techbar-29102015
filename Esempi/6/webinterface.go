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

/* UX Interface to check files to dalete
 */
func hello(res http.ResponseWriter, req *http.Request) {

	Info.Println("Start hello function")
	Info.Println(res)
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
	io.WriteString(res,
		fmt.Sprintf(
			"<div>"+
				"<button onclick=myFunction(\"https://localhost/api/delete/%s\")>DELETE</button>"+
				"<a href=http://www.univpm.it/Entra/Engine/RAServeFile.php/f%s> %s </a>"+
				"</div>",
			"5",
			"/Prova",
			"Prova"))

}

func deletefile(w http.ResponseWriter, r *http.Request) {

	Info.Println("Start delete function")
	Info.Println(r)
	if r.Method == "POST" {
		Info.Println("Start POST Submit")
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err == nil {
			Info.Printf("Id %d mark to delete", id)
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
	// Command Line flags definition
	debug = flag.Int("d", 0, "debugging level: 0 = Null, 8 = INFO")
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

	// Operation Modules

	Info.Printf("Define Route")
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/api/delete/{id}", deletefile)
	Info.Printf("Start httpd deamon on http://localhost:9000")
	http.ListenAndServe(":9000", nil)

}
