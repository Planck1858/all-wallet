package tinkoff

import (
	"context"
	"net/http"

	sdk "github.com/ssummers02/invest-api-go-sdk"
)

const (
	defaultUrl = "hinvest-public-api.tinkoff.ru:443"
)

type TinkoffClient struct {
	token, url string
	client     *http.Client
}

func New(token string) *TinkoffClient {
	return &TinkoffClient{
		token: token,
		url:   defaultUrl,
		client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   15,
		},
	}
}

func (t *TinkoffClient) GetAccounts(ctx context.Context) ([]Account, error) {
	sdk.

	return acc, err
}

func (t *TinkoffClient) doRequest(url, method string, body interface{}) (resp interface{}, respCode int, err error) {
	return resp, respCode, nil
}
