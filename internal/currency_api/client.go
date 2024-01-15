package currency_api

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	apiV1Url = "https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1"

	sep        = "-"
	dateFormat = "2006-01-02"
)

type Client struct {
	log           *zap.SugaredLogger
	client        *http.Client
	currencyCache map[string]float64 // key: usd-eur-2023-12-30, value: cash
	allCurrencies map[string]struct{}

	rwMut *sync.RWMutex
}

func New(log *zap.SugaredLogger, allCurrencies map[string]struct{}) *Client {
	return &Client{
		log: log,
		client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   time.Second * 10,
		},
		currencyCache: make(map[string]float64),
		allCurrencies: allCurrencies,
		rwMut:         &sync.RWMutex{},
	}
}

func (c *Client) GetCurrency(from, to string, date time.Time) (float64, error) {
	toValue := c.getCurrencyFromCache(from, to, date)
	if toValue > 0 {
		return toValue, nil
	}

	var dateStr string
	now := time.Now().UTC().Round(time.Hour * 24)
	if date.Round(time.Hour * 24).Equal(now) {
		dateStr = "latest"
	} else {
		dateStr = date.Format(dateFormat)
	}

	url := apiV1Url + "/" + dateStr + "/currencies/" + strings.ToLower(from) + "/" + strings.ToLower(to) + ".json"
	resp, err := c.client.Get(url)
	if err != nil {
		return 0, errors.Wrap(err, "get currency")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.Errorf("status code: %d, url: %s", resp.StatusCode, url)
	}

	respData := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return 0, errors.Wrap(err, "decode json")
	}

	switch v := respData[to].(type) {
	case float64:
		c.setCurrencyToCache(from, to, v, date)
		return v, nil
	}

	return 0, errors.New("unknown currency")
}

func (c *Client) CheckCurrency(cur string) bool {
	_, ok := c.allCurrencies[strings.ToUpper(cur)]
	return ok
}

func (c *Client) getCurrencyFromCache(from, to string, date time.Time) float64 {
	c.rwMut.RLock()
	defer c.rwMut.RUnlock()

	v, ok := c.currencyCache[from+sep+to+sep+date.Format(dateFormat)]
	if ok {
		return v
	}

	return 0
}

func (c *Client) setCurrencyToCache(from, to string, cash float64, date time.Time) {
	c.rwMut.Lock()
	c.currencyCache[from+sep+to+sep+date.Format(dateFormat)] = cash
	c.rwMut.Unlock()
}
