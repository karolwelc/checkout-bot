package zalando

import (
	"encoding/json"
	"html"
	"strings"

	"github.com/ProjektPLTS/pkg/customErrors"
	tls "github.com/ProjektPLTS/pkg/tls"
)

func (task *ZalTask) submitPaymentWeb() error {
	jsonResp, err := task.getPaymentWeb()
	if err != nil {
		return err
	}
	headers := tls.Headers{
		"Host":            `purchase-session.client-api.payment.zalando.com`,
		"Content-Type":    `application/json`,
		"Origin":          `https://www.zalando.` + strings.ToLower(task.Address.Delivery.Country),
		"Content-Length":  `45`,
		"Accept-Encoding": `gzip, deflate, br`,
		"Connection":      `keep-alive`,
		"Accept":          `*/*`,
		"User-Agent":      `Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148`,
		"Referer":         `https://www.zalando.cz/`,
		"Authorization":   `Bearer ` + jsonResp.Model.PurchaseSessionToken,
		"Accept-Language": `cs-CZ`,
	}
	jsonBody := `{"options":{"selected":[],"not_selected":[]}}`

	opts := tls.Options{
		URL:    jsonResp.Model.PurchaseSessionURL + "/checkout/payment?payment_method_id=paypal",
		Method: "POST",

		Headers: headers,
		Body:    jsonBody,

		ClientSettings: &task.ClientSettings,
	}
	resp, err := tls.Do(opts)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}
	return nil
}

type paymentResp struct {
	Model struct {
		PurchaseSessionURL   string `json:"purchaseSessionUrl"`
		PurchaseSessionToken string `json:"purchaseSessionToken"`
	} `json:"model"`
}

func (task *ZalTask) getPaymentWeb() (paymentResp, error) {
	headers := tls.Headers{
		"Host":                       `www.zalando.` + strings.ToLower(task.Address.Delivery.Country),
		"Accept":                     `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`,
		"X-Device-OS":                `ios`,
		"X-Zalando-Footer-Mode":      `mobile`,
		"x-zalando-checkout-app":     `app`,
		"Accept-Language":            `cs-CZ,cs;q=0.9`,
		"Accept-Encoding":            `gzip, deflate, br`,
		"x-device-platform":          `ios`,
		"x-zalando-checkout-uuid":    `3361EFE1-A3A8-4258-933B-0D8A8EE1A37A`,
		"Referer":                    `https://www.zalando.cz/checkout/address?clientId=&appName=zalando&appId=de.zalando.iphone&appVersion=22.11.1&appCountry=CZ&uId=3969c50c247fc3a3b78e39a17a152e04`,
		"x-zalando-checkout-webview": `WKWebView`,
		"x-app-version":              `22.11.1`,
		"User-Agent":                 `Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148`,
		"Connection":                 `keep-alive`,
		"X-Zalando-Header-Mode":      `mobile`,
	}

	opts := tls.Options{
		URL:    "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/checkout/payment",
		Method: "GET",

		Headers: headers,

		ClientSettings: &task.ClientSettings,
		Jar:            task.Cookies,
	}
	resp, err := tls.Do(opts)
	if err != nil {
		return paymentResp{}, err
	}
	if resp.StatusCode != 200 {
		return paymentResp{}, customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}

	jsonStringDecoded := strings.Split(strings.Split(resp.Body, `</script><div data-props='`)[1], `' data-translations`)[0]
	jsonDecoded := (html.UnescapeString(jsonStringDecoded))
	var jsonResp paymentResp

	if err := json.Unmarshal([]byte(jsonDecoded), &jsonResp); err != nil {
		return jsonResp, err
	}
	task.Cookies = resp.Request.Jar
	return jsonResp, nil
}
