package main

type UserMinimal struct {
	ID          uint
	DisplayName string
}

type RegisterRequest struct {
	Usernames   []string `json:"usernames"`
	DisplayName string   `json:"display_name"`
}
