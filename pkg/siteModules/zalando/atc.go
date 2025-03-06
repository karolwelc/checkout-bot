package zalando

import (
	"encoding/json"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/ProjektPLTS/pkg/customErrors"
	tls "github.com/ProjektPLTS/pkg/tls"
)

func (task *ZalTask) atcApp() error {
	headers := tls.Headers{
		"Host":                 `www.zalando.` + strings.ToLower(task.Address.Delivery.Country),
		"x-ts":                 "1663175691502",
		"User-Agent":           `zalando/22.11.1 (iPhone; iOS 15.6.1; Scale/3.00)`,
		"x-sig":                `c04c3b41e970020373ad8c39baa0194ba30e282e`, //not checked
		"ot-tracer-traceid":    `fd805790d43d9d45`,
		"x-logged-in":          `true`,
		"x-device-type":        `smartphone`,
		"x-zalando-client-id":  `1e780ba7-e417-4f8a-b8ca-e4510129a87d`, //todo
		"x-frontend-type":      `mobile-app`,
		"x-zalando-mobile-app": `3580f92a4bafb890i`,
		"x-zalando-consent-id": `2A24E2AF-C407-47A5-B493-EA69FD0F8BD1`,
		"x-os-version":         `15.6.1`,
		"x-device-os":          `ios`,
		"x-sales-channel":      `b773b421-c719-4dfd-afc8-e97da508a88d`,
		"ot-tracer-sampled":    `true`,
		"x-app-domain":         "47",
		"x-app-version":        "22.11.1",
		"x-device-platform":    `ios`,
		"Accept-Language":      `*`,
		"ot-tracer-spanid":     `f8386bf252a3f470`,
		"x-uuid":               `3361EFE1-A3A8-4258-933B-0D8A8EE1A37A`, //todo
		"Accept":               `application/json`,
		"Content-Type":         `application/json`,
	}
	jsonBody := `{"appVersion":"22.11.1","appDomain":47,"signature":"282ce494d1ddbbf0eb1499e849fccf90e255fcc5","timeStamp":1663175691502,"devicePlatform":"ios","uuid":"3361EFE1-A3A8-4258-933B-0D8A8EE1A37A","items":[{"sku":"` + task.Product.Id + `","simpleSku":"` + task.Product.SelectedSku + `","business_partner_id":"810d1d00-4312-43e5-bd31-d8373fdd24c7"}]}`

	opts := tls.Options{
		URL:    "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/api/mobile/v3/cart.json",
		Method: "POST",

		Headers: headers,
		Body:    jsonBody,

		ClientSettings: &task.ClientSettings,
		Jar:            task.Cookies,
	}
	resp, err := tls.Do(opts)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}
	task.Cookies = resp.Request.Jar
	return nil
}

