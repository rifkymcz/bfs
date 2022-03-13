package shopee

import (
	_ "embed"
	"fmt"
	"log"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type CheckoutValidationError struct{ json jsoniter.Any }

func (c CheckoutValidationError) Code() int   { return c.json.Get("error").ToInt() }
func (c CheckoutValidationError) Msg() string { return c.json.Get("error_msg").ToString() }
func (c CheckoutValidationError) Error() string {
	return fmt.Sprintf("checkout validation error: code=%d, validation_code=%d, %s", c.Code(), c.ValidationCode(), c.Msg())
}
func (c CheckoutValidationError) ValidationCode() int {
	return c.json.Get("data", "validation_error").ToInt()
}

func (c Client) ValidateCheckout(item CheckoutableItem) error {
	type obj = map[string]interface{}
	type arr = []interface{}
	resp, err := c.Client.R().
		SetBody(obj{
			"shop_orders": arr{
				obj{
					"shop_info": obj{"shop_id": item.ShopID()},
					"item_infos": arr{
						obj{
							"item_id":  item.ItemID(),
							"model_id": item.ChosenModel().ModelID(),
							"quantity": 1,
						},
					},
					"buyer_address": nil,
				},
			},
		}).
		Post("/api/v4/pdp/buy_now/validate_checkout")
	if err != nil {
		return err
	}

	json := jsoniter.Get(resp.Body())
	if json.Get("error").ToInt() != 0 || json.Get("data", "validation").ToInt() != 0 {
		return CheckoutValidationError{json}
	}

	return nil
}

type PlaceOrderError struct{ json jsoniter.Any }

func (p PlaceOrderError) Type() string { return p.json.Get("error").ToString() }
func (p PlaceOrderError) Msg() string  { return p.json.Get("error_msg").ToString() }
func (p PlaceOrderError) Error() string {
	return fmt.Sprintf("%s: %s", p.json.Get("error").ToString(), p.json.Get("error_msg").ToString())
}

//go:embed place_order.json
var placeOrderPayload []byte

func (c Client) PlaceOrder(r CheckoutGetResult) error {
	resp, err := c.Client.R().
		SetBody(r).
		Post("/api/v4/checkout/place_order")
	if err != nil {
		return err
	}

	json := jsoniter.Get(resp.Body())
	if json.Get("error").ToString() != "" {
		return PlaceOrderError{json}
	}
	return nil
}

// TODO
func (c Client) checkoutGetQuick() {
	panic("TODO")
}

//go:embed checkout_get.json
var checkoutGetPayload []byte

type CheckoutGetResult struct {
	json jsoniter.Any
	ts   int64
}

func (c CheckoutGetResult) MarshalJSON() ([]byte, error) {
	data, err := jsoniter.Marshal(c.json)
	if err != nil {
		return data, err
	}
	return append(data[:len(data)-1], (`,"timestamp":` + strconv.FormatInt(c.ts, 10) + "}")...), nil
}

func (c Client) CheckoutGet(
	addr AddressInfo,
	item CheckoutableItem,
	payment PaymentChannelData,
	logistic LogisticChannelInfo,
) (CheckoutGetResult, error) {
	var data map[string]interface{}
	if err := jsoniter.Unmarshal(checkoutGetPayload, &data); err != nil {
		return CheckoutGetResult{}, err
	}

	ts := time.Now().Unix()

	type p = []interface{}
	for _, field := range []struct {
		path p
		v    interface{}
	}{
		{p{"timestamp"}, ts},
		{p{"shoporders", 0, "shop", "shopid"}, item.ShopID()},
		{p{"shoporders", 0, "items", 0, "itemid"}, item.ItemID()},
		{p{"shoporders", 0, "items", 0, "modelid"}, item.ChosenModel().ModelID()},
		{p{"selected_payment_channel_data"}, payment},
		{p{"shipping_orders", 0, "buyer_address_data", "addressid"}, addr.ID()},
		{p{"shipping_orders", 0, "selected_logistic_channelid"}, logistic.ChannelID()},
	} {
		setJsonField(data, field.path, field.v)
	}

	resp, err := c.Client.R().
		SetBody(data).
		Post("/api/v4/checkout/get")
	if err != nil {
		return CheckoutGetResult{}, err
	}

	return CheckoutGetResult{jsoniter.Get(resp.Body()), ts}, nil
}

func setJsonField(json interface{}, path []interface{}, v interface{}) {
	for i, accessor := range path[:len(path)-1] {
		switch jsontyp := json.(type) {
		case map[string]interface{}:
			json = jsontyp[accessor.(string)]
		case []interface{}:
			json = jsontyp[accessor.(int)]
		default:
			log.Printf("path: %#+v\n", path)
			log.Printf("json: %#+v\n", json)
			panic(fmt.Sprint("invalid accessor: ", accessor, " at index ", i))
		}
	}
	field := path[len(path)-1]
	switch json := json.(type) {
	case map[string]interface{}:
		json[field.(string)] = v
	case []interface{}:
		json[field.(int)] = v
	}
}

func getJsonField(json interface{}, path []interface{}) jsoniter.Any {
	for i, accessor := range path[:len(path)-1] {
		switch jsontyp := json.(type) {
		case map[string]interface{}:
			json = jsontyp[accessor.(string)]
		case []interface{}:
			json = jsontyp[accessor.(int)]
		default:
			log.Printf("path: %#+v\n", path)
			log.Printf("json: %#+v\n", json)
			panic(fmt.Sprint("invalid accessor: ", accessor, " at index ", i))
		}
	}
	field := path[len(path)-1]
	switch json := json.(type) {
	case map[string]interface{}:
		return jsoniter.Wrap(json[field.(string)])
	case []interface{}:
		return jsoniter.Wrap(json[field.(int)])
	}
	return nil
}
