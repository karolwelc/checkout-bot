package zalando

import (
	"encoding/json"
	"html"
	"net/url"
	"strings"

	"github.com/ProjektPLTS/pkg/customErrors"
	tls "github.com/ProjektPLTS/pkg/tls"
)

type nextStepResp struct {
	URL string `json:"url"`
}

func (task *ZalTask) nextStepApp() (nextStepResp, error) {
	var xXsrf string
	url, _ := url.Parse("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/api/checkout/next-step")
	for _, c := range task.Cookies.Cookies(url) {
		if c.Name == "frsx" {
			xXsrf = c.Value
		}
	}
	headers := tls.Headers{
		"Host":                    `www.zalando.` + strings.ToLower(task.Address.Delivery.Country),
		"Accept":                  `application/json`,
		"x-device-os":             `ios`,
		"x-zalando-footer-mode":   `mobile`,
		"x-zalando-checkout-app":  `app`,
		"x-xsrf-token":            xXsrf,
		"Accept-Language":         `cs-CZ,cs;q=0.9`,
		"Accept-Encoding":         `gzip, deflate, br`,
		"x-zalando-checkout-uuid": `3361EFE1-A3A8-4258-933B-0D8A8EE1A37A`,
		"Content-Type":            `application/json`,
		"User-Agent":              `Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148`,
		"Referer":                 `https://www.zalando.cz/checkout/address?clientId=&appName=zalando&appId=de.zalando.iphone&appVersion=22.11.1&appCountry=CZ&uId=dde84987fd9d511746e5b7f18efcca5b`,
		"x-app-version":           `22.11.1`,
		"Connection":              `keep-alive`,
		"x-checkout-type":         `web`,
		"x-zalando-header-mode":   `mobile`,
	}

	opts := tls.Options{
		URL:    "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/api/checkout/next-step",
		Method: "GET",

		Headers: headers,

		ClientSettings: &task.ClientSettings,
		Jar:            task.Cookies,
	}
	resp, err := tls.Do(opts)
	if err != nil {
		return nextStepResp{}, err
	}
	if resp.StatusCode != 200 {
		return nextStepResp{}, customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}
	var jsonResp nextStepResp
	if err := json.Unmarshal([]byte(resp.Body), &jsonResp); err != nil {
		return jsonResp, err
	}

	return jsonResp, nil
}

type checkoutData struct {
	Model struct {
		CheckoutID string `json:"checkoutId"`
		ETag       string `json:"eTag"`
	} `json:"model"`
}

func (task *ZalTask) getCheckoutWeb() error {
	headers := tls.Headers{
		"authority":                 "www.zalando." + strings.ToLower(task.Address.Delivery.Country),
		"accept":                    `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`,
		"referer":                   "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/checkout/payment",
		"sec-ch-ua":                 task.Sec_ch_ua,
		"sec-ch-ua-mobile":          `?0`,
		"sec-ch-ua-platform":        task.Sec_ch_ua_platform,
		"sec-fetch-dest":            `document`,
		"sec-fetch-mode":            `navigate`,
		"sec-fetch-site":            `same-origin`,
		"sec-fetch-user":            `?1`,
		"upgrade-insecure-requests": `1`,
		"user-agent":                task.UserAgent,
	}

	opts := tls.Options{
		Method: "GET",
		URL:    "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/checkout/confirm",

		Headers: headers,

		ClientSettings: &task.ClientSettings,
		Jar:            task.Cookies,
	}
	resp, err := tls.Do(opts)
	if err != nil {
		return err
	}

	//Known status codes on checkout get
	// 302 - session expired - redirect to login page / (TODO) handle it by going back to login if status code == 302

	if resp.StatusCode != 200 {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}

	jsonEncoded := strings.Split(strings.Split(resp.Body, `data-props='`)[1], `' data-translations`)[0]

	jsonDecoded := (html.UnescapeString(jsonEncoded))

	var data checkoutData
	err = json.Unmarshal([]byte(jsonDecoded), &data)
	if err != nil {
		return err
	}
	task.checkoutId = data.Model.CheckoutID
	task.eTag = data.Model.ETag
	if task.checkoutId == "" || task.eTag == "" {
		return customErrors.EmptyField
	}
	task.Cookies = resp.Request.Jar
	return nil
}

func (task *ZalTask) postCheckoutWeb() error {
	var xXsrf string
	url, _ := url.Parse("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/api/checkout/buy-now")
	for _, c := range task.Cookies.Cookies(url) {
		if c.Name == "frsx" {
			xXsrf = c.Value
		}
	}
	headers := tls.Headers{
		"authority":              "www.zalando." + strings.ToLower(task.Address.Delivery.Country),
		"accept":                 `application/json`,
		"accept-language":        `pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7'`,
		"cache-control":          `no-cache`,
		"dnt":                    `1`,
		"origin":                 "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country),
		"pragma":                 `no-cache`,
		"referer":                "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/checkout/confirm",
		"sec-ch-ua":              task.Sec_ch_ua,
		"sec-ch-ua-mobile":       `?0`,
		"sec-ch-ua-platform":     task.Sec_ch_ua_platform,
		"sec-fetch-dest":         `empty`,
		"sec-fetch-mode":         `cors`,
		"sec-fetch-site":         `same-origin`,
		"user-agent":             task.UserAgent,
		"x-checkout-type":        `web`,
		"x-xsrf-token":           xXsrf,
		"content-type":           "application/json",
		"x-zalando-checkout-app": `web`,
		"x-zalando-footer-mode":  `desktop`,
		"x-zalando-header-moden": `desktop`,
	}
	body := `{"checkoutId": "` + task.checkoutId + `","eTag": ` + task.eTag + `}`
	opts := tls.Options{
		Method: "POST",
		URL:    "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/api/checkout/buy-now",

		Headers: headers,
		Body:    body,

		ClientSettings: &task.ClientSettings,
		Jar:            task.Cookies,
	}
	resp, err := tls.Do(opts)
	if err != nil {
		return err
	}

	//Known status codes on checkout post
	// 302 - session expired - redirect to login page
	// 401 - Unauthorized - x-xsrf-token is invalid
	// 400 - bad request - checkoutId is invalid
	// 403 - banned proxy(i think, 99%) !!!!!!!!!!!!!!!!!!!!!!!!!!!! (TODO) handle all of them respectively

	if resp.StatusCode != 200 {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}

	return nil
}

func (task *ZalTask) checkout() error {
	err := task.getCheckoutWeb()
	if err != nil {
		return err
	}
	err = task.postCheckoutWeb()
	if err != nil {
		return err
	}
	return nil
}
