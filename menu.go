package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func menu() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n[:blue]CC-Food")

	menu := tview.NewList().
		AddItem("購買", "", 'b', func() {
			pages.AddAndSwitchToPage("buyPage", buyPage(), true)
		}).
		AddItem("結清帳款", "", 'p', func() {
			pages.AddAndSwitchToPage("checkoutPage", checkoutPage(), true)
		}).
		AddItem("商品管理", "", 'P', func() {
			pages.AddAndSwitchToPage("productManagePage", ProductManagePage(), true)
		}).
		AddItem("我的資料", "", 'm', func() {
			pages.AddAndSwitchToPage("mePage", mePage(), true)
		}).
		AddItem("匯入舊系統資料", "", 'i', func() {
			pages.AddAndSwitchToPage("importPage", importPage(), true)
		}).
		AddItem("登出", "", 'o', func() {
			pages.AddAndSwitchToPage("loginPage", loginPage(), true)
		})

	flex.AddItem(header, 0, 1, false)
	flex.AddItem(
		tview.NewFlex().
			AddItem(nil, 0, 3, false).
			AddItem(menu, 0, 1, true).
			AddItem(nil, 0, 3, false),
		0, 3, true,
	)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.AddAndSwitchToPage("loginPage", loginPage(), true)
		}
		return event
	})

	return flex
}
