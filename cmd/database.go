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

	transactions, err := os.ReadFile("./sql/Transactions.sql")
	check(err)

	accounts, err := os.ReadFile("./sql/Accounts.sql")
	check(err)

	categories, err := os.ReadFile("./sql/Categories.sql")
	check(err)

	trigger_insert, err := os.ReadFile("./sql/Update_Balance_On_Transaction_Insert.sql")
	check(err)

	trigger_update, err := os.ReadFile("./sql/Update_Balance_On_Transaction_Update.sql")
	check(err)

	trigger_delete, err := os.ReadFile("./sql/Update_Balance_On_Transaction_Delete.sql")
	check(err)

	// Check if database file exist, if exist return.
	fileInfo, _ := os.Stat(url)
	if fileInfo != nil {
		log.Println("Database file exist")
		return nil
	}

	os.Create(url)

	db, err := sql.Open("sqlite3", url)
	check(err)

	_, err = db.Exec(string(transactions))
	check(err)

	_, err = db.Exec(string(accounts))
	check(err)

	_, err = db.Exec(string(categories))
	check(err)

	_, err = db.Exec(string(trigger_insert))
	check(err)

	_, err = db.Exec(string(trigger_update))
	check(err)

	_, err = db.Exec(string(trigger_delete))
	check(err)

	defer db.Close()
	return nil
}
