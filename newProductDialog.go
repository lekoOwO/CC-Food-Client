package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewProductDialogPage(barcode string, backPage string) tview.Primitive {
	var form *tview.Form
	form = tview.NewForm().
		AddInputField("品名", "", 20, nil, nil).
		AddInputField("單價", "", 20, nil, nil).
		AddInputField("條碼", barcode, 20, nil, nil).
		AddButton("送出", func() {
			name := form.GetFormItem(0).(*tview.InputField).GetText()
			price, _ := strconv.ParseInt(form.GetFormItem(1).(*tview.InputField).GetText(), 10, 64)
			barcode := form.GetFormItem(2).(*tview.InputField).GetText()
			req := NewProductRequest{
				Name:    name,
				Price:   price,
				Barcode: barcode,
			}
			_, err := CreateProduct(req)
			if err != nil {
				return
			}
			pages.HidePage("newProductDialogPage")
			pages.SwitchToPage(backPage)
		}).
		AddButton("取消", func() {
			pages.HidePage("newProductDialogPage")
			pages.SwitchToPage(backPage)
		})
	form.SetBorder(true).SetTitle("新增商品")

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

	return modal(form, 30, 11)
}
