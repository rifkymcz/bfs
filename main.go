package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alimsk/bfs/shopee"
	"github.com/charmbracelet/lipgloss"
)

var (
	cookieFilename = flag.String("cookie", "cookie", "cookie file")
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			if w, ok := err.(FatalWrapper); ok {
				log.Fatal(E(w.Error()))
			}
			panic(err)
		}
	}()

	log.SetFlags(0)
	flag.Parse()

	c := loadCookie()
	checkcookie(c)
	defer func() { saveCookie(c.Client.GetClient().Jar) }()

	deliveryAddr := addrinfo(c)
	item := iteminfo(c)
	citem := model(item)
	paymentData := payment()
	logistic := logistics(c, deliveryAddr, item)

	fmt.Println()
	clog := log.New(os.Stdout, "", log.Ltime)
	if !citem.IsFlashSale() {
		clog.Print("menunggu flash sale... ")
		fsaletime := time.Unix(citem.UpcomingFsaleStartTime(), 0)
		time.Sleep(time.Until(fsaletime))
	}

	start := time.Now()

	clog.Println("validasi checkout")
	err := c.ValidateCheckout(citem)
	fatalIf(err)
	clog.Println("checkout get")
	cget, err := c.CheckoutGet(deliveryAddr, citem, paymentData, logistic)
	fatalIf(err)
	clog.Println("place order")
	err = c.PlaceOrder(cget)
	fatalIf(err)

	spent := time.Since(start)

	clog.Println(OK("sukses"))
	fmt.Println("waktu checkout:", ternary(spent.Seconds() < 1.3, OK, W)(spent.String()))
}

func checkcookie(c shopee.Client) shopee.AccountInfo {
	fmt.Println("mengambil informasi akun ")
	acc, err := c.FetchAccountInfo()
	fatalIf(err)

	fmt.Println("masuk sebagai", H(acc.Username()))
	return acc
}

func addrinfo(c shopee.Client) shopee.AddressInfo {
	fmt.Println("mengambil informasi alamat ")
	addrs, err := c.FetchAddresses()
	fatalIf(err)
	i, deliveryAddr := addrs.DeliveryAddress()
	if i == -1 {
		fatalIf(errors.New("alamat pengiriman tidak diatur"))
	}
	return deliveryAddr
}

func iteminfo(c shopee.Client) shopee.Item {
	fmt.Println()
	url := input("url: ")
	fmt.Println("mengambil informasi item ")
	item, err := shopee.FetchItemFromURL(url)
	fatalIf(err)

	fmt.Println("nama: ", H(item.Name()))
	if !item.HasUpcomingFsale() && !item.IsFlashSale() {
		fatalIf(errors.New("tidak ada flash sale untuk item ini"))
	}
	fmt.Println("harga:", Num(strconv.FormatInt(item.Price(), 10)))
	fmt.Println("stok: ", Num(strconv.Itoa(item.Stock())))
	if item.Stock() == 0 {
		fatalIf(errors.New("stok kosong"))
	}

	return item
}

func model(item shopee.Item) shopee.CheckoutableItem {
	if len(item.Models()) == 1 {
		return shopee.ChooseModel(item, item.Models()[0].ModelID())
	}

	fmt.Println()
	fmt.Println("pilih model/varian")
	for i, m := range item.Models() {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top,
			"\n"+Num(strconv.Itoa(i))+". ",
			border(""+
				ternary(m.Stock() > 0, H, E)(m.Name())+"\n"+
				"harga: "+Num(strconv.FormatInt(m.Price(), 10))+"\n"+
				"stok: "+ternary(m.Stock() > 0, Num, E)(strconv.Itoa(m.Stock())),
			)),
		)
	}

	var i int
	for {
		i = inputInt("pilih: ", len(item.Models()))
		if item.Models()[i].Stock() != 0 {
			break
		}
		fmt.Println(E("stok varian kosong"))
	}
	return shopee.ChooseModel(item, item.Models()[i].ModelID())
}

