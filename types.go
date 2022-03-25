package main

import "time"

type UserMinimal struct {
	ID          uint
	DisplayName string
}

type User struct {
	UserMinimal
	Usernames []Username
}

type Username struct {
	ID     uint64
	Name   string
	UserID uint64
}

type RegisterRequest struct {
	Usernames   []string `json:"usernames"`
	DisplayName string   `json:"display_name"`
}

type BuyRequestDetail struct {
	ProductID uint64
	Quantity  uint64
}

type BuyRequest struct {
	UserID  uint
	Details []BuyRequestDetail
}

type NewProductRequest struct {
	Name    string `json:"name"`
	Price   int64  `json:"price"`
	Barcode string `json:"barcode"`
}

type NewUsernameRequest struct {
	Name string `json:"name"`
}

type Product struct {
	ID      uint64
	Name    string
	Price   uint64
	Barcode string
}

type Products struct {
	Products []Product
}

type CartItem struct {
	ProductID uint64
	Quantity  uint64
}

type PurchaseDetail struct {
	ID        uint64
	ProductID uint64
	Quantity  uint64
	Total     uint64
}

type Purchase struct {
	ID              uint64
	UserID          uint64
	CreatedAt       time.Time
	PurchaseDetails []PurchaseDetail
}

type PayRequest struct {
	UserID      uint     `json:"user_id"`
	PurchaseIDs []uint64 `json:"purchase_ids"`
}
