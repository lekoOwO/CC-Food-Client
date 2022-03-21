package main

import (
	"github.com/rivo/tview"
)

const (
	APIEndPoint = "http://172.21.192.1:8080"
)

var pages *tview.Pages
var userID uint = 0
var lastPage string = "menu"

var pageFuncs = map[string]func() *tview.Flex{}

func init() {
	pageFuncs["menu"] = menu
	pageFuncs["loginPage"] = loginPage
	pageFuncs["registerAndLoginPage"] = registerAndLoginPage
	pageFuncs["buyPage"] = buyPage
}

func main() {
	app := tview.NewApplication()
	pages = tview.NewPages()

	pages.AddPage("menu", menu(), true, true)
	pages.AddPage("loginPage", loginPage(), true, false)
	pages.AddPage("registerAndLoginPage", registerAndLoginPage(), true, false)
	pages.AddPage("buyPage", buyPage(), true, false)
	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}
