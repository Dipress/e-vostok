package main

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// getBalance gets balance via e-vastok.ru
func getBalance(url string) string {
	res, err := http.Get(url)

	if err != nil {
		errors.Wrapf(err, "get request failed")
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		errors.Wrapf(err, "collect data failed")
	}
	return string(result)
}
