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
