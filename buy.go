package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var RefreshProducts bool = false

func buySuccessDialogPage(total int) tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	text := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n購買成功!\n金額: " + strconv.Itoa(total) + " 元")

	button := tview.NewButton("登出").SetSelectedFunc(func() {
		pages.AddAndSwitchToPage("loginPage", loginPage(), true)
	})

	flex = flex.
		AddItem(nil, 0, 2, false).
		AddItem(text, 0, 3, false).
		AddItem(button, 0, 2, true)

	flex.SetBorder(true).SetTitle("購買成功").SetTitleAlign(tview.AlignCenter)
	return modal(flex, 30, 11)
}

func buyPage() tview.Primitive {
	products := GetProducts()
	cart := []CartItem{}

	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	header := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("\n[yellow:]ESC[white:]:返回主選單\n🛈 使用[yellow:]方向鍵[white:]操控游標，[yellow:]Tab[white:] 操作下拉選單，在數量欄位按 [yellow:]Enter[white:] 將商品加入購物車")

	errText := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetText("")

	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false)

	initTable := func() {
		table.Clear()
		table.SetCell(0, 0, tview.NewTableCell("商品名稱").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(0, 1, tview.NewTableCell("數量").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(0, 2, tview.NewTableCell("金額").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	}
	initTable()

	var form *tview.Form
	total := 0

	drawTable := func() {
		initTable()
		total = 0
		for i, item := range cart {
			product := products.GetProductByID(item.ProductID)
			price := int(item.Quantity) * int(product.Price)
			table.SetCell(i+1, 0, tview.NewTableCell(product.Name).SetAlign(tview.AlignCenter))
			table.SetCell(i+1, 1, tview.NewTableCell(strconv.Itoa(int(item.Quantity))).SetAlign(tview.AlignCenter))
			table.SetCell(i+1, 2, tview.NewTableCell(strconv.Itoa(price)).SetAlign(tview.AlignCenter))
			total += price
		}
		table.SetCell(len(cart)+1, 0, tview.NewTableCell("總金額").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(len(cart)+1, 1, tview.NewTableCell("").SetAlign(tview.AlignCenter))
		table.SetCell(len(cart)+1, 2, tview.NewTableCell(strconv.Itoa(total)).SetAlign(tview.AlignCenter))
	}

	showAddProductDialog := func(barcode string) {
		app.SetFocus(form.GetFormItem(0))
		pages.AddAndSwitchToPage("newProductDialogPage", NewProductDialogPage(barcode, "buyPage", nil), true)
		RefreshProducts = true
		pages.ShowPage("buyPage")
	}

	add := func() {
		i, _ := form.GetFormItemByLabel("商品選擇").(*tview.DropDown).GetCurrentOption()
		quantityText := form.GetFormItemByLabel("數量").(*tview.InputField).GetText()
		if quantityText == "" {
			quantityText = "1"
		}
		quantity, err := strconv.Atoi(quantityText)
		if err != nil {
			return
		}
		cart = append(cart, CartItem{
			ProductID: products.Products[i].ID,
			Quantity:  int64(quantity),
		})
		form.GetFormItemByLabel("商品條碼").(*tview.InputField).SetText("")
		form.GetFormItemByLabel("商品選擇").(*tview.DropDown).SetCurrentOption(0)
		form.GetFormItemByLabel("數量").(*tview.InputField).SetText("")
		drawTable()
		app.SetFocus(form.GetFormItem(form.GetFormItemIndex("商品條碼")))
	}
	form = tview.NewForm().
		AddInputField("商品條碼", "", 20, nil, nil).
		AddDropDown("商品選擇", products.GetProductNames(), 0, nil).
		AddInputField("數量", "", 20, nil, nil).
		AddButton("結帳", func() {
			if total == 0 {
				errText.SetText("[red:]請先將商品加入購物車")
				app.SetFocus(form.GetFormItemByLabel("數量"))
				return
			}
			var brd []BuyRequestDetail
			for _, item := range cart {
				brd = append(brd, BuyRequestDetail{
					ProductID: item.ProductID,
					Quantity:  item.Quantity,
				})
			}
			br := BuyRequest{
				UserID:  userID,
				Details: brd,
			}
			err := Buy(br)
			if err != nil {
				errText.SetText("[:red]購買失敗：" + err.Error())
			} else {
				pages.AddAndSwitchToPage("buySuccessPage", buySuccessDialogPage(total), true)
				pages.ShowPage("buyPage")
			}
		}).
		AddButton("新增商品", func() {
			showAddProductDialog("")
		})

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.SwitchToPage("menu")
			return nil
		}
		if RefreshProducts {
			products = GetProducts()
			form.GetFormItemByLabel("商品選擇").(*tview.DropDown).SetOptions(products.GetProductNames(), nil)
			RefreshProducts = false
		}

		i, j := form.GetFocusedItemIndex()

		if event.Key() == tcell.KeyUp && i != -1 {
			if i == 0 {
				app.SetFocus(form.GetButton(0))
				return nil
			}
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
		if event.Key() == tcell.KeyEnter && i == form.GetFormItemIndex("商品條碼") {
			barcode := form.GetFormItemByLabel("商品條碼").(*tview.InputField).GetText()
			if barcode == "" {
				return event
			}

			i := products.GetIndexByBarcode(barcode)
			if i != -1 {
				form.GetFormItemByLabel("商品選擇").(*tview.DropDown).SetCurrentOption(i)
				add()
				return nil
			} else {
				showAddProductDialog(barcode)
				return nil
			}
		}
		if event.Key() == tcell.KeyEnter && i == form.GetFormItemIndex("數量") {
			add()
			return nil
		}
		if event.Key() == tcell.KeyRight && j == form.GetButtonIndex("送出") {
			app.SetFocus(table)
			return nil
		}
		return event
	})
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft {
			app.SetFocus(form.GetFormItem(0))
			return nil
		}
		if event.Key() == tcell.KeyDelete || event.Key() == tcell.KeyBackspace {
			i, _ := table.GetSelection()
			i -= 1
			if i < len(cart) && i >= 0 {
				cart = append(cart[:i], cart[i+1:]...)
				drawTable()
			}
			return nil
		}
		return event
	})
	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			pages.SwitchToPage("menu")
		}
	})

	flex = flex.
		AddItem(header, 0, 2, false).
		AddItem(errText, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(form, 0, 3, true).
			AddItem(table, 0, 3, false),
			0, 5, true)
	return flex
}
