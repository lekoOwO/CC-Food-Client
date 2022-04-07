package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func GetUserByUsername(username string) (*UserMinimal, error) {
	res, err := http.Get(APIEndPoint + "/user/byUsername/" + url.QueryEscape(username))

	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, error(fmt.Errorf("%d", res.StatusCode))
	}

	var data UserMinimal

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func GetUserByID(id uint) (*User, error) {
	res, err := http.Get(APIEndPoint + "/user/" + fmt.Sprintf("%d", id))

	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, error(fmt.Errorf("%d", res.StatusCode))
	}

	var data User

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func Register(displayName string, usernames []string) (uint, error) {
	data := RegisterRequest{
		DisplayName: displayName,
		Usernames:   usernames,
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(data)
	if err != nil {
		return 0, err
	}

	res, err := http.Post(APIEndPoint+"/user", "application/json", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != 200 {
		return 0, error(fmt.Errorf("%d", res.StatusCode))
	}

	var minUser UserMinimal

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&minUser)
	if err != nil {
		return 0, err
	}

	return minUser.ID, nil
}

func GetProductByID(id uint64) (*Product, error) {
	res, err := http.Get(APIEndPoint + "/product/" + url.QueryEscape(fmt.Sprintf("%d", id)))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, error(fmt.Errorf("%d", res.StatusCode))
	}

	var data Product

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func GetProducts() *Products {
	res, err := http.Get(APIEndPoint + "/product")

	if err != nil {
		return nil
	}
	if res.StatusCode != 200 {
		return nil
	}

	var data []Product

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return nil
	}

	return &Products{Products: data}
}

func Buy(buyRequest BuyRequest) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(buyRequest)
	if err != nil {
		return err
	}

	res, err := http.Post(APIEndPoint+"/purchase", "application/json", b)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return error(fmt.Errorf("%d", res.StatusCode))
	}

	return nil
}

func CreateProduct(req NewProductRequest) (*Product, error) {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(req)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(APIEndPoint+"/product", "application/json", b)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, error(fmt.Errorf("%d", res.StatusCode))
	}

	var data Product

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func EditProduct(id uint64, req NewProductRequest) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(req)
	if err != nil {
		return err
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPut, APIEndPoint+"/product/"+url.QueryEscape(fmt.Sprintf("%d", id)), b)
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return error(fmt.Errorf("%d", res.StatusCode))
	}

	return nil
}

func DelectProduct(id uint64) error {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodDelete, APIEndPoint+"/product/"+url.QueryEscape(fmt.Sprintf("%d", id)), nil)
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return error(fmt.Errorf("%d", res.StatusCode))
	}

	return nil
}

func GetUnpaidPurchasesByUserID(userID uint) ([]Purchase, error) {
	res, err := http.Get(APIEndPoint + "/purchase/notPaid/" + url.QueryEscape(fmt.Sprintf("%d", userID)))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, error(fmt.Errorf("%d", res.StatusCode))
	}

	var data []Purchase

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetUnpaidPurchases() ([]Purchase, error) {
	res, err := http.Get(APIEndPoint + "/purchase/notPaid/")
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, error(fmt.Errorf("%d", res.StatusCode))
	}

	var data struct {
		Purchases []Purchase `json:"purchases"`
	}

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data.Purchases, nil
}

func GetPurchase(id uint64) (*Purchase, error) {
	res, err := http.Get(APIEndPoint + "/purchase/id/" + url.QueryEscape(fmt.Sprintf("%d", id)))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, error(fmt.Errorf("%d", res.StatusCode))
	}

	var data Purchase

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func pay(pr PayRequest) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(pr)
	if err != nil {
		return err
	}

	res, err := http.Post(APIEndPoint+"/purchase/pay", "application/json", b)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return error(fmt.Errorf("%d", res.StatusCode))
	}

	return nil
}

func DeleteUsername(id uint64) error {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodDelete, APIEndPoint+"/username/"+url.QueryEscape(fmt.Sprintf("%d", id)), nil)
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return error(fmt.Errorf("%d", res.StatusCode))
	}
	return nil
}

func NewUsername(id uint, username string) error {
	nur := NewUsernameRequest{
		Name: username,
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(nur)
	if err != nil {
		return err
	}

	res, err := http.Post(APIEndPoint+"/username/"+url.QueryEscape(fmt.Sprintf("%d", id)), "application/json", b)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return error(fmt.Errorf("%d", res.StatusCode))
	}

	return nil
}
