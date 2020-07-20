package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

func existsTable(db *sql.DB, tableName string) bool {
	if db == nil {
		return false
	}

	stmtOut, err := db.Query("SELECT squareNumber FROM " + tableName)
	if err != nil {
		return false
	}
	defer stmtOut.Close()
	return true
}

func createSquareTable(db *sql.DB) bool {
	stmtOut, err := db.Query(fmt.Sprint(`
		CREATE TABLE SQUARENUMBERS(
		number integer PRIMARY KEY,
		squareNumber integer NOT NULL)
	`))
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stmtOut.Close()
	return true
}

func dropSquareTable(db *sql.DB) bool {
	stmtOut, err := db.Query(fmt.Sprint("DROP TABLE SQUARENUMBERS"))
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stmtOut.Close()
	return true
}

var maxRows = 20

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/test")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println("Connected!!")
	if existsTable(db, "SQUARENUMBERS") {
		dropSquareTable(db)
	}
	createSquareTable(db)
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO SQUARENUMBERS VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT squareNumber FROM SQUARENUMBERS WHERE number = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	// Insert square numbers for 0-24 in the database
	for i := 0; i < maxRows; i++ {
		_, err = stmtIns.Exec(i, i*i) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	var squareNum int // we "scan" the result in here

	// Query the square-number of 13
	err = stmtOut.QueryRow(13).Scan(&squareNum) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 13 is: %d\n", squareNum)
	fmt.Println()
	fmt.Println("Table data")
	fmt.Println()
	// Query another number.. 1 maybe?
	err = stmtOut.QueryRow(1).Scan(&squareNum) // WHERE number = 1
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 1 is: %d\n", squareNum)

	selectTable(db, "SQUARENUMBERS", "")

	selectTable(db, "", "show tables")

	if existsTable(db, "SQUARENUMBERS") {
		dropSquareTable(db)
	}
}

func selectTable(db *sql.DB, name string, sqlQuery string) {
	var stmtQryOut *sql.Rows
	var err error
	var entity string
	if sqlQuery == "" {
		entity = "table: " + name
		stmtQryOut, err = db.Query(fmt.Sprintf("SELECT number, squareNumber FROM %s", name))
	} else {
		entity = "query: " + sqlQuery
		stmtQryOut, err = db.Query(sqlQuery)
	}
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer func() {
		_ = stmtQryOut.Close()
	}()
	fmt.Println()
	fmt.Printf("Report for %s\n", entity)
	fmt.Println()
	length, columns, lineLen, lineSep := formatHeader(stmtQryOut, err)
	// Make a slice for the values
	var rowNum int64
	for stmtQryOut.Next() {
		rowNum++
		formatRow(rowNum, columns, stmtQryOut, length, lineLen, lineSep, err)
	}
	if rowNum == 0 {
		fmt.Println("No rows selected.")
	}
	fmt.Println()
	fmt.Println()
}

func formatHeader(stmtQryOut *sql.Rows, err error) (length int, columns int, lineLen int, lineSep string) {
	if cols, err := stmtQryOut.Columns(); err == nil {
		columns = len(cols)
		for _, col := range cols {
			if len(col) > length {
				length = len(col)
			}
		}
		length += 2
		lineLen = 7 + columns*(length+1)
		lineSep = strings.Repeat("-", lineLen)
		fmt.Println(lineSep)
		fmt.Print("|")
		fmt.Print(pad("#", 5))
		fmt.Print("|")
		for _, col := range cols {
			fmt.Print(pad(col, length) + "|")
		}
		fmt.Println()
		fmt.Println(lineSep)
	} else {
		fmt.Println(err.Error())
	}
	return length, columns, lineLen, lineSep
}

func formatRow(rowNum int64, columns int, rows *sql.Rows, length, lineLen int, lineSep string, err error) {
	values := make([]sql.RawBytes, columns)
	scanArgs := make([]interface{}, len(values))
	for i, _ := range values {
		scanArgs[i] = &values[i]
	}
	err = rows.Scan(scanArgs...)
	if err == nil {
		fmt.Print("|")
		fmt.Print(pad(fmt.Sprintf("%v", rowNum), 5))
		fmt.Print("|")
		for _, col := range values {
			var value string
			if col == nil {
				value = "NULL"
			} else {
				str := string(col)
				value = str
			}
			fmt.Print(pad(value, length) + "|")
		}
		fmt.Println()
		fmt.Println(lineSep)
	} else {
		fmt.Print("|")
		fmt.Print(pad(fmt.Sprintf("%v", rowNum), 5))
		fmt.Print("|")
		fmt.Println(pad(err.Error(), lineLen-8))
		fmt.Print("|")
		fmt.Println()
		fmt.Println(lineSep)
	}
}

func pad(s string, l int) string {
	if len(s) < l {
		return strings.Repeat(" ", l-len(s)) + s
	} else if len(s) > l {
		return s[:l]
	}
	return s

}
