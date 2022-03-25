package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func importFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("files", fi.Name())
	if err != nil {
		return err
	}
	part.Write(fileContents)

	err = writer.Close()
	if err != nil {
		return err
	}

	res, err := http.NewRequest("POST", APIEndPoint+"/import/", body)
	if err != nil {
		return err
	}
	if res.Response.StatusCode != 200 {
		return error(fmt.Errorf("%d", res.Response.StatusCode))
	}
	return nil
}

func importPage() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	var form *tview.Form
	form = tview.NewForm().
		AddInputField("來源資料夾", "", 20, nil, nil).
		AddButton("Save", func() {
			rootDir := form.GetFormItem(0).(*tview.InputField).GetText()
			files, err := filepath.Glob(path.Join(rootDir, "*.json"))
			if err != nil {
				return
			}
			for _, file := range files {
				importFile(file)
			}
			pages.SwitchToPage("menu")
		})
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.SwitchToPage("menu")
			return nil
		}
		i, j := form.GetFocusedItemIndex()

		if event.Key() == tcell.KeyUp && i > 0 {
			app.SetFocus(form.GetFormItem(i - 1))
			return nil
		}
		if event.Key() == tcell.KeyDown {
			if i == form.GetFormItemCount()-1 {
				app.SetFocus(form.GetButton(0))
			} else {
				app.SetFocus(form.GetFormItem(i + 1))
			}
			return nil
		}
		if event.Key() == tcell.KeyUp && i == -1 {
			app.SetFocus(form.GetFormItem(0))
			return nil
		}
		if event.Key() == tcell.KeyRight && j < form.GetButtonCount()-1 {
			app.SetFocus(form.GetButton(j + 1))
			return nil
		}
		if event.Key() == tcell.KeyLeft && j > 0 {
			app.SetFocus(form.GetButton(j - 1))
			return nil
		}
		return event
	})
	flex = flex.
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(form, 0, 1, true).
			AddItem(nil, 0, 1, false),
			0, 1, true).
		AddItem(nil, 0, 1, false)
	return flex
}
