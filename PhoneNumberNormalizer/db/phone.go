package db

import (
	"database/sql"
)

type DB struct {
	db *sql.DB
}

//Phone represents the phone_numbers table
type Phone struct {
	ID     int
	Number string
}

//Open opens a connection to a db
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

//Close closes a connection to a db
func (db *DB) Close() error {
	return db.db.Close()
}

//Migrate creates table
func Migrate(driverName, dataSourceName string) error {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	err = createPhoneNumbersTable(db)
	if err != nil {
		return err
	}
	return db.Close()
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

//Reset creates DB - deletes if exists
func Reset(driverName, dataSourceName, dbName string) error {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	//drop the connection
	return db.Close()
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

func (db *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, number := range data {
		if _, err := insertPhoneNumber(db.db, number); err != nil {
			return err
		}
	}
	return nil
}
func insertPhoneNumber(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	// result, err := db.Exec(statement, phone)
	// id, err := result.LastInsertId()

	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

//GetAllPhones iterates over entries in the db, get all records
func (db *DB) GetAllPhones() ([]Phone, error) {
	rows, err := db.db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return ret, nil
}

//UpdatePhone updates given Phone instance
func (db *DB) UpdatePhone(p *Phone) error {
	statement := "UPDATE phone_numbers SET value=$2 WHERE id=$1 "
	_, err := db.db.Exec(statement, p.ID, p.Number)
	return err
}

//DeletePhone deletes record by id
func (db *DB) DeletePhone(id int) error {
	statement := "DELETE FROM phone_numbers WHERE id=$1 "
	_, err := db.db.Exec(statement, id)
	return err
}

//FindPhone gets phone by value of number
func (db *DB) FindPhone(number string) (*Phone, error) {
	var p Phone
	row := db.db.QueryRow("SELECT * FROM phone_numbers WHERE value = $1", number)
	err := row.Scan(&p.ID, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}