type productInfo struct {
	Data struct {
		Product struct {
			Name             string `json:"name"`
			Sku              string `json:"sku"`
			PrimaryThumbnail struct {
				URI string `json:"uri"`
			} `json:"primaryThumbnail"`
			ComingSoon  bool      `json:"comingSoon"`
			ReleaseDate time.Time `json:"releaseDate"`
			Simples     []struct {
				Sku       string `json:"sku"`
				Size      string `json:"size"`
				AllOffers []struct {
					Stock struct {
						Quantity string `json:"quantity"`
					} `json:"stock"`
				} `json:"allOffers"`
			} `json:"simples"`
			SimplesWithStock []struct {
				Sku       string `json:"sku"`
				Size      string `json:"size"`
				AllOffers []struct {
					Stock struct {
						Quantity string `json:"quantity"`
					} `json:"stock"`
				} `json:"allOffers"`
			} `json:"simplesWithStock"`
			Family struct {
				Products struct {
					Edges []struct {
						Node struct {
							StandardColorThumbnail struct {
								URI string `json:"uri"`
							} `json:"standardColorThumbnail"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"products"`
				Rating struct {
					Average float64 `json:"average"`
				} `json:"rating"`
				Reviews struct {
					TotalCount int `json:"totalCount"`
				} `json:"reviews"`
			} `json:"family"`
		} `json:"product"`
	} `json:"data"`
}

func (task *ZalTask) getProductInfo() error {
	var productPage string
	if !task.finishedPreload {
		productPage = task.DummyProductURL
	} else {
		productPage = task.Product.URL

	}
	headers := tls.Headers{
		"Host":                      "www.zalando." + task.Address.Delivery.Country,
		"cache-control":             "max-age=0",
		"sec-ch-ua":                 "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\"",
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "\"Windows\"",
		"upgrade-insecure-requests": "1",
		"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"sec-fetch-site":            "none",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-user":            "?1",
		"sec-fetch-dest":            "document",
		"accept-language":           "en,cs-CZ;q=0.9,cs;q=0.8,be;q=0.7,sk;q=0.6,zh-TW;q=0.5,zh;q=0.4,en-US;q=0.3",
	}
	opts := tls.Options{
		URL:    productPage,
		Method: "GET",

		Headers: headers,

		ClientSettings: &task.ClientSettings,
	}
	resp, err := tls.Do(opts)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}
	jsonRespStr := strings.Split(strings.Split(strings.Split(resp.Body, `d79ad15b0b125e06d5b4430437c09a6667f5a3bcad11c9ac0d7e324759f20aba\",\"variables\":{\"i`)[1], `"}}":`)[1], `,"{\"id\":\"`)[0]
	var jsonResp productInfo
	if err := json.Unmarshal([]byte(jsonRespStr), &jsonResp); err != nil {
		return err
	}
	task.Product.Name = jsonResp.Data.Product.Name
	task.Product.Id = jsonResp.Data.Product.Sku
	task.timer = jsonResp.Data.Product.ComingSoon

	task.time = jsonResp.Data.Product.ReleaseDate
	task.Product.Image = jsonResp.Data.Product.Family.Products.Edges[0].Node.StandardColorThumbnail.URI
	for _, size := range jsonResp.Data.Product.SimplesWithStock {
		task.Product.Skus = append(task.Product.Skus, size.Sku)
		task.Product.Sizes = append(task.Product.Sizes, size.Size)
		task.Product.Stock = append(task.Product.Stock, size.AllOffers[0].Stock.Quantity)
	}

	return nil
}

func (task *ZalTask) chooseSize() error {
	//TODO (KONIK): maybe add priority to MANY stock
	variants := task.Product.Skus
	for i := len(variants) - 1; i >= 0; i-- {
		if task.Product.Stock[i] == "OUT_OF_STOCK" {
			variants = append(variants[:i], variants[i+1:]...)
		}
	}
	wantedSize := task.WantedSize
	if len(variants) == 0 {
		return customErrors.OOS
	}
	for i := 0; i < len(variants); i++ {
		if task.Product.Sizes[i] == wantedSize {
			task.Product.SelectedSku = variants[i]
			task.Product.SelectedSize = task.Product.Sizes[i]
			return nil
		}
	}
	if task.RandomSizeIfOOS {
		randomIndex := rand.Intn(len(variants))
		task.Product.SelectedSku = variants[randomIndex]
		task.Product.SelectedSize = task.Product.Sizes[randomIndex]
		return nil
	} else {
		return customErrors.OOS
	}
}

func (task *ZalTask) addToCart() error {
	//only get product info if empty
	if task.Product.SelectedSize == "" {
		err := task.getProductInfo()
		if err != nil {
			return err
		}
		err = task.chooseSize()
		if err != nil {
			return err
		}
	}
	err := task.atcApp()
	if err != nil {
		return err
	}
	return nil
}

type cartResp struct {
	ID     string `json:"id"`
	Groups []struct {
		Articles []struct {
			SimpleSku string `json:"simple_sku"`
		} `json:"articles"`
	} `json:"groups"`
}

func (task *ZalTask) getCartedProducts() ([]string, error) {
	var xXsrf string
	url, _ := url.Parse("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/api/cart-gateway/carts")
	for _, c := range task.Cookies.Cookies(url) {
		if c.Name == "frsx" {
			xXsrf = c.Value
		}
	}
	headers := tls.Headers{
		"Host":               "www.zalando." + strings.ToLower(task.Address.Delivery.Country),
		"sec-ch-ua":          "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\"",
		"accept":             "application/json",
		"content-type":       "application/json",
		"x-xsrf-token":       xXsrf,
		"sec-ch-ua-mobile":   "?0",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		"sec-ch-ua-platform": "\"Windows\"",
		"origin":             "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country),
		"sec-fetch-site":     "same-origin",
		"sec-fetch-mode":     "cors",
		"sec-fetch-dest":     "empty",
		"referer":            "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/cart",
		"accept-language":    "en,cs-CZ;q=0.9,cs;q=0.8,be;q=0.7,sk;q=0.6,zh-TW;q=0.5,zh;q=0.4,en-US;q=0.3",
	}

	postData := `{}`

	opts := tls.Options{
		URL:    "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/api/cart-gateway/carts",
		Method: "POST",

		Headers: headers,

		Body: postData,

		ClientSettings: &task.ClientSettings,
		Jar:            task.Cookies,
	}

	resp, err := tls.Do(opts)
	if err != nil {
		return []string{}, err
	}
	if resp.StatusCode != 200 {
		return []string{}, customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}
	var jsonResp cartResp
	if err := json.Unmarshal([]byte(resp.Body), &jsonResp); err != nil {
		return []string{}, err
	}
	task.cartId = jsonResp.ID
	var products []string

	if len(jsonResp.Groups) == 0 {
		return []string{}, nil
	}

	for _, p := range jsonResp.Groups[0].Articles {
		products = append(products, p.SimpleSku)
	}

	task.Cookies = resp.Request.Jar
	return products, nil
}

func (task *ZalTask) deleteProducts(products []string) error {
	var xXsrf string
	headers := tls.Headers{
		"Host":         "www.zalando." + strings.ToLower(task.Address.Delivery.Country),
		"sec-ch-ua":    "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\"",
		"accept":       "application/json",
		"content-type": "application/json",

		"sec-ch-ua-mobile":   "?0",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		"sec-ch-ua-platform": "\"Windows\"",
		"origin":             "www.zalando." + strings.ToLower(task.Address.Delivery.Country),
		"sec-fetch-site":     "same-origin",
		"sec-fetch-mode":     "cors",
		"sec-fetch-dest":     "empty",
		"referer":            "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/cart",
		"accept-language":    "en,cs-CZ;q=0.9,cs;q=0.8,be;q=0.7,sk;q=0.6,zh-TW;q=0.5,zh;q=0.4,en-US;q=0.3",
	}
	opts := tls.Options{
		Method: "DELETE",

		ClientSettings: &task.ClientSettings,

		Jar: task.Cookies,
	}
	for _, product := range products {
		url, _ := url.Parse("https://www.zalando." + task.Address.Delivery.Country + "/api/cart-gateway/carts/" + task.cartId + "/items/" + product)
		for _, c := range task.Cookies.Cookies(url) {
			if c.Name == "frsx" {
				xXsrf = c.Value
			}
		}

		headers["x-xsrf-token"] = xXsrf
		opts.URL = "https://www.zalando." + task.Address.Delivery.Country + "/api/cart-gateway/carts/" + task.cartId + "/items/" + product
		opts.Headers = headers

		resp, err := tls.Do(opts)
		if err != nil {
			return err
		}
		if resp.StatusCode != 204 {
			return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
		}
		task.Cookies = resp.Request.Jar
	}
	return nil
}
