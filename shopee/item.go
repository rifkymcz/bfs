package shopee

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

var ErrUrlNotMatch = errors.New("url doesn't match")

var (
	// link dari apk android
	prodLinkAppRe = regexp.MustCompile(`https?://shopee\.[\w.]+/product/(\d+)/(\d+)`)
	// link dari web
	prodLinkWebRe = regexp.MustCompile(`https?://shopee\.[\w.]+/.+\.(\d+)\.(\d+)`)
)

// parse product url
func ParseProdURL(urlstr string) (shopid, itemid int64, err error) {
	for _, re := range [...]*regexp.Regexp{prodLinkAppRe, prodLinkWebRe} {
		if match := re.FindStringSubmatch(urlstr); len(match) > 0 {
			shopid, _ = strconv.ParseInt(match[1], 10, 64)
			itemid, _ = strconv.ParseInt(match[2], 10, 64)
			return
		}
	}

	return 0, 0, ErrUrlNotMatch
}

// shopee product item
type Item struct {
	json        jsoniter.Any
	modelsCache []Model
}

func FetchItem(shopid, itemid int64) (Item, error) {
	req, err := http.NewRequest("GET", "https://mall.shopee.co.id/api/v4/item/get", nil)
	if err != nil {
		return Item{}, err
	}

	q := req.URL.Query()
	q.Set("itemid", strconv.FormatInt(itemid, 10))
	q.Set("shopid", strconv.FormatInt(shopid, 10))
	req.URL.RawQuery = q.Encode()

	for k, v := range map[string]string{
		"referer":           "https://mall.shopee.co.id",
		"x-api-source":      "rn",
		"x-shopee-language": "id",
		"if-none-match-":    "*",
		"user-agent":        ua,
	} {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Item{}, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Item{}, err
	}
	data := jsoniter.Get(b)
	if err := data.Get("error").GetInterface(); err != nil {
		return Item{}, fmt.Errorf("%v", err)
	}
	return Item{data.Get("data"), nil}, nil
}

func FetchItemFromURL(urlstr string) (Item, error) {
	shopid, itemid, err := ParseProdURL(urlstr)
	if err != nil {
		return Item{}, err
	}
	return FetchItem(shopid, itemid)
}

func (i Item) ShopID() int64          { return i.json.Get("shopid").ToInt64() }
func (i Item) ItemID() int64          { return i.json.Get("itemid").ToInt64() }
func (i Item) PriceMin() int64        { return i.json.Get("price_min").ToInt64() }
func (i Item) PriceMax() int64        { return i.json.Get("price_max").ToInt64() }
func (i Item) Price() int64           { return i.json.Get("price").ToInt64() }
func (i Item) Stock() int             { return i.json.Get("stock").ToInt() }
func (i Item) Name() string           { return i.json.Get("name").ToString() }
func (i Item) IsFlashSale() bool      { return i.json.Get("flash_sale").GetInterface() != nil }
func (i Item) HasUpcomingFsale() bool { return i.json.Get("upcoming_flash_sale").GetInterface() != nil }
func (i Item) UpcomingFsaleStartTime() int64 {
	return i.json.Get("upcoming_flash_sale", "start_time").ToInt64()
}

type Model struct{ json jsoniter.Any }

func (i *Item) Models() []Model {
	if i.modelsCache != nil {
		return i.modelsCache
	}

	models := i.json.Get("models")
	out := make([]Model, models.Size())
	for i := 0; i < models.Size(); i++ {
		out[i] = Model{models.Get(i)}
	}
	i.modelsCache = out
	return out
}

func (m Model) ItemID() int64  { return m.json.Get("itemid").ToInt64() }
func (m Model) Name() string   { return m.json.Get("name").ToString() }
func (m Model) Stock() int     { return m.json.Get("stock").ToInt() }
func (m Model) ModelID() int64 { return m.json.Get("modelid").ToInt64() }
func (m Model) Price() int64   { return m.json.Get("price").ToInt64() }

type CheckoutableItem struct {
	Item
	chosenModel int
}

func (c CheckoutableItem) ChosenModel() Model { return c.Models()[c.chosenModel] }

func ChooseModel(item Item, modelId int64) CheckoutableItem {
	var modelIndex int
	for i, m := range item.Models() {
		if m.ModelID() == modelId {
			modelIndex = i
			break
		}
	}
	return CheckoutableItem{item, modelIndex}
}
