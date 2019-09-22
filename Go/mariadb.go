package main

import (
	"fmt"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Create the database handle, confirm driver is present
	db, err := sql.Open("mysql", "root:root@tcp(localhost)/test_db")
	if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
	defer db.Close()
	var language = "golang"
	rows, err := db.Query("select link from resourcelist where tech_id = (select id from techtools where name = ?)", language)
	
	if err != nil {
		fmt.Printf("failed to enumerate tables: %v", err)
	}
	for rows.Next() {
		var table string
		if rows.Scan(&table) == nil {
			fmt.Printf(table)
		}
	}
}