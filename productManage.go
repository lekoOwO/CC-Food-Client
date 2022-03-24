package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func EditProductPage(id uint64) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	product, err := GetProductByID(id)
	if err != nil {
		panic(err)
	}

	var form *tview.Form
	form = tview.NewForm().
		AddInputField("品名", product.Name, 20, nil, nil).
		AddInputField("單價", strconv.FormatInt(int64(product.Price), 10), 20, nil, nil).
		AddInputField("條碼", product.Barcode, 20, nil, nil).
		AddButton("送出", func() {
			name := form.GetFormItem(0).(*tview.InputField).GetText()
			price, err := strconv.ParseInt(form.GetFormItem(1).(*tview.InputField).GetText(), 10, 64)
			if err != nil {
				return
			}
			barcode := form.GetFormItem(2).(*tview.InputField).GetText()
			req := NewProductRequest{
				Name:    name,
				Price:   price,
				Barcode: barcode,
			}
			err = EditProduct(id, req)
			if err != nil {
				return
			}
			pages.AddAndSwitchToPage("productManagePage", ProductManagePage(), true)
		})
	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.SwitchToPage("productManagePage")
		}
		return event
	})
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
	form.SetBorder(true).SetTitle("編輯商品")
	flex = flex.
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(form, 0, 1, true).
			AddItem(nil, 0, 1, false),
			0, 3, true).
		AddItem(nil, 0, 1, false)
	return flex
}

func DeleteProductConfirmDialogPage(id uint64) tview.Primitive {
	m := tview.NewModal().
		SetText("你確定要刪除這個商品嗎?").
		AddButtons([]string{"是", "否"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "是" {
				DelectProduct(id)
			}
			pages.AddAndSwitchToPage("productManagePage", ProductManagePage(), true)
		})
	return modal(m, 20, 10)
}

func ProductManagePage() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	refresh := false

	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n[yellow:]+[white:]:新增商品\t[yellow:]DELETE[white:]:刪除商品\t[yellow:]ENTER[white:]:編輯商品\t[yellow:]ESC[white:]:返回")

	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	products := GetProducts()
	drawTable := func() {
		table.Clear()
		table.SetCell(0, 0, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(0, 1, tview.NewTableCell("商品名稱").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(0, 2, tview.NewTableCell("單價").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(0, 3, tview.NewTableCell("條碼").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		products = GetProducts()
		for i, product := range products.Products {
			table.SetCell(i+1, 0, tview.NewTableCell(strconv.FormatUint(product.ID, 10)).SetAlign(tview.AlignCenter))
			table.SetCell(i+1, 1, tview.NewTableCell(product.Name).SetAlign(tview.AlignCenter))
			table.SetCell(i+1, 2, tview.NewTableCell(strconv.FormatUint(product.Price, 10)).SetAlign(tview.AlignCenter))
			table.SetCell(i+1, 3, tview.NewTableCell(product.Barcode).SetAlign(tview.AlignCenter))
		}
		table.Select(1, 0)
	}
	drawTable()
	table.SetFixed(1, 0)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if refresh {
			refresh = false
			drawTable()
		}
		i, _ := table.GetSelection()
		if i == 0 {
			return event
		}
		id := products.Products[i-1].ID
		if event.Key() == tcell.KeyEnter {
			refresh = true
			pages.AddAndSwitchToPage("editProductPage", EditProductPage(id), true)
		} else if event.Key() == tcell.KeyDelete {
			pages.AddAndSwitchToPage("DeleteProductConfirmDialogPage", DeleteProductConfirmDialogPage(id), true)
			pages.ShowPage("productManagePage")
		} else if event.Key() == tcell.KeyEscape {
			pages.SwitchToPage("menu")
		} else if event.Key() == tcell.KeyRune && event.Rune() == '+' {
			pages.AddAndSwitchToPage("newProductDialogPage", NewProductDialogPage("", "productManagePage"), true)
			refresh = true
			pages.ShowPage("productManagePage")
		}
		return event
	})

	flex = flex.
		AddItem(header, 0, 2, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 2, false).
			AddItem(table, 0, 1, true).
			AddItem(nil, 0, 2, false),
			0, 4, true)
	return flex
}
