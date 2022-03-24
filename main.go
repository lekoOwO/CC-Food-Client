package main

import (
	"os"

	"github.com/rivo/tview"
)

var APIEndPoint = os.Getenv("CC_FOOD_API")

var app *tview.Application
var pages *tview.Pages
var userID uint = 0

func main() {
	app = tview.NewApplication()
	pages = tview.NewPages()

	pages.AddPage("menu", menu(), true, false)
	pages.AddPage("loginPage", loginPage(), true, true)
	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}
