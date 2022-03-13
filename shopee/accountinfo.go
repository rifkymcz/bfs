package shopee

import (
	jsoniter "github.com/json-iterator/go"
)

type AccountInfo struct{ json jsoniter.Any }

func (c Client) FetchAccountInfo() (AccountInfo, error) {
	resp, err := c.Client.R().
		SetQueryParams(map[string]string{
			"need_cart":    "0",
			"skip_address": "1",
		}).
		Get("/api/v1/account_info")
	if err != nil {
		return AccountInfo{}, err
	}

	json := jsoniter.Get(resp.Body())
	if json.Size() == 0 {
		return AccountInfo{}, InvalidCookieError("invalid or expired cookie")
	}

	return AccountInfo{json}, nil
}

func (a AccountInfo) ShopID() int64    { return a.json.Get("shopid").ToInt64() }
func (a AccountInfo) Username() string { return a.json.Get("username").ToString() }
func (a AccountInfo) UserID() int64    { return a.json.Get("userid").ToInt64() }

type AddressInfo struct {
	json           jsoniter.Any
	isDeliveryAddr bool
}

type Addresses []AddressInfo

func (c Client) FetchAddresses() (Addresses, error) {
	resp, err := c.Client.R().Get("/api/v1/addresses")
	if err != nil {
		return nil, err
	}

	json := jsoniter.Get(resp.Body())
	deliveryAddrId := json.Get("delivery_address_id").ToInt64()
	addrs := json.Get("addresses")
	out := make(Addresses, addrs.Size())
	for i := 0; i < addrs.Size(); i++ {
		out[i] = AddressInfo{
			json:           addrs.Get(i),
			isDeliveryAddr: addrs.Get(i, "id").ToInt64() == deliveryAddrId,
		}
	}
	return out, nil
}

// returns index -1 if not found
func (a Addresses) DeliveryAddress() (int, AddressInfo) {
	for i, addr := range a {
		if addr.IsDeliveryAddress() {
			return i, addr
		}
	}
	return -1, AddressInfo{}
}

func (a AddressInfo) City() string            { return a.json.Get("city").ToString() }
func (a AddressInfo) District() string        { return a.json.Get("district").ToString() }
func (a AddressInfo) State() string           { return a.json.Get("state").ToString() }
func (a AddressInfo) Country() string         { return a.json.Get("country").ToString() }
func (a AddressInfo) Town() string            { return a.json.Get("town").ToString() }
func (a AddressInfo) Address() string         { return a.json.Get("address").ToString() }
func (a AddressInfo) Zipcode() string         { return a.json.Get("zipcode").ToString() }
func (a AddressInfo) GeoString() string       { return a.json.Get("geoString").ToString() }
func (a AddressInfo) ID() int64               { return a.json.Get("id").ToInt64() }
func (a AddressInfo) IsDeliveryAddress() bool { return a.isDeliveryAddr }
