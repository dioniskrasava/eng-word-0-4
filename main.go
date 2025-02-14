package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

const (
	windowWidth  = 300
	windowHeight = 200
	tableWidth   = 600
	tableHeight  = 400
)

func main() {
	db := initDB()
	defer db.Close()

	app := app.New()
	window := createMainWindow(app, db)

	window.Resize(fyne.NewSize(windowWidth, windowHeight))
	window.ShowAndRun()
}

// createMainWindow создает главное окно приложения
func createMainWindow(app fyne.App, db *sql.DB) fyne.Window {
	window := app.NewWindow("Английские слова")

	// Элементы интерфейса
	titleLabel := widget.NewLabelWithStyle("Добавьте новое слово", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	wordEntry := widget.NewEntry()
	wordEntry.SetPlaceHolder("Английское слово")
	translationEntry := widget.NewEntry()
	translationEntry.SetPlaceHolder("Русский перевод")

	// Кнопки
	addButton := widget.NewButton("Добавить", func() { addWordToDB(db, wordEntry, translationEntry) })
	viewButton := widget.NewButton("Посмотреть", func() { showWordsWindow(app, db) })

	// Компоновка интерфейса
	wordRow := container.NewGridWithColumns(1, wordEntry)
	translationRow := container.NewGridWithColumns(1, translationEntry)
	buttonRow := container.NewGridWithColumns(2, addButton, viewButton)

	content := container.NewGridWithRows(4, titleLabel, wordRow, translationRow, buttonRow)
	window.SetContent(content)

	return window
}

// initDB инициализирует подключение к базе данных SQLite
func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./bdwords.db")
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
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
		log.Fatal("Ошибка создания таблицы:", err)
	}

	return db
}

// addWordToDB добавляет слово и перевод в базу данных
func addWordToDB(db *sql.DB, wordEntry *widget.Entry, translationEntry *widget.Entry) {
	word := strings.ToLower(wordEntry.Text)
	translation := strings.ToLower(translationEntry.Text)

	// Вставка данных в базу данных
	stmt, err := db.Prepare("INSERT INTO words(word, translation) values(?, ?)")
	if err != nil {
		log.Fatal("Ошибка подготовки запроса:", err)
	}
	_, err = stmt.Exec(word, translation)
	if err != nil {
		log.Fatal("Ошибка выполнения запроса:", err)
	}

	// Очистка полей ввода
	wordEntry.SetText("")
	translationEntry.SetText("")
}

// showWordsWindow создает окно для просмотра слов
func showWordsWindow(app fyne.App, db *sql.DB) {
	data := fetchWordsFromDB(db)

	newWindow := app.NewWindow("Просмотр слов")
	newWindow.Resize(fyne.NewSize(tableWidth, tableHeight))

	if len(data) == 0 {
		newWindow.SetContent(widget.NewLabel("Нет данных для отображения"))
		newWindow.Show()
		return
	}

	table := createTable(data)
	closeButton := widget.NewButton("Закрыть", func() {
		newWindow.Close()
	})

	newWindow.SetContent(container.NewBorder(nil, closeButton, nil, nil, container.NewScroll(table)))
	newWindow.Show()
}

// fetchWordsFromDB извлекает слова из базы данных
func fetchWordsFromDB(db *sql.DB) [][]string {
	rows, err := db.Query("SELECT id, word, translation FROM words")
	if err != nil {
		log.Fatal("Ошибка выполнения запроса:", err)
	}
	defer rows.Close()

	var data [][]string
	for rows.Next() {
		var id int
		var word string
		var translation string
		err = rows.Scan(&id, &word, &translation)
		if err != nil {
			log.Fatal("Ошибка чтения данных:", err)
		}
		data = append(data, []string{fmt.Sprint(id), word, translation})
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Ошибка после чтения данных:", err)
	}

	return data
}

// createTable создает таблицу для отображения данных
func createTable(data [][]string) *widget.Table {
	calculateMaxColumnWidth := func(columnIndex int) float32 {
		maxWidth := float32(0)
		for _, row := range data {
			text := row[columnIndex]
			label := widget.NewLabel(text)
			label.Refresh()
			width := label.MinSize().Width
			if width > maxWidth {
				maxWidth = width
			}
		}
		return maxWidth
	}

	maxWordWidth := calculateMaxColumnWidth(1)

	table := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
		},
	)

	table.SetColumnWidth(1, maxWordWidth+20) // Добавляем отступ

	return table
}
