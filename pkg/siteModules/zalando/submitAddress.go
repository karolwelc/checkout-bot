package zalando

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/ProjektPLTS/pkg/customErrors"
	"github.com/ProjektPLTS/pkg/task"
	taskpkg "github.com/ProjektPLTS/pkg/task"
	tls "github.com/ProjektPLTS/pkg/tls"
)

type getProfile struct {
	Data struct {
		Customer struct {
			Addresses struct {
				Nodes []struct {
					ID                       string `json:"id"`
					City                     string `json:"city"`
					IsDefaultDeliveryAddress bool   `json:"isDefaultDeliveryAddress"`
					PostalCode               string `json:"postalCode"`
					Country                  struct {
						Typename string `json:"__typename"`
						Name     string `json:"name"`
						Code     string `json:"code"`
					} `json:"country"`
					Name struct {
						First      string `json:"first"`
						Last       string `json:"last"`
						Salutation string `json:"salutation"`
					} `json:"name"`
					IsDefaultBillingAddress bool   `json:"isDefaultBillingAddress"`
					Street                  string `json:"street"`
					Additional              string `json:"additional"`
				} `json:"nodes"`
			} `json:"addresses"`
		} `json:"customer"`
	} `json:"data"`
}

func (task *ZalTask) validateAddressWeb(address task.AddressType) (string, error) {
	var xXsrf string
	url, _ := url.Parse("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/api/checkout/validate-address")
	for _, c := range task.Cookies.Cookies(url) {
		if c.Name == "frsx" {
			xXsrf = c.Value
		}
	}
	headers := tls.Headers{
		"sec-ch-ua":              task.Sec_ch_ua,
		"x-xsrf-token":           xXsrf,
		"sec-ch-ua-mobile":       `?0`,
		"x-zalando-header-mode":  `desktop`,
		"x-zalando-checkout-app": `web`,
		"content-type":           `application/json`,
		"x-zalando-footer-mode":  `desktop`,
		"accept":                 `application/json`,
		"user-agent":             task.UserAgent,
		"x-checkout-type-uuid":   `cdadbc8f-1754-5ad3-852a-d806ed07fe47`,
		"sec-ch-ua-platform":     task.Sec_ch_ua_platform,
		"origin":                 `https://www.zalando.` + strings.ToLower(task.Address.Delivery.Country),
		"sec-fetch-site":         `same-origin`,
		"sec-fetch-mode":         `cors`,
		"sec-fetch-dest":         `empty`,
		"referer":                `https://www.zalando.` + strings.ToLower(task.Address.Delivery.Country) + `/checkout/address`,
		"accept-encoding":        `gzip, deflate, br`,
		"accept-language":        `en-US,en;q=0.9`,
	}
	jsonBody := `{"address":{"address":{"street":"` + address.Street1 + `","city":"` + address.City + `","zip":"` + address.ZipCode + `","id":"29489806","country_code":"` + strings.ToUpper(address.Country) + `","last_name":"` + address.LastName + `","first_name":"` + address.FirstName + `","salutation":"Ms"}}}`

	opts := tls.Options{
		URL:    "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/api/checkout/validate-address",
		Method: "POST",

		Headers: headers,
		Body:    jsonBody,

		ClientSettings: &task.ClientSettings,
		Jar:            task.Cookies,
	}
	resp, err := tls.Do(opts)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}
	task.Cookies = resp.Request.Jar
	return resp.Body, nil
}

func (task *ZalTask) submitAddressWeb(address task.AddressType, defaultShipping bool, defaultBilling bool) error {
	validatedAddy, err := task.validateAddressWeb(address)
	if err != nil {
		return err
	}
	var xXsrf string
	url, _ := url.Parse("https://www.zalando." + strings.ToLower(address.Country) + "/api/checkout/create-or-update-address")
	for _, c := range task.Cookies.Cookies(url) {
		if c.Name == "frsx" {
			xXsrf = c.Value
		}
	}
	headers := tls.Headers{
		"sec-ch-ua":              task.Sec_ch_ua,
		"x-xsrf-token":           xXsrf,
		"sec-ch-ua-mobile":       `?0`,
		"x-zalando-header-mode":  `desktop`,
		"x-zalando-checkout-app": `web`,
		"content-type":           `application/json`,
		"x-zalando-footer-mode":  `desktop`,
		"accept":                 `application/json`,
		"user-agent":             task.UserAgent,
		"x-checkout-type-uuid":   `cdadbc8f-1754-5ad3-852a-d806ed07fe47`,
		"sec-ch-ua-platform":     task.Sec_ch_ua_platform,
		"origin":                 `https://www.zalando.cz`,
		"sec-fetch-site":         `same-origin`,
		"sec-fetch-mode":         `cors`,
		"sec-fetch-dest":         `empty`,
		"referer":                `https://www.zalando.` + strings.ToLower(address.Country) + `/checkout/address`,
		"accept-encoding":        `gzip, deflate, br`,
		"accept-language":        `en-US,en;q=0.9`,
	}
	jsonBody := `{"address":{"street":"` + address.Street1 + `","city":"` + address.City + `","zip":"` + address.ZipCode + `","id":"` + address.ID + `","country_code":"` + address.Country + `","last_name":"` + address.LastName + `","first_name":"` + address.FirstName + `","salutation":"Mr"},"addressDestination":` + validatedAddy + `,"isDefaultShipping":` + strconv.FormatBool(defaultShipping) + `,"isDefaultBilling":` + strconv.FormatBool(defaultBilling) + `}`

	opts := tls.Options{
		URL:    "https://www.zalando." + strings.ToLower(address.Country) + "/api/checkout/create-or-update-address",
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

func (task *ZalTask) defaultAddressExists() error {
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
	jsonBody := `{"extensions":{"persistedQuery":{"sha256Hash":"a0c635ee7463c16181bfceaf6459befae83418d71bf017edfb447d70d2978197","version":1}},"id":"a0c635ee7463c16181bfceaf6459befae83418d71bf017edfb447d70d2978197","operationName":"GetUserProfile","variables":null}`

	opts := tls.Options{
		URL:    "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/api/graphql/mobile",
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
	var getProfile getProfile
	if err := json.Unmarshal([]byte(resp.Body), &getProfile); err != nil {
		return err
	}

	for _, address := range getProfile.Data.Customer.Addresses.Nodes {
		if address.IsDefaultDeliveryAddress {
			task.Address.Delivery = taskpkg.AddressType{
				ID:        address.ID,
				City:      address.City,
				ZipCode:   address.PostalCode,
				FirstName: address.Name.First,
				LastName:  address.Name.Last,
				Country:   address.Country.Code,
				Street1:   address.Street,
				Street2:   address.Additional,
			}
		}
		if address.IsDefaultBillingAddress {
			task.Address.Billing = taskpkg.AddressType{
				ID:        address.ID,
				City:      address.City,
				ZipCode:   address.PostalCode,
				FirstName: address.Name.First,
				LastName:  address.Name.Last,
				Country:   address.Country.Code,
				Street1:   address.Street,
				Street2:   address.Additional,
			}
		}
	}
	task.Cookies = resp.Request.Jar

	return nil
}
