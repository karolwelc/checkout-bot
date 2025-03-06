package asos

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/ProjektPLTS/pkg/akamai"
	"github.com/ProjektPLTS/pkg/customErrors"
	tls "github.com/ProjektPLTS/pkg/tls"
	"github.com/ProjektPLTS/pkg/tls/http/cookiejar"
)

func (task *AsosTask) submitSensorData(baseUrl string) error {
	var _abck, bm_sz string
	url, _ := url.Parse(baseUrl)
	cookies := task.Cookies.Cookies(url)
	for _, c := range cookies {
		if c.Name == "_abck" {
			_abck = c.Value
		} else if c.Name == "bm_sz" {
			bm_sz = c.Value
		}
	}

	sensorData, err := akamai.GenerateSensorData(baseUrl, _abck, bm_sz, task.configSC)
	if err != nil {
		return err
	}

	jsonBody := `{"sensor_data": "` + sensorData + `"}`
	headers := tls.Headers{
		"content-length":     strconv.Itoa(len(jsonBody)), //add
		"sec-ch-ua":          task.Sec_ch_ua,
		"sec-ch-ua-mobile":   `?0`,
		"user-agent":         task.UserAgent,
		"sec-ch-ua-platform": task.Sec_ch_ua_platform,
		"content-type":       `text/plain;charset=UTF-8`,
		"accept":             `*/*`,
		"sec-fetch-site":     `same-origin`,
		"sec-fetch-mode":     `cors`,
		"sec-fetch-dest":     `empty`,
		"accept-encoding":    `gzip, deflate, br`,
		"accept-language":    `en-US,en;q=0.9`,
	}
	headerOrder := tls.HeaderOrder{
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"user-agent",
		"sec-ch-ua-platform",
		"content-type",
		"accept",
		"origin",
		"sec-fetch-site",
		"sec-fetch-mode",
		"sec-fetch-dest",
		"referer",
		"accept-encoding",
		"accept-language",
		"cookie",
		"content-length",
	}

	reqOptions := tls.Options{
		URL:    baseUrl + task.akamaiEndpoint,
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
	task.Cookies = resp.Request.Jar

	return nil
}

func isCookieValid(cookieJar *cookiejar.Jar, cookieUrl string) bool {
	var _abck string
	url, _ := url.Parse(cookieUrl)
	cookies := cookieJar.Cookies(url)

	for _, c := range cookies {
		if c.Name == "_abck" {
			_abck = c.Value
		}
	}
	return strings.Split(_abck, "~")[1] == "0"
}
