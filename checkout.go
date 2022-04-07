package main

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func PurchaseDetailPage(id uint64) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n[yellow:]ESC[white:]:返回")

	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	purchase, err := GetPurchase(id)
	if err != nil {
		pages.SwitchToPage("menu")
		return flex
	}
	products := GetProducts()

	table.SetCell(0, 0, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 1, tview.NewTableCell("商品").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 2, tview.NewTableCell("數量").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	table.SetCell(0, 3, tview.NewTableCell("價格").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))

	var total int64 = 0
	for i, item := range purchase.PurchaseDetails {
		product := products.GetProductByID(item.ProductID)
		productName := ""
		if product != nil {
			productName = product.Name
		}
		table.SetCell(i+1, 0, tview.NewTableCell(strconv.FormatUint(item.ProductID, 10)).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 1, tview.NewTableCell(productName).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 2, tview.NewTableCell(strconv.FormatInt(item.Quantity, 10)).SetAlign(tview.AlignCenter))
		table.SetCell(i+1, 3, tview.NewTableCell(strconv.FormatInt(item.Total, 10)).SetAlign(tview.AlignCenter))
		total += item.Total
	}
	table.SetCell(len(products.Products)+1, 0, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	table.SetCell(len(products.Products)+1, 1, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	table.SetCell(len(products.Products)+1, 2, tview.NewTableCell("").SetAlign(tview.AlignCenter))
	table.SetCell(len(products.Products)+1, 3, tview.NewTableCell(strconv.FormatInt(total, 10)).SetAlign(tview.AlignCenter))

	table.SetFixed(1, 0)
	table.Select(1, 0)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.SwitchToPage("checkoutPage")
			return nil
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

func checkoutPage() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n[yellow:]Enter[white:]:選擇/取消選擇\t[yellow:]d[white:]:詳細資料\t[yellow:]ESC[white:]:返回")

	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	totalTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter)
	reminderTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n\n[red:]請確認您已繳納正確金額，再按下送出鍵結帳！[white:]")
	// reminderTextView.SetBorder(true)

	errorTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter)

	rightPanel.AddItem(reminderTextView, 0, 2, false)
	rightPanel.AddItem(totalTextView, 0, 1, false)
	rightPanel.AddItem(errorTextView, 0, 1, false)
	rightPanel.AddItem(nil, 0, 2, false)

	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	purchases, err := GetUnpaidPurchasesByUserID(userID)
	if err != nil {
		panic(err)
	}
	selections := []bool{}
	for i := 0; i < len(purchases); i++ {
		selections = append(selections, true)
	}

	drawTable := func() {
		table.Clear()
		table.SetCell(0, 0, tview.NewTableCell("選擇").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(0, 1, tview.NewTableCell("訂單日期").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(0, 2, tview.NewTableCell("金額").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))

		var total int64 = 0
		for i, purchase := range purchases {
			var _total int64 = 0
			for _, item := range purchase.PurchaseDetails {
				_total += item.Total
			}
			if selections[i] {
				table.SetCell(i+1, 0, tview.NewTableCell("✓").SetAlign(tview.AlignCenter))
				total += _total
			} else {
				table.SetCell(i+1, 0, tview.NewTableCell(""))
			}

			table.SetCell(i+1, 1, tview.NewTableCell(purchase.CreatedAt.String()))
			table.SetCell(i+1, 2, tview.NewTableCell(strconv.FormatInt(_total, 10)).SetAlign(tview.AlignRight))
		}
		table.SetCell(len(purchases)+1, 0, tview.NewTableCell(""))
		table.SetCell(len(purchases)+1, 1, tview.NewTableCell(""))
		table.SetCell(len(purchases)+1, 2, tview.NewTableCell(strconv.FormatInt(total, 10)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignRight))
		totalTextView.SetText("總金額: " + strconv.FormatInt(total, 10))
	}

	table.SetSelectedFunc(func(row, column int) {
		if row == 0 || row == len(purchases)+1 {
			return
		}
		selections[row-1] = !selections[row-1]
		drawTable()
	})

	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			pages.SwitchToPage("menu")
		}
	})
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRight {
			app.SetFocus(rightPanel)
			return nil
		}
		if event.Key() == tcell.KeyRune && event.Rune() == 'd' {
			i, _ := table.GetSelection()
			if i == 0 || i == len(purchases)+1 {
				return event
			}
			purchaseID := purchases[i-1].ID
			pages.AddAndSwitchToPage("purchaseDetailPage", PurchaseDetailPage(purchaseID), true)
			return nil
		}
		return event
	})
	rightPanel.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft {
			app.SetFocus(table)
			return nil
		}
		return event
	})

	drawTable()
	table.SetFixed(1, 0)
	table.Select(1, 0)

	submitButton := tview.NewButton("結帳").SetSelectedFunc(func() {
		pr := PayRequest{
			UserID:      userID,
			PurchaseIDs: []uint64{},
		}

		for i, purchase := range purchases {
			if selections[i] {
				pr.PurchaseIDs = append(pr.PurchaseIDs, purchase.ID)
			}
		}

		if len(pr.PurchaseIDs) == 0 {
			return
		}

		err := pay(pr)
		if err != nil {
			errorTextView.SetText(fmt.Sprintf("[:red]%s[:white]", err.Error()))
		}

		pages.SwitchToPage("menu")
	})
	rightPanel.AddItem(tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(submitButton, 0, 1, true).
		AddItem(nil, 0, 1, false), 0, 1, true)

	flex = flex.
		AddItem(header, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(table, 0, 2, true).
			AddItem(rightPanel, 0, 2, true).
			AddItem(nil, 0, 1, false),
			0, 4, true)
	return flex
}
