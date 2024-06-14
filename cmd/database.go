package cmd

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() error {
	defer ErrorModal()
	url := "./database.db"

	// Check if database file exist, if exist return.
	fileInfo, _ := os.Stat(url)
	if fileInfo != nil {
		log.Println("Database file exist")
		return nil
	}

	os.Create(url)

	db, err := sql.Open("sqlite3", url)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	_, err = db.Exec(`CREATE TABLE "Transactions" (
	"id"	INTEGER,
	"account_id"	INTEGER NOT NULL,
	"category_id"	INTEGER NOT NULL,
	"transaction_type"	TEXT NOT NULL,
	"date"	TEXT NOT NULL,
	"amount"	NUMERIC NOT NULL,
	"balance"	NUMERIC NOT NULL,
	FOREIGN KEY("category_id") REFERENCES "Categories"("id") ON DELETE CASCADE,
	FOREIGN KEY("account_id") REFERENCES "Accounts"("id") ON DELETE CASCADE,
	PRIMARY KEY("id" AUTOINCREMENT)
)`)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	_, err = db.Exec(`CREATE TABLE "Accounts" (
  "id"	INTEGER,
  "title"	TEXT NOT NULL UNIQUE,
  "currency"	TEXT NOT NULL,
  "balance"	NUMERIC NOT NULL DEFAULT 0,
  PRIMARY KEY("id" AUTOINCREMENT)
)`)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	_, err = db.Exec(`CREATE TABLE "Categories" (
	"id"	INTEGER,
	"parent_id"	INTEGER,
	"title"	TEXT NOT NULL UNIQUE,
	PRIMARY KEY("id" AUTOINCREMENT)
)`)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer db.Close()
	return nil
}
