package cmd

import (
	"database/sql"
	"log"
	m "main/modal"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() error {
	defer m.ErrorModal(source.Pages, source.Modal)
	url := "./database.db"
	dir := "./sql/init"

	// Check if database file exist, if exist return.
	fileInfo, _ := os.Stat(url)
	if fileInfo != nil {
		log.Println("Database file exist")
		return nil
	}

	os.Create(url)

	db, err := sql.Open("sqlite3", url)
	check(err)

	files, err := os.ReadDir(dir)
	check(err)

	for _, file := range files {
		// TODO: change to filepath
		query, err := os.ReadFile(dir + "/" + file.Name())
		check(err)
		_, err = db.Exec(string(query))
		check(err)
	}

	defer db.Close()
	return nil
}
