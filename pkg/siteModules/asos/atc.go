package asos

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/ProjektPLTS/pkg/customErrors"
	tls "github.com/ProjektPLTS/pkg/tls"
	"github.com/google/uuid"
)

type productInfo []struct {
	ProductID   int    `json:"productId"`
	ProductCode string `json:"productCode"`
	Variants    []struct {
		ID        int  `json:"id"`
		VariantID int  `json:"variantId"`
		IsInStock bool `json:"isInStock"`
	} `json:"variants"`
}

func (task *AsosTask) getProductInfo() error {
	task.Product.Skus = []string{}
	headers := tls.Headers{
		"Host":               "www.asos.com",
		"sec-ch-ua":          "\"Google Chrome\";v=\"105\", \"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"105\"",
		"asos-c-version":     "1.0.0-99ece1f62a83-9662",
		"sec-ch-ua-mobile":   "?0",
		"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
		"asos-c-name":        "asos-web-productpage",
		"sec-ch-ua-platform": "\"macOS\"",
		"accept":             "*/*",
		"sec-fetch-site":     "same-origin",
		"sec-fetch-mode":     "cors",
		"sec-fetch-dest":     "empty",
		"referer":            "https://www.asos.com/vans/vans-theory-sk8-hi-tapered-trainers-in-black-white/prd/203116306?ctaref=we+recommend+grid_1&featureref1=we+recommend+pers",
		"accept-language":    "en-GB,en-US;q=0.9,en;q=0.8",
	}

	opts := tls.Options{
		URL:    "https://www.asos.com/api/product/catalogue/v3/stockprice?productIds=" + task.Product.URL + "&store=" + task.storeCode + "&konik=" + fmt.Sprint(time.Now().UnixMilli()),
		Method: "GET",

		Headers:        headers,
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
	var productInfo productInfo
	if err := json.Unmarshal([]byte(resp.Body), &productInfo); err != nil {
		return err
	}
	for _, v := range productInfo[0].Variants {
		if v.IsInStock {
			task.Product.Skus = append(task.Product.Skus, strconv.Itoa(v.VariantID))
		}
	}
	if len(task.Product.Skus) == 0 {
		fmt.Println(task.TaskPrefix() + "OOS monitoring, sleeping for " + strconv.Itoa(task.Config.MonitorDelay) + "ms...")
		time.Sleep(time.Duration(task.Config.MonitorDelay) * time.Millisecond)
		return task.getProductInfo()
	}
	randIndx := rand.Intn(len(task.Product.Skus))
	task.Product.SelectedSku = task.Product.Skus[randIndx]
	return nil
}

func (task *AsosTask) atc() error {
	err := task.getProductInfo()
	if err != nil {
		return err
	}
	headers := tls.Headers{
		"Host":               "www.asos.com",
		"sec-ch-ua":          "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\"",
		"asos-c-plat":        "Web",
		"asos-bag-origin":    "EUW",
		"authorization":      "Bearer " + task.bearerToken,
		"asos-c-name":        "Asos.Commerce.Bag.Sdk",
		"asos-c-store":       "ROE",
		"x-requested-with":   "XMLHttpRequest",
		"sec-ch-ua-platform": "\"macOS\"",
		"sec-ch-ua-mobile":   "?0",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		"content-type":       "application/json",
		"accept":             "application/json, text/javascript, */*; q=0.01",
		"asos-c-ismobile":    "false",
		"asos-c-ver":         "5.5.195",
		"asos-c-istablet":    "false",
		"asos-cid":           uuid.NewString(),
		"origin":             "https://www.asos.com",
		"sec-fetch-site":     "same-origin",
		"sec-fetch-mode":     "cors",
		"sec-fetch-dest":     "empty",
		"accept-language":    "en,cs-CZ;q=0.9,cs;q=0.8,be;q=0.7,sk;q=0.6,zh-TW;q=0.5,zh;q=0.4,en-US;q=0.3",
	}

	postData := `{"variantId":` + task.Product.SelectedSku + `}`
	fmt.Println(task.Cookies)
	opts := tls.Options{
		URL:    "https://www.asos.com/api/commerce/bag/v4/bags/" + task.bagId + "/product?expand=summary,total&lang=en-GB",
		Method: "POST",

		Headers: headers,

		Body: postData,

		ClientSettings: &task.ClientSettings,
		Jar:            task.Cookies,
	}

	resp, err := tls.Do(opts)
	if err != nil {
		return err
	}
	if !(resp.StatusCode == 201 || resp.StatusCode == 200) {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}
	task.Cookies = resp.Request.Jar
	fmt.Println("-----BODY-----")
	fmt.Println(resp.Body)
	fmt.Println("-----REQUEST HEADERS-----")
	fmt.Println(resp.Request.Headers)
	fmt.Println("-----RESPONSE HEADERS-----")
	fmt.Println(resp.Headers)
	fmt.Println("-----COOKIES-----")
	url, _ := url.Parse("https://my.asos.com/identity")
	for _, v := range task.Cookies.Cookies(url) {
		if v.Name == "asos-anon12" {
			fmt.Println(v.Name)
			fmt.Println(v.Value)
		} else if v.Name == "asos-ts121" {
			fmt.Println(v.Name)
			fmt.Println(v.Value)
		} else if v.Name == "idsvr.session" {
			fmt.Println(v.Name)
			fmt.Println(v.Value)
		} else if v.Name == "idsrv" {
			fmt.Println(v.Name)
			fmt.Println(v.Value)
		}
	}
	fmt.Println(task.Cookies)
	return nil
}

type getBagResp struct {
	Bag struct {
		ID string `json:"id"`
	} `json:"bag"`
}

func (task *AsosTask) getBag() error {
	headers := tls.Headers{
		"Host":               "www.asos.com",
		"sec-ch-ua":          "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\"",
		"asos-c-istablet":    "false",
		"asos-c-plat":        "Web",
		"sec-ch-ua-mobile":   "?0",
		"authorization":      "Bearer " + task.bearerToken,
		"asos-c-name":        "Asos.Commerce.Bag.Sdk",
		"content-type":       "application/json",
		"accept":             "application/json, text/javascript, */*; q=0.01",
		"asos-c-ismobile":    "false",
		"asos-c-store":       "ROE",
		"asos-c-ver":         "5.5.195",
		"x-requested-with":   "XMLHttpRequest",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		"asos-cid":           uuid.NewString(),
		"sec-ch-ua-platform": "\"Windows\"",
		"origin":             "https://www.asos.com",
		"sec-fetch-site":     "same-origin",
		"sec-fetch-mode":     "cors",
		"sec-fetch-dest":     "empty",
		"referer":            "https://www.asos.com/nike/nike-dunk-high-retro-trainers-in-white-sail-and-black/prd/202278985",
		"accept-language":    "en-US,en;q=0.9",
	}

	postData := `{"currency":"EUR","lang":"en-GB","sizeSchema":"EU","country":"CZ","originCountry":"CZ"}`
	var customerId string
	url, _ := url.Parse("https://my.asos.com/identity")
	for _, v := range task.Cookies.Cookies(url) {
		if v.Name == "asos-anon12" {
			customerId = v.Value
		}
	}
	opts := tls.Options{
		URL:    "https://www.asos.com/api/commerce/bag/v4/customers/" + customerId + "/bags/getbag?expand=summary,total&lang=en-GB&keyStoreDataversion=dup0qtf-35",
		Method: "POST",

		Headers: headers,

		Body: postData,

		ClientSettings: &task.ClientSettings,
		Jar:            task.Cookies,
	}

	resp, err := tls.Do(opts)
	if err != nil {
		return err
	}
	if !(resp.StatusCode == 201 || resp.StatusCode == 200) {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}

	var getBagResp getBagResp
	if err := json.Unmarshal([]byte(resp.Body), &getBagResp); err != nil {
		return err
	}
	task.bagId = getBagResp.Bag.ID

	return nil
}
