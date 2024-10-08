package cmd

import (
	"database/sql"
	s "github.com/VageLO/money-tree/structs"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func initDB() error {
	url := source.Config.Database

	// Check if database file exist, if exist return.
	fileInfo, _ := os.Stat(url)
	if fileInfo != nil {
		log.Println("Database file exist")
		return nil
	}

	_, err := os.Create(url)
	check(err)

	db, err := sql.Open("sqlite3", url)
	check(err)

	array := []string{
		s.InitTransactions,
		s.InitAccounts,
		s.InitCategories,
		s.UpdateOnDelete,
		s.UpdateOnInsert,
		s.UpdateOnUpdate,
		s.UpdateToAccount,
	}

	for _, query := range array {
		_, err = db.Exec(query)
		check(err)
	}

	defer db.Close()
	return nil
}
