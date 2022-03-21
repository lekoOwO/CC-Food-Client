package main

import (
	"github.com/rivo/tview"
)

func menu() *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n[:blue]CC-Food")

	menu := tview.NewList().
		AddItem("購買", "", 'a', func() {
			pages.AddAndSwitchToPage("loginPage", loginPage(), true)
		}).
		AddItem("結清帳款", "", 'b', nil)

	flex.AddItem(header, 0, 1, false)
	flex.AddItem(
		tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(menu, 0, 1, true).
			AddItem(nil, 0, 1, false),
		0, 1, true,
	)

	return flex
}
