package zalando

import (
	"net/url"
	"strings"

	"github.com/ProjektPLTS/pkg/customErrors"
	tls "github.com/ProjektPLTS/pkg/tls"
)

func (task *ZalTask) submitDiscountCode() error {
	var xXsrf string
	url, _ := url.Parse("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/api/checkout/redeem")
	for _, c := range task.Cookies.Cookies(url) {
		if c.Name == "frsx" {
			xXsrf = c.Value
		}
	}
	headers := tls.Headers{
		"Host":                   `www.zalando.` + strings.ToLower(task.Address.Delivery.Country),
		"x-xsrf-token":           xXsrf,
		"sec-ch-ua":              task.Sec_ch_ua,
		"sec-ch-ua-mobile":       "?0",
		"x-zalando-header-mode":  "desktop",
		"x-zalando-checkout-app": "web",
		"content-type":           "application/json",
		"x-zalando-footer-mode":  "desktop",
		"accept":                 "application/json",
		"user-agent":             task.UserAgent,
		"x-checkout-type-uuid":   "cdadbc8f-1754-5ad3-852a-d806ed07fe47",
		"sec-ch-ua-platform":     task.Sec_ch_ua_platform,
		"origin":                 "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country),
		"sec-fetch-site":         "same-origin",
		"sec-fetch-mode":         "cors",
		"sec-fetch-dest":         "empty",
		"referer":                "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/checkout/confirm",
		"accept-language":        "en-US,en;q=0.9",
	}
	jsonBody := `{"code":"` + task.Discount + `","pageRenderFlowId":""}`

	opts := tls.Options{
		URL:    "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/api/checkout/redeem",
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
