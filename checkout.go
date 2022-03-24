package main

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func checkoutPage() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	rightPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	totalTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter)
	reminderTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n[:blue]請確認您已繳納正確金額，再按下送出鍵結帳！[:white]")
	// reminderTextView.SetBorder(true)

	errorTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter)

	rightPanel.AddItem(reminderTextView, 0, 3, false)
	rightPanel.AddItem(totalTextView, 0, 1, false)
	rightPanel.AddItem(errorTextView, 0, 1, false)

	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	purchases, err := GetUnpaidPurchases(userID)
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

		var total uint64 = 0
		for i, purchase := range purchases {
			var _total uint64 = 0
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
			table.SetCell(i+1, 2, tview.NewTableCell(strconv.FormatUint(_total, 10)).SetAlign(tview.AlignRight))
		}
		table.SetCell(len(purchases)+1, 0, tview.NewTableCell(""))
		table.SetCell(len(purchases)+1, 1, tview.NewTableCell(""))
		table.SetCell(len(purchases)+1, 2, tview.NewTableCell(strconv.FormatUint(total, 10)).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignRight))
		totalTextView.SetText("總金額: " + strconv.FormatUint(total, 10))
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
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(table, 0, 2, true).
			AddItem(rightPanel, 0, 2, true).
			AddItem(nil, 0, 1, false),
			0, 4, true)
	return flex
}
