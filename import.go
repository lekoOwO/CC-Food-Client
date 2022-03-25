package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createForm(form map[string]string) (string, io.Reader, error) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	defer mp.Close()
	for key, val := range form {
		if strings.HasPrefix(val, "@") {
			val = val[1:]
			file, err := os.Open(val)
			if err != nil {
				return "", nil, err
			}
			defer file.Close()
			part, err := mp.CreateFormFile(key, val)
			if err != nil {
				return "", nil, err
			}
			io.Copy(part, file)
		} else {
			mp.WriteField(key, val)
		}
	}
	return mp.FormDataContentType(), body, nil
}

func importFile(path string) error {
	form := map[string]string{"files": "@" + path}
	ct, body, err := createForm(form)
	if err != nil {
		return err
	}
	resp, err := http.Post(APIEndPoint+"/import", ct, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("%d", resp.StatusCode)
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
