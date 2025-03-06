package zalando

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/ProjektPLTS/pkg/customErrors"

	tls "github.com/ProjektPLTS/pkg/tls"
)

func (task *ZalTask) getLoginPage() error {
	headers := tls.Headers{
		"sec-ch-ua":                 task.Sec_ch_ua,
		"sec-ch-ua-mobile":          `?0`,
		"sec-ch-ua-platform":        task.Sec_ch_ua_platform,
		"upgrade-insecure-requests": `1`,
		"user-agent":                task.UserAgent,
		"accept":                    `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`,
		"sec-fetch-site":            `none`,
		"sec-fetch-mode":            `navigate`,
		"sec-fetch-user":            `?1`,
		"sec-fetch-dest":            `document`,
		"accept-encoding":           `gzip, deflate, br`,
		"accept-language":           `en-US,en;q=0.9`,
	}
	headerOrder := tls.HeaderOrder{
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"sec-ch-ua-platform",
		"upgrade-insecure-requests",
		"user-agent",
		"accept",
		"sec-fetch-site",
		"sec-fetch-mode",
		"sec-fetch-user",
		"sec-fetch-dest",
		"accept-encoding",
		"accept-language",
	}
	reqOptions := tls.Options{
		URL:    "https://www.zalando." + strings.ToLower(task.Address.Delivery.Country) + "/login",
		Method: "GET",

		Headers:     headers,
		HeaderOrder: headerOrder,

		ClientSettings:  &task.ClientSettings,
		FollowRedirects: true,
	}
	resp, err := tls.Do(reqOptions)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}
	task.Cookies = resp.Request.Jar
	task.loginUrl = resp.Request.URL
	task.akamaiEndpoint = strings.Split(strings.Split(resp.Body, `<script type="text/javascript"  src="`)[1], `"></script>`)[0]
	return nil
}

func (task *ZalTask) getAkamaiScript() error {
	headers := tls.Headers{
		"sec-ch-ua":          task.Sec_ch_ua,
		"sec-ch-ua-mobile":   `?0`,
		"user-agent":         task.UserAgent,
		"sec-ch-ua-platform": task.Sec_ch_ua_platform,
		"accept":             `*/*`,
		"sec-fetch-site":     `same-origin`,
		"sec-fetch-mode":     `no-cors`,
		"sec-fetch-dest":     `script`,
		"referer":            task.loginUrl,
		"accept-encoding":    `gzip, deflate, br`,
		"accept-language":    `en-US,en;q=0.9`,
	}
	headerOrder := tls.HeaderOrder{
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"user-agent",
		"sec-ch-ua-platform",
		"accept",
		"sec-fetch-site",
		"sec-fetch-mode",
		"sec-fetch-dest",
		"referer",
		"accept-encoding",
		"accept-language",
		"cookie",
	}
	reqOptions := tls.Options{
		URL:    "https://accounts.zalando.com" + task.akamaiEndpoint,
		Method: "GET",

		Headers:     headers,
		HeaderOrder: headerOrder,

		ClientSettings: &task.ClientSettings,
		Jar:            task.Cookies,
	}
	resp, err := tls.Do(reqOptions)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}
	task.Cookies = resp.Request.Jar
	task.configSC = resp.Body
	return nil
}

type loginResp struct {
	Status      bool   `json:"status"`
	RedirectURL string `json:"redirect_url"`
}

