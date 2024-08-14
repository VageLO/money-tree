package cmd

import (
	"database/sql"
	"log"
	m "github.com/VageLO/money-tree/modal"
	s "github.com/VageLO/money-tree/structs"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() error {
	defer m.ErrorModal(source.Pages, source.Modal)
	url := "./database.db"

	// Check if database file exist, if exist return.
	fileInfo, _ := os.Stat(url)
	if fileInfo != nil {
		log.Println("Database file exist")
		return nil
	}

	os.Create(url)

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
