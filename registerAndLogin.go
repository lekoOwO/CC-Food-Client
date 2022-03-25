package main

import (
	"github.com/rivo/tview"
)

func registerAndLoginPage() tview.Primitive {
	userID = 0

	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n[:blue]CC-Food")

	errMessage := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter)

	var form *tview.Form
	form = tview.NewForm().
		AddInputField("顯示名稱", "", 20, nil, nil).
		AddInputField("帳號", "", 20, func(data string, _ rune) bool {
			_, err := GetUserByUsername(data)
			result := err != nil
			if !result {
				errMessage.SetText("[:red]帳號已被使用")
			} else {
				errMessage.SetText("")
			}
			return result
		}, nil).
		AddInputField("卡號", "", 20, func(data string, _ rune) bool {
			if data == "" {
				return true
			}
			_, err := GetUserByUsername(data)
			result := err != nil
			if !result {
				errMessage.SetText("[:red]卡號已被使用")
			} else {
				errMessage.SetText("")
			}
			return err != nil
		}, nil).
		AddButton("Save", func() {
			displayName := form.GetFormItem(0).(*tview.InputField).GetText()
			var usernames []string
			for i := 1; i <= 2; i++ {
				data := form.GetFormItem(i).(*tview.InputField).GetText()
				if data != "" {
					usernames = append(usernames, data)
				}
			}
			if len(usernames) == 0 {
				return
			}
			id, err := Register(displayName, usernames)
			if err != nil {
				errMessage.SetText("[:red]" + err.Error())
				return
			}
			userID = id
			pages.SwitchToPage("menu")
		}).
		AddButton("Quit", func() {
			pages.SwitchToPage("menu")
		})
	form.SetBorder(true).SetTitle("註冊帳號").SetTitleAlign(tview.AlignLeft)

	flex = flex.
		AddItem(header, 0, 1, false).
		AddItem(errMessage, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(form, 0, 1, true).
			AddItem(nil, 0, 1, false),
			0, 4, true).
		AddItem(nil, 0, 4, false)
	return flex
}
