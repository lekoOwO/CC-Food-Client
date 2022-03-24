package main

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
)

func writeToLog(log string) {
	f, err := os.OpenFile("/tmp/ccf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	if _, err := f.WriteString(log); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func (p *Products) GetProductNames() []string {
	var names []string
	for _, product := range p.Products {
		names = append(names, product.Name)
	}
	return names
}

func (p *Products) GetProductByBarcode(barcode string) *Product {
	for _, product := range p.Products {
		if product.Barcode == barcode {
			return &product
		}
	}
	return nil
}

func (p *Products) GetProductByID(id uint64) *Product {
	for _, product := range p.Products {
		if product.ID == id {
			return &product
		}
	}
	return nil
}

func (p *Products) GetIndexByBarcode(barcode string) int {
	for i, product := range p.Products {
		if product.Barcode == barcode {
			return i
		}
	}
	return -1
}

func GetBarcode(id int) string {
	return fmt.Sprintf("CC-Food-%d", id)
}

func modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}
