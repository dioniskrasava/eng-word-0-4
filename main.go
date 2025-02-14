package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := initDB()
	defer db.Close()

	a := app.New()
	w := a.NewWindow("Управление словарем")
	w.Resize(fyne.NewSize(500, 200))

	// Устанавливаем светлую тему
	a.Settings().SetTheme(theme.LightTheme())

	// Элементы интерфейса
	lab1 := widget.NewLabelWithStyle("Добавьте новое слово", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	wordEntry := widget.NewEntry()
	wordEntry.SetPlaceHolder("Английское слово")
	translationEntry := widget.NewEntry()
	translationEntry.SetPlaceHolder("Русский перевод")

	// Кнопки
	btn1 := widget.NewButton("Добавить", func() { addDB(db, wordEntry, translationEntry) })
	btn2 := widget.NewButton("Удалить", func() {})
	btn3 := widget.NewButton("Редактировать", func() {
		// Выполняем запрос для получения данных из таблицы
		rows, err := db.Query("SELECT id, word, translation FROM words")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Создаем слайс для хранения данных
		var data [][]string

		// Читаем данные из результата запроса
		for rows.Next() {
			var id int
			var word string
			var translation string
			err = rows.Scan(&id, &word, &translation)
			if err != nil {
				log.Fatal(err)
			}
			data = append(data, []string{fmt.Sprint(id), word, translation})
		}

		// Проверяем на ошибки после чтения
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}

		// Создаем новое окно
		newWindow := a.NewWindow("Окно редактирования")
		newWindow.Resize(fyne.NewSize(600, 400))

		// Если данных нет, выводим сообщение
		if len(data) == 0 {
			newWindow.SetContent(widget.NewLabel("Нет данных для отображения"))
			newWindow.Show()
			return
		}

		// Создаем таблицу для отображения данных
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

		// Кнопка для закрытия окна
		closeBtn := widget.NewButton("Закрыть", func() {
			newWindow.Close()
		})

		// Устанавливаем содержимое нового окна
		newWindow.SetContent(container.NewBorder(nil, closeBtn, nil, nil, container.NewScroll(table)))
		newWindow.Show()
	})

	// Компоновка интерфейса
	rows1 := container.NewGridWithColumns(1, wordEntry)
	rows2 := container.NewGridWithColumns(1, translationEntry)
	rows3 := container.NewGridWithColumns(3, btn1, btn2, btn3) // строка кнопок

	rowsCont := container.NewGridWithRows(5, lab1, rows1, rows2, rows3)
	w.SetContent(rowsCont)

	w.ShowAndRun()
}

// initDB инициализирует подключение к базе данных SQLite
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

// addDB добавляет слово и перевод в базу данных
func addDB(db *sql.DB, wordEntry *widget.Entry, translationEntry *widget.Entry) {
	word := strings.ToLower(wordEntry.Text)
	translation := strings.ToLower(translationEntry.Text)

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