func payment() shopee.PaymentChannelData {
	fmt.Println()
	fmt.Println("pilih metode pembayaran")
	for i, p := range shopee.PaymentChannelList {
		fmt.Print(Num(strconv.Itoa(i)), ". ", H(p.Name), "\n")
	}

	paymentCh := shopee.PaymentChannelList[inputInt("pilih: ", len(shopee.PaymentChannelList))]

	if len(paymentCh.Options) > 0 {
		fmt.Println()
		for i, opt := range paymentCh.Options {
			fmt.Print(Num(strconv.Itoa(i)), ". ", H(opt.Name), "\n")
		}
		i := inputInt("pilih: ", len(paymentCh.Options))
		return paymentCh.ApplyOpt(paymentCh.Options[i])
	} else {
		return paymentCh.Apply()
	}
}

func logistics(c shopee.Client, addr shopee.AddressInfo, item shopee.Item) shopee.LogisticChannelInfo {
	fmt.Println()
	fmt.Println("mengambil informasi logistik ")
	logistics, err := c.FetchShippingInfo(addr, item)
	fatalIf(err)

	if len(logistics) > 1 {
		for {
			fmt.Println("pilih channel logistik")
			for i, lc := range logistics {
				if lc.HasWarning() {
					fmt.Print(Num(strconv.Itoa(i)), ". ", E(lc.Name()), "\n")
				} else {
					fmt.Print(Num(strconv.Itoa(i)), ". ", OK(lc.Name()), "\n")
				}
			}
			i := inputInt("pilih: ", len(logistics))
			if l := logistics[i]; !l.HasWarning() {
				return l
			}
			fmt.Println(E("channel tersebut tidak bisa digunakan, pilih channel lain"))
			fmt.Println()
		}
	} else if len(logistics) == 1 {
		l := logistics[0]
		fmt.Println("channel", H(l.Name()), "dipilih secara otomatis")
		if l.HasWarning() {
			fatalIf(errors.New("channel " + H(l.Name()) + " tidak bisa digunakan"))
		}
		return l
	} else {
		fatalIf(errors.New("tidak ada channel logistik tersedia"))
	}

	panic("unreachable")
}

func saveCookie(jar http.CookieJar) {
	f, err := os.Create(*cookieFilename)
	fatalIf(err)
	defer f.Close()

	err = gob.NewEncoder(f).Encode(jar.Cookies(shopee.ShopeeUrl))
	fatalIf(err)
}

func loadCookie() (c shopee.Client) {
	data, err := os.ReadFile(*cookieFilename)
	if os.IsNotExist(err) {
		err = errors.New("file cookie tidak ditemukan")
	}
	fatalIf(err)

	cookieReader := bytes.NewBuffer(data)

	var cookies []*http.Cookie
	if err = gob.NewDecoder(cookieReader).Decode(&cookies); err == nil {
		jar, _ := cookiejar.New(nil)
		jar.SetCookies(shopee.ShopeeUrl, cookies)
		c, err = shopee.New(jar)
		fatalIf(err)
	} else {
		c, err = shopee.NewFromCookieString(strings.TrimSpace(cookieReader.String()))
		fatalIf(err)
	}

	return
}

type FatalWrapper struct{ error }

func fatalIf(err error) {
	if err != nil {
		panic(FatalWrapper{err})
	}
}

func input(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		os.Exit(1)
	}
	fatalIf(scanner.Err())
	return scanner.Text()
}

// 0 <= i < max
func inputInt(prompt string, max int) (i int) {
	var err error
	for {
		i, err = strconv.Atoi(input(prompt))
		if err != nil {
			fmt.Println(E("masukkan angka!"))
		} else if !(0 <= i && i < max) {
			fmt.Println(E("masukkan angka dari 0 sampai " + strconv.Itoa(max-1)))
		} else {
			return
		}
	}
}
