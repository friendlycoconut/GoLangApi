package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
	a  App
)

func main() {
	a = App{}
	a.Initialize()

	a.Run(":8080")
}
