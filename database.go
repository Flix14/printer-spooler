package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

//Adaptar
func dbExample() {
	database, _ := sql.Open("sqlite3", "./files.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS files (id INTEGER PRIMARY KEY, name TEXT, status TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO files (name, status) VALUES (?, ?)")
	statement.Exec("file1", "printed")
	rows, _ := database.Query("SELECT id, name, status FROM files")
	var id int
	var name string
	var status string
	for rows.Next() {
		rows.Scan(&id, &name, &status)
		fmt.Println(strconv.Itoa(id) + ": " + name + " " + status)
	}
}
