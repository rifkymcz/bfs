package shopee

import (
	"fmt"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

type LogisticChannelInfo struct{ json jsoniter.Any }

func (c Client) FetchShippingInfo(addr AddressInfo, item Item) ([]LogisticChannelInfo, error) {
	resp, err := c.Client.R().
		SetQueryParams(map[string]string{
			"buyer_zipcode": addr.Zipcode(),
			"city":          addr.City(),
			"district":      addr.District(),
			"itemid":        strconv.FormatInt(item.ItemID(), 10),
			"shopid":        strconv.FormatInt(item.ShopID(), 10),
			"state":         addr.State(),
			"town":          addr.Town(),
		}).
		Get("/api/v4/pdp/get_shipping")
	if err != nil {
		return nil, err
	}

	json := jsoniter.Get(resp.Body())

	if err := json.Get("error").GetInterface(); err != nil {
		return nil, fmt.Errorf("%v: %v", json.Get("error").GetInterface(), json.Get("error_msg").GetInterface())
	}

	channels := json.Get("data", "ungrouped_channel_infos")
	out := make([]LogisticChannelInfo, channels.Size())
	for i := 0; i < channels.Size(); i++ {
		out[i] = LogisticChannelInfo{channels.Get(i)}
	}

	return out, nil
}

func (c LogisticChannelInfo) ChannelID() int64 { return c.json.Get("channel_id").ToInt64() }
func (c LogisticChannelInfo) Name() string     { return c.json.Get("name").ToString() }
func (c LogisticChannelInfo) HasWarning() bool { return c.json.Get("warning").GetInterface() != nil }
func (c LogisticChannelInfo) Warning() string  { return c.json.Get("warning", "warning_msg").ToString() }
