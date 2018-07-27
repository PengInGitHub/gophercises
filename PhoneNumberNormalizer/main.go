package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "pengchengliu"
	password = "thepassword"
	dbname   = "gophercise_phone"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	// db, err := sql.Open("postgres", psqlInfo)
	// handleErr(err)

	// //create DB
	// err = resetDB(db, dbname)
	// handleErr(err)
	// db.Close()

	//disconnect DB, create tables
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	handleErr(err)
	defer db.Close()
	handleErr(createPhoneNumbersTable(db))
	_, err = insertPhoneNumber(db, "1234567890")
	handleErr(err)
	_, err = insertPhoneNumber(db, "123 456 7891")
	handleErr(err)
	_, err = insertPhoneNumber(db, "(123) 456 7892")
	handleErr(err)
	_, err = insertPhoneNumber(db, "(123) 456-7893")
	handleErr(err)
	_, err = insertPhoneNumber(db, "123-456-7890")
	handleErr(err)
	_, err = insertPhoneNumber(db, "(123)456-7892")
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func insertPhoneNumber(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
    CREATE TABLE IF NOT EXISTS phone_numbers (
      id SERIAL,
      value VARCHAR(255)
    )`
	_, err := db.Exec(statement)
	return err
}

func createDB(db *sql.DB, dbname string) error {
	_, err := db.Exec("CREATE DATABASE " + dbname)
	if err != nil {
		return err
	}
	return nil
}

func resetDB(db *sql.DB, dbname string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + dbname)
	if err != nil {
		return err
	}
	return createDB(db, dbname)
}

func normalize(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch <= '9' && ch >= '0' {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

func normalizeREGEXP(phone string) string {
	re := regexp.MustCompile("\\D") //find all non int
	return re.ReplaceAllString(phone, "")
}
