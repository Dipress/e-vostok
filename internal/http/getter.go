package http

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Balance struct {
	client *http.Client
}

func NewBalance() *Balance {
	return &Balance{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (b Balance) Get(url string) (string, error) {
	res, err := b.client.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "get request failed")
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "collect data failed")
	}

	defer res.Body.Close()

	return "The Balance is " + string(result), nil
}
