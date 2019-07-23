package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"

	"time"

)



func main() {


	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_userusername)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	fmt.Println("# Inserting values")




	var lastInsertId int
	err = db.QueryRow("INSERT INTO userinfo(useruserusername,depart,created) VALUES($1,$2,$3) returning uid;", "ILKO", "Computer Science", "2012-12-09").Scan(&lastInsertId)
	checkErr(err)
	fmt.Println("last inserted id =", lastInsertId)

	fmt.Println("# Updating")
	stmt, err := db.Prepare("update userinfo set useruserusername=$1 where uid=$2")
	checkErr(err)

	res, err := stmt.Exec("ILKO", lastInsertId)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect, "rows changed")

	fmt.Println("# Querying")
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var uid int
		var useruserusername string
		var departname string
		var created time.Time
		err = rows.Scan(&uid, &useruserusername, &departname, &created)
		checkErr(err)
		fmt.Println("uid | useruserusername | departname | created ")
		fmt.Printf("%3v | %8v | %6v | %6v\n", uid, useruserusername, departname, created)
	}

	/*fmt.Println("# Deleting")
	stmt, err = db.Prepare("delete from userinfo1 where uid=$1")
	checkErr(err)

	res, err = stmt.Exec(lastInsertId)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)*/

	fmt.Println(affect, "rows changed")
}

//

