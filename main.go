package main

import (
	"github.com/rivo/tview"
)

const (
	APIEndPoint = "http://172.21.192.1:8080"
)

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
