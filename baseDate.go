package main

import (
	"database/sql"
	"log"

	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./bdwords.db")
	if err != nil {
		log.Fatal(err)
	}

	// Создание таблицы, если она еще не существует
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS words (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		word TEXT,
		translation TEXT
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func addDB(db *sql.DB, wordEntry *widget.Entry, translationEntry *widget.Entry) {
	word := wordEntry.Text
	translation := translationEntry.Text

	// Вставка данных в базу данных
	stmt, err := db.Prepare("INSERT INTO words(word, translation) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(word, translation)
	if err != nil {
		log.Fatal(err)
	}

	// Очистка полей ввода
	wordEntry.SetText("")
	translationEntry.SetText("")
}
