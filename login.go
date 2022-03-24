package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func loginPage() tview.Primitive {
	userID = 0

	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n[:blue]CC-Food")

	var loginBox *tview.InputField
	loginBox = tview.NewInputField().
		SetLabel("帳號: ").
		SetLabelColor(tcell.ColorWhite).
		SetFieldStyle(tcell.StyleDefault.Background(tcell.ColorBlue)).
		SetFieldTextColor(tcell.ColorWhite).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				username := loginBox.GetText()
				userMinimal, err := GetUserByUsername(username)
				if err != nil {
					// no user found
					if err.Error() == "404" {
						pages.AddAndSwitchToPage("registerAndLoginPage", registerAndLoginPage(), true)
						return
					} else {
						pages.AddAndSwitchToPage("loginPage", loginPage(), true)
						return
					}
				}
				userID = userMinimal.ID
				pages.SwitchToPage("menu")
			}
		})

	flex = flex.
		AddItem(header, 0, 3, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 2, false).
			AddItem(loginBox, 0, 1, true).
			AddItem(nil, 0, 2, false),
			0, 4, true)
	return flex
}
