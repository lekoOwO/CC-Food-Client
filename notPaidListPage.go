package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NotPaidListPage() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n[yellow:]ESC[white:]:返回")

	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	drawTable := func() {
		table.Clear()
		table.SetCell(0, 0, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(0, 1, tview.NewTableCell("用戶名").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(0, 2, tview.NewTableCell("總金額").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		purchases, err := GetUnpaidPurchases()
		if err != nil {
			panic(err)
		}

		userTotal := make(map[uint64]int64)
		for _, p := range purchases {
			if _, ok := userTotal[p.UserID]; !ok {
				userTotal[p.UserID] = 0
			}
			for _, pd := range p.PurchaseDetails {
				userTotal[p.UserID] += pd.Total
			}
		}

		usernames := make(map[uint64]string)
		for i, _ := range userTotal {
			user, err := GetUserByID(uint(i))
			if err != nil {
				panic(err)
			}
			usernames[i] = user.Usernames[0].Name
		}

		c := 0
		for i, t := range userTotal {
			table.SetCell(c+1, 0, tview.NewTableCell(strconv.FormatUint(i, 10)).SetAlign(tview.AlignCenter))
			table.SetCell(c+1, 1, tview.NewTableCell(usernames[i]).SetAlign(tview.AlignCenter))
			table.SetCell(c+1, 2, tview.NewTableCell(strconv.FormatInt(t, 10)).SetAlign(tview.AlignCenter))
		}
		table.Select(1, 0)
	}
	drawTable()
	table.SetFixed(1, 0)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.SwitchToPage("menu")
		}
		return event
	})

	flex = flex.
		AddItem(header, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 5, false).
			AddItem(table, 0, 2, true).
			AddItem(nil, 0, 5, false),
			0, 4, true)
	return flex
}
