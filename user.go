package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func mePage() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	menu := tview.NewList().
		AddItem("別名管理", "", 'a', func() {
			pages.AddAndSwitchToPage("myAliasPage", myAliasPage(), true)
		})

	flex.AddItem(nil, 0, 1, false)
	flex.AddItem(
		tview.NewFlex().
			AddItem(nil, 0, 3, false).
			AddItem(menu, 0, 1, true).
			AddItem(nil, 0, 3, false),
		0, 3, true,
	)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.SwitchToPage("menu")
		}
		return event
	})

	return flex
}

func myAliasPage() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	refresh := false

	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n[yellow:]+[white:]:新增別名\t[yellow:]DELETE[white:]:刪除別名\t[yellow:]ESC[white:]:返回")

	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)
	drawTable := func() {
		table.Clear()
		table.SetCell(0, 0, tview.NewTableCell("ID").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(0, 1, tview.NewTableCell("別名").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		me, err := GetUserByID(userID)
		if err != nil {
			pages.AddAndSwitchToPage("loginPage", loginPage(), true)
			return
		}
		for i, username := range me.Usernames {
			table.SetCell(i+1, 0, tview.NewTableCell(strconv.FormatUint(username.ID, 10)).SetAlign(tview.AlignCenter))
			table.SetCell(i+1, 1, tview.NewTableCell(username.Name).SetAlign(tview.AlignCenter))
		}
		table.Select(1, 0)
	}
	drawTable()
	table.SetFixed(1, 0)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.SwitchToPage("mePage")
		}
		if refresh {
			refresh = false
			drawTable()
		}

		i, _ := table.GetSelection()
		if i == 0 {
			return event
		}
		id, err := strconv.ParseUint(table.GetCell(i, 0).Text, 10, 64)
		if err != nil {
			return nil
		}
		if event.Key() == tcell.KeyDelete {
			if err := DeleteUsername(id); err != nil {
				pages.AddAndSwitchToPage("loginPage", loginPage(), true)
				return nil
			} else {
				drawTable()
			}
		} else if event.Key() == tcell.KeyRune && event.Rune() == '+' {
			pages.AddAndSwitchToPage("addUsernamePage", addUsernamePage("myAliasPage", func() { drawTable() }), true)
			pages.ShowPage("myAliasPage")
			refresh = true
			return nil
		}
		return event
	})

	flex.AddItem(header, 0, 1, false)
	flex.AddItem(
		tview.NewFlex().
			AddItem(nil, 0, 3, false).
			AddItem(table, 0, 1, true).
			AddItem(nil, 0, 3, false),
		0, 3, true,
	)

	return flex
}

func addUsernamePage(backPage string, callback func()) tview.Primitive {
	var form *tview.Form
	form = tview.NewForm().
		AddInputField("別名", "", 20, nil, nil).
		AddButton("送出", func() {
			name := form.GetFormItem(0).(*tview.InputField).GetText()
			err := NewUsername(userID, name)
			if err != nil {
				return
			}
			pages.HidePage("addUsernamePage")
			pages.SwitchToPage(backPage)
			callback()
		}).
		AddButton("取消", func() {
			pages.HidePage("addUsernamePage")
			pages.SwitchToPage(backPage)
			callback()
		})
	form.SetBorder(true).SetTitle("新增別名")

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
