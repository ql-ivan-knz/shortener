package storage

import (
	"errors"
	"fmt"
)

var errNoExistedURL = errors.New("url_store: no existed url")

var store = make(map[string]string)

func Get(key string) (string, error) {
	if _, ok := store[key]; ok {
		return store[key], nil
	}

	return "", errNoExistedURL
}

func Set(key, value string) {
	store[key] = value
}

func GetAll() {
	fmt.Println(store)
}
