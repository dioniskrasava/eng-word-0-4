package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme" // Импортируем пакет для работы с темами
	"fyne.io/fyne/v2/widget"
)

func main() {

	db := initDB()
	defer db.Close()

	a := app.New()
	w := a.NewWindow("Hello")
	w.Resize(fyne.NewSize(500, 150))

	// Устанавливаем светлую тему
	a.Settings().SetTheme(theme.LightTheme())

	//-----------------------------------------------------------------------------------------------
	lab1 := widget.NewLabelWithStyle("Напишите слово", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	//-----------------------------------------------------------------------------------------------
	wordEntry := widget.NewEntry()
	wordEntry.SetPlaceHolder("Англиское слово")
	translationEntry := widget.NewEntry()
	translationEntry.SetPlaceHolder("Русский перевод")
	//-----------
	btn1 := widget.NewButton("Добавить", func() { addDB(db, wordEntry, translationEntry) })
	btn2 := widget.NewButton("Удалить", func() {})
	btn3 := widget.NewButton("Редактировать", func() {})

	//-----------------------------------------------------------------------------------------------

	// СТОЛБЕЦ 1
	rows1 := container.NewGridWithColumns(1, wordEntry)
	rows2 := container.NewGridWithColumns(1, translationEntry)
	rows3 := container.NewGridWithColumns(3, btn1, btn2, btn3) // строка кнопок

	// СТРОКИ
	rowsCont := container.NewGridWithRows(5, lab1, rows1, rows2, rows3)

	w.SetContent(rowsCont)

	w.ShowAndRun()
}