func (task *ZalTask) postLogin() error {
	request := strings.Split(strings.Split(task.loginUrl, "&request=")[1], "&")[0]
	jsonBody := `{"email":"` + task.Email + `","secret":"` + task.Password + `","request":"` + request + `"}`
	headers := tls.Headers{
		"content-length":     strconv.Itoa(len(jsonBody)),
		"sec-ch-ua":          task.Sec_ch_ua,
		"dnt":                `1`,
		"x-csrf-token":       task.csrfToken,
		"sec-ch-ua-mobile":   `?0`,
		"user-agent":         task.UserAgent,
		"content-type":       `application/json`,
		"accept":             `application/json`,
		"sec-ch-ua-platform": task.Sec_ch_ua_platform,
		"origin":             `https://accounts.zalando.com`,
		"sec-fetch-site":     `same-origin`,
		"sec-fetch-mode":     `cors`,
		"sec-fetch-dest":     `empty`,
		"referer":            task.loginUrl,
		"accept-encoding":    `gzip, deflate, br`,
		"accept-language":    `en-GB,en-US;q=0.9,en;q=0.8`,
	}
	headerOrder := tls.HeaderOrder{
		"content-length",
		"sec-ch-ua",
		"dnt",
		"x-csrf-token",
		"sec-ch-ua-mobile",
		"user-agent",
		"content-type",
		"accept",
		"sec-ch-ua-platform",
		"origin",
		"sec-fetch-site",
		"sec-fetch-mode",
		"sec-fetch-dest",
		"referer",
		"accept-encoding",
		"accept-language",
		"cookie",
	}
	reqOptions := tls.Options{
		URL:    "https://accounts.zalando.com/api/login",
		Method: "POST",
		Body:   jsonBody,

		Headers:     headers,
		HeaderOrder: headerOrder,

		ClientSettings: &task.ClientSettings,
		Jar:            task.Cookies,
	}
	resp, err := tls.Do(reqOptions)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}
	var jsonResp loginResp
	if err := json.Unmarshal([]byte(resp.Body), &jsonResp); err != nil {
		return err
	}
	if !jsonResp.Status {
		return errors.New("unexpected error")
	}
	task.redirect_url = jsonResp.RedirectURL
	task.Cookies = resp.Request.Jar
	return nil
}

func (task *ZalTask) followRedirect() error {
	headers := tls.Headers{
		"sec-ch-ua":          task.Sec_ch_ua,
		"dnt":                `1`,
		"x-csrf-token":       task.csrfToken,
		"sec-ch-ua-mobile":   `?0`,
		"user-agent":         task.UserAgent,
		"content-type":       `application/json`,
		"accept":             `application/json`,
		"sec-ch-ua-platform": task.Sec_ch_ua_platform,
		"origin":             `https://accounts.zalando.com`,
		"sec-fetch-site":     `same-origin`,
		"sec-fetch-mode":     `cors`,
		"sec-fetch-dest":     `empty`,
		"referer":            task.loginUrl,
		"accept-encoding":    `gzip, deflate, br`,
		"accept-language":    `en-GB,en-US;q=0.9,en;q=0.8`,
	}
	headerOrder := tls.HeaderOrder{
		"content-length",
		"sec-ch-ua",
		"dnt",
		"x-csrf-token",
		"sec-ch-ua-mobile",
		"user-agent",
		"content-type",
		"accept",
		"sec-ch-ua-platform",
		"origin",
		"sec-fetch-site",
		"sec-fetch-mode",
		"sec-fetch-dest",
		"referer",
		"accept-encoding",
		"accept-language",
		"cookie",
	}
	reqOptions := tls.Options{
		URL:    "https://accounts.zalando.com" + task.redirect_url,
		Method: "GET",

		Headers:     headers,
		HeaderOrder: headerOrder,

		ClientSettings: &task.ClientSettings,
		Jar:            task.Cookies,

		FollowRedirects: true,
	}
	resp, err := tls.Do(reqOptions)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}

	task.Cookies = resp.Request.Jar
	task.akamaiEndpoint = strings.Split(strings.Split(resp.Body, `</noscript><script type="text/javascript"  src="`)[1], `"></script></body>`)[0]
	return nil
}

func (task *ZalTask) login() error {

	err := task.getLoginPage()
	if err != nil {
		return err
	}
	err = task.getAkamaiScript()
	if err != nil {
		return err
	}
	url, _ := url.Parse(task.loginUrl)
	cookies := task.Cookies.Cookies(url)
	for _, c := range cookies {
		if c.Name == "csrf-token" {
			task.csrfToken = c.Value
		}
	}

	// Submit sensor data until valid
	for !isCookieValid(task.Cookies, task.loginUrl) {
		err = task.submitSensorData("https://accounts.zalando.com")
		if err != nil {
			return err
		}
	}

	err = task.postLogin()
	if err != nil {
		return err
	}
	err = task.followRedirect()
	if err != nil {
		return err
	}
	return nil
}
