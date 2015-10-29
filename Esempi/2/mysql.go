package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := sql.Open("mysql", "Entra:w78KIJ10R@/EntraNIA")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Execute the query
	//rows, err := db.Query("SELECT T13_Id_Obj FROM EntraNIA.Oggetti " +
	//                      "WHERE (T13_Id_SottoSito='001' " +
	//                      "AND T13_StatoPagina='0' AND T13_Id_Oggetto='13');")
	rows, err := db.Query(`SELECT T25_Dir, T25_Name, T25_Ext, T25_Path
		FROM EntraNIA.Download
		INNER JOIN EntraNIA.Download_Obj
		ON EntraNIA.Download.T25_Id_Obj=EntraNIA.Download_Obj.T13_id_Obj;`)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
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
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		value := make([]string, 1)
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = append(value, "NULL")
			} else {
				value = append(value, string(col))
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
}
