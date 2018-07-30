package main

import (
	"bytes"
	"fmt"
	"regexp"

	phonedb "github.com/gophercises/PhoneNumberNormalizer/db"
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
	//create DB
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	handleErr(phonedb.Reset("postgres", psqlInfo, dbname))

	//disconnect DB, create tables
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	handleErr(phonedb.Migrate("postgres", psqlInfo))

	//open db
	db, err := phonedb.Open("postgres", psqlInfo)
	handleErr(err)
	defer db.Close()

	//insert values into table
	err = db.Seed()
	handleErr(err)

	//retrive all records
	phones, err := db.GetAllPhones()
	handleErr(err)
	for _, p := range phones {
		number := normalize(p.Number)
		if p.Number != number {
			existing, err := db.FindPhone(number)
			handleErr(err)
			if existing != nil {
				//does exist, delete
				fmt.Println("Deleting ... ", p.Number)
				handleErr(db.DeletePhone(p.ID))
			} else {
				//does not exist, update
				fmt.Println("Updating ... ", p.Number)
				p.Number = number
				handleErr(db.UpdatePhone(&p))
			}
		} else {
			fmt.Println("No changes required")
		}
	}
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
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
	re := regexp.MustCompile("\\D") //only int left
	return re.ReplaceAllString(phone, "")
}
