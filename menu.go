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
		AddItem("購買", "", 'a', func() {
			pages.AddAndSwitchToPage("buyPage", buyPage(), true)
		}).
		AddItem("結清帳款", "", 'b', func() {
			pages.AddAndSwitchToPage("checkoutPage", checkoutPage(), true)
		}).
		AddItem("商品管理", "", 'c', func() {
			pages.AddAndSwitchToPage("productManagePage", ProductManagePage(), true)
		}).
		AddItem("登出", "", 'd', func() {
			pages.AddAndSwitchToPage("loginPage", loginPage(), true)
		})

	flex.AddItem(header, 0, 2, false)
	flex.AddItem(
		tview.NewFlex().
			AddItem(nil, 0, 5, false).
			AddItem(menu, 0, 1, true).
			AddItem(nil, 0, 5, false),
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
