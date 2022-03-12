package shopee

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
)

var ShopeeUrl, _ = url.Parse("https://mall.shopee.co.id")

const ua = "Android app Shopee appver=28320 app_type=1"

type Client struct {
	Client *resty.Client
}

func New(cookie http.CookieJar) (Client, error) {
	var csrftoken string
	for _, c := range cookie.Cookies(ShopeeUrl) {
		if c.Name == "csrftoken" {
			csrftoken = c.Value
			break
		}
	}
	if csrftoken == "" {
		return Client{}, errors.New("csrftoken not found in cookie")
	} else if len(csrftoken) != 32 {
		return Client{}, errors.New("invalid csrftoken")
	}

	client := resty.New().
		SetCookieJar(cookie).
		SetBaseURL(ShopeeUrl.String()).
		SetHeaders(map[string]string{
			"Referer":           ShopeeUrl.String(),
			"If-None-Match-":    "*",
			"X-Api-Source":      "rn",
			"X-Shopee-Language": "id",
			"User-Agent":        ua,
			"Content-Type":      "application/json",
			"Accept":            "application/json",
			"X-Csrftoken":       csrftoken,
		})
	client.JSONMarshal = jsoniter.Marshal
	client.JSONUnmarshal = jsoniter.Unmarshal
	// client.SetDebug(true) // uncomment for debugging
	return Client{client}, nil
}

func NewFromCookieString(cookie string) (Client, error) {
	h := http.Header{}
	h.Add("Cookie", cookie)
	r := http.Request{Header: h}
	c, _ := cookiejar.New(nil)
	c.SetCookies(ShopeeUrl, r.Cookies())

	return New(c)
}
