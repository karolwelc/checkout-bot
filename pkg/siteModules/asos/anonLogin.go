package asos

import (
	"encoding/json"
	"strings"

	"github.com/ProjektPLTS/pkg/customErrors"
	tls "github.com/ProjektPLTS/pkg/tls"
	nonce "github.com/larrybattle/nonce-golang"
)

type anonResp struct {
	AccessToken string `json:"access_token"`
}

func (task *AsosTask) anonLogin() error {
	headers := tls.Headers{
		"Host":               "my.asos.com",
		"sec-ch-ua":          "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\"",
		"accept":             "application/json",
		"sec-ch-ua-mobile":   "?0",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		"sec-ch-ua-platform": "\"Windows\"",
		"origin":             "https://www.asos.com",
		"sec-fetch-site":     "same-site",
		"sec-fetch-mode":     "cors",
		"sec-fetch-dest":     "empty",
		"referer":            "https://www.asos.com/",
		"accept-language":    "en-US,en;q=0.9",
	}

	opts := tls.Options{
		URL:    "https://my.asos.com/identity/connect/authorize?nonce=" + nonce.NewToken() + "&client_id=D91F2DAA-898C-4E10-9102-D6C974AFBD59&redirect_uri=https://www.asos.com&response_type=id_token%20token&scope=openid%20sensitive%20profile&ui_locales=en-GB&acr_values=0&response_mode=json&store=ROE&country=CZ&keyStoreDataversion=dup0qtf-35&lang=en-GB",
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

	task.Cookies = resp.Request.Jar
	var anonResp anonResp
	if err := json.Unmarshal([]byte(resp.Body), &anonResp); err != nil {
		return err
	}
	task.bearerToken = anonResp.AccessToken
	return nil
}

func (task *AsosTask) getAkamaiScriptAnon() error {
	headers := tls.Headers{
		"Host":               "www.asos.com",
		"sec-ch-ua":          "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\"",
		"sec-ch-ua-mobile":   "?0",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		"sec-ch-ua-platform": "\"Windows\"",
		"accept":             "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8",
		"sec-fetch-site":     "same-origin",
		"sec-fetch-mode":     "no-cors",
		"sec-fetch-dest":     "image",
		"referer":            "https://www.asos.com/men/",
		"accept-language":    "en-US,en;q=0.9",
	}

	opts := tls.Options{
		URL:    "https://www.asos.com/men/",
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
	task.Cookies = resp.Request.Jar
	task.configSC = resp.Body
	task.akamaiEndpoint = strings.Split(strings.Split(resp.Body, ` <script type="text/javascript"  src="`)[1], `"></script></body>`)[0]
	return nil
}
