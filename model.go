package main

import (
	"database/sql"
)

type employee struct {
	Uid        int    `json:"uid"` //field - user ident. number (autoincremented)
	Username   string `json:"username"`
	Departname string `json:"departname"`
}

func (p *employee) getEmployee(db *sql.DB) error {
	return db.QueryRow("SELECT username, departname FROM userinfo WHERE uid=$1",
		p.Uid).Scan(&p.Username, &p.Departname)
}

func (p *employee) updateEmployee(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE userinfo SET username=$1, departname=$2 WHERE uid=$3",
			p.Username, p.Departname, p.Uid)

	return err
}

func (p *employee) deleteEmployee(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM userinfo WHERE uid=$1", p.Uid)

	return err
}

func (p *employee) createEmployee(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO userinfo(username, departname) VALUES($1, $2) RETURNING uid",
		p.Username, p.Departname).Scan(&p.Uid)

	if err != nil {
		return err
	}

	return nil
}

func getEmployees(DB *sql.DB, start, count int) ([]employee, error) {
	rows, err := DB.Query(
		"SELECT uid, username, departname FROM userinfo LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	employees := []employee{}

	for rows.Next() {
		var p employee
		if err := rows.Scan(&p.Uid, &p.Username, &p.Departname); err != nil {
			return nil, err
		}
		employees = append(employees, p)
	}

	return employees, nil
}
