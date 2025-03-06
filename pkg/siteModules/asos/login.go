package asos

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ProjektPLTS/pkg/customErrors"
	securedtouch "github.com/ProjektPLTS/pkg/securedTouch"
	tls "github.com/ProjektPLTS/pkg/tls"
	"github.com/google/uuid"
)

func (task *AsosTask) postLogin() error {
	headers := tls.Headers{
		"Host":                      "my.asos.com",
		"cache-control":             "max-age=0",
		"sec-ch-ua":                 task.Sec_ch_ua,
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        task.Sec_ch_ua_platform,
		"upgrade-insecure-requests": "1",
		"origin":                    "https://my.asos.com",
		"content-type":              "application/x-www-form-urlencoded",
		"user-agent":                task.UserAgent,
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"sec-fetch-site":            "same-origin",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-user":            "?1",
		"sec-fetch-dest":            "document",
		"referer":                   task.loginUrl,
		"accept-language":           "en-US,en;q=0.9",
	}

	postData := `idsrv.xsrf=` + task.xsrf + `&SecuredTouchToken=` + task.securedTouch + `&Username=` + task.Email + `&Password=` + task.Password

	opts := tls.Options{
		URL:    task.loginUrl + "&checkout=False&showAllOptions=False",
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
	if resp.StatusCode != 200 {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}
	task.Cookies = resp.Request.Jar
	return nil
}

func (task *AsosTask) getLoginpage() error {
	headers := tls.Headers{
		"Host":                      "my.asos.com",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                task.UserAgent,
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Accept-Language":           "en-US,en;q=0.9",
	}

	opts := tls.Options{
		URL:    "http://my.asos.com/identity/login",
		Method: "GET",

		Headers:        headers,
		ClientSettings: &task.ClientSettings,

		FollowRedirects: true,
		Jar:             task.Cookies,
	}

	resp, err := tls.Do(opts)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}
	task.Cookies = resp.Request.Jar
	task.loginUrl = resp.Request.URL
	task.appSessId = strings.Split(strings.Split(resp.Body, `sessionId: '`)[1], "'")[0]
	task.xsrf = strings.Split(strings.Split(resp.Body, `name="idsrv.xsrf" type="hidden" value="`)[1], `" /><input`)[0]
	task.akamaiEndpoint = strings.Split(strings.Split(strings.Split(resp.Body, `<script type="text/javascript" nonce="`)[1], `" src="`)[1], `"></script></body>`)[0]
	return nil
}

func (task *AsosTask) getAkamaiScript() error {
	headers := tls.Headers{
		"Host":               "my.asos.com",
		"sec-ch-ua":          "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\"",
		"sec-ch-ua-mobile":   "?0",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		"sec-ch-ua-platform": "\"Windows\"",
		"accept":             "*/*",
		"sec-fetch-site":     "same-origin",
		"sec-fetch-mode":     "no-cors",
		"sec-fetch-dest":     "script",
		"referer":            task.loginUrl,
		"accept-language":    "en-US,en;q=0.9",
	}

	opts := tls.Options{
		URL:    "https://my.asos.com" + task.akamaiEndpoint,
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
	return nil
}

func (task *AsosTask) getSecuredTouch() error {
	task.instanceuuid = uuid.New().String()
	err := task.securedTouchStarter()
	if err != nil {
		return nil
	}
	inputData := securedtouch.InputData{
		AppSessId:    task.appSessId,
		LocationHref: task.loginUrl,
		StToken:      task.securedTouch,
		CheckSum:     task.checkSum,
		DeviceId:     task.deviceID,
	}
	err = task.securedTouchInteractions(inputData)
	if err != nil {
		return nil
	}
	err = task.securedTouchPong()
	if err != nil {
		return nil
	}

	return nil
}

type startResp struct {
	Token    string `json:"token"`
	Checksum string `json:"checksum"`
	DeviceID string `json:"deviceId"`
}

func (task *AsosTask) securedTouchStarter() error {

	headers := tls.Headers{
		"Host":               "st.asos.com",
		"sec-ch-ua":          "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\"",
		"sec-ch-ua-mobile":   "?0",
		"authorization":      "YjIxMzVjdDIxSnVsVnlP",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		"clientepoch":        strconv.Itoa(int(time.Now().Unix())),
		"content-type":       "application/json",
		"accept":             "application/json",
		"attempt":            "0",
		"instanceuuid":       task.instanceuuid,
		"clientversion":      "3.13.2w",
		"sec-ch-ua-platform": "\"Windows\"",
		"origin":             "https://my.asos.com",
		"sec-fetch-site":     "same-site",
		"sec-fetch-mode":     "cors",
		"sec-fetch-dest":     "empty",
		"referer":            "https://my.asos.com/",
		"accept-language":    "en-US,en;q=0.9",
	}

	postData := `{"device_id":"Id-` + task.instanceuuid + `","clientVersion":"3.13.2w","deviceType":"Chrome(106.0.0.0)-Windows(10)","authToken":""}`

	opts := tls.Options{
		URL:    "https://st.asos.com/SecuredTouch/rest/services/v2/starter/asos",
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

	if resp.StatusCode != 200 {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}

	var startResp startResp
	if err := json.Unmarshal([]byte(resp.Body), &startResp); err != nil {
		return err
	}
	task.securedTouch = startResp.Token
	task.checkSum = startResp.Checksum
	task.deviceID = startResp.DeviceID

	task.Cookies = resp.Request.Jar
	return nil
}

func (task *AsosTask) securedTouchInteractions(data securedtouch.InputData) error {

	//motionData := securedtouch.GenSecuredTouchData(data)
	motionData := securedtouch.GenSecuredTouchData(data)
	os.Setenv("HTTPS_PROXY", "http://localhost:8888")
	req, err := http.NewRequest("POST", "https://st.asos.com/SecuredTouch/rest/services/v2/interactions/asos", strings.NewReader(string(motionData)))
	if err != nil {
		return err
	}
	/* url, _ := url.Parse("https://st.asos.com")
	cookies := task.Cookies.Cookies(url)
	var _abck, bm_sz, ak_bmsc, st_deviceID string
	for _, c := range cookies {
		if c.Name == "_abck" {
			_abck = c.Value
		} else if c.Name == "bm_sz" {
			bm_sz = c.Value
		} else if c.Name == "ak_bmsc" {
			ak_bmsc = c.Value
		} else if c.Name == "st_deviceId" {
			st_deviceID = c.Value
		}
	} */

	req.Header = http.Header{
		"Host":               {`st.asos.com`},
		"User-Agent":         {`Mozilla/5.0 (Windows NT 10.0; Win64; x64) Applewebkit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36`},
		"accept-encoding":    {`gzip, deflate, br`},
		"accept":             {`application/json`},
		"Connection":         {`keep-alive`},
		"accept-language":    {`de-DE, de;q=0.9`},
		"attempt":            {`0`},
		"authorization":      {task.securedTouch},
		"clientepoch":        {`1665431713140`},
		"clientversion":      {`3.13.2w`},
		"content-encoding":   {`gzip`},
		"content-type":       {`application/json`},
		"encrypted":          {`1`},
		"instanceuuid":       {task.instanceuuid},
		"origin":             {`https://my.asos.com`},
		"referer":            {`https://my.asos.com/`},
		"sec-ch-ua":          {`"" Not A; Brand"; v="99", "Chromium"; v="102", "Google Chrome"; v="102"`},
		"sec-ch-ua-platform": {`""Windows"`},
		"sec-fetch-dest":     {`empty`},
		"sec-fetch-mode":     {`cors`},
		"sec-fetch-site":     {`same-site`},
		"Cookie":             {`_abck=B31B03093F18C068F85271B0D99A6F0D~-1~YAAQDk4SAjgLK76DAQAA2TJ3wwhGBNhQKX5IhbQL2FRzzeWjkivZzOpwedA/c5qjLnd1fNOsmcwJP+hdqKznJll497atl8oQxAxvvV+hCoXxvbEaVzi+MdMxUIdlXmDoK+BWn5BqNkE+MwEmfSyrdyxqKqQXh1OzJFJGfIFyIa8CBZDuWhkuo4p+0Zx508Jl6bhLOVqFLeCixMynODzSGh7rmnBqS0MVWR53G1Eb8ik+mXxcmQjkxN+OuPUWWxHQ1GKZQBp2Mv/SLTz7qoEjgvkXFZV9zk5fvDkggHwgbNtQhHNIay+wWqc8Ddvh23IGcRDeVDzSAsSJ97tlW48yM7lTWpPXXroYNfexwxF+iWW5n5qVpPnle/qu2i9RKWiSfkq16Tc=~-1~-1~1665435213; ak_bmsc=A5752BA2FDA0DB75E0707F3878D2F2F2~000000000000000000000000000000~YAAQDk4SAjwLK76DAQAAOTN3wxG1xroKGBmUtblgupXvCU1TL4+0+6WHb4+MmM76uLiPZQM6aDwEVSV2hjSDIb4ZsVVt3m3ItfEJFuX+cgIWCYn0sppH2m5MfoYYg0pNPzehJ9OA4zNGqtMqVIAWZ09xgWDB5N7Rg1EeUbOOyl+1AGnYX1i7X+RwtPDdVNM7RuihCOWOB1rR1QXDA/Gjjqq4kUd4raTPZtkVbVwF3EzjOWII/NOYZs8UqIFZWz8uRNsQf2zSBc2LwNWthosVDR79g0sbUpderhJRpXgKINjFIPulz/HpPRW73RfM7bXU1I9lld+Hc2Sbuo9JGPe+VVpBT1wPe5F+cGd4MrlSFATb7BYqXVVZUchTKA==; bm_sz=B397047E5D56B734E69315E8AEB94126~YAAQDk4SAjALK76DAQAAAzJ3wxHcFXrsz/9Ld8ui43voiMZpPBCgm2BBcG1QjvNATFHuHDTnl2Ljm1oLvQGMTBiPby/xh9P28s6Mh+GOze2YXRWMHdJbNnLMru2pm+vLwxuB/2ZsywTMzgyw4kbDAHtgvpf7HOY54OFQrOvmWXLdJpsq3KNR5u2UwDelECFU015CL2KEvFOqu9ouDXA/E8E6uETkAZqO3SUsiGiF2zWfXdZSzH+IVmi/DEewB2pK2UVcQRPmfCuriS2RQARJVQRvD6bWrdZ962RzAJOxZGFL~3225908~3687990; st_deviceId=ZjVmNjlmZDMtN2NlNy00ZmZmLWExZWUtNWVmNDQ0ZjdjMGQz`},
		"Content-Length":     {`3779`},
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	/* headers := tls.Headers{
		"Host":               "st.asos.com",
		"sec-ch-ua":          "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\"",
		"sec-ch-ua-mobile":   "?0",
		"authorization":      task.securedTouch,
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		"clientepoch":        "1665342901181",
		"content-type":       "application/json",
		"accept":             "application/json",
		"attempt":            "0",
		"instanceuuid":       task.instanceuuid,
		"encrypted":          "1",
		"clientversion":      "3.13.2w",
		"sec-ch-ua-platform": "\"Windows\"",
		"origin":             "https://my.asos.com",
		"sec-fetch-site":     "same-site",
		"sec-fetch-mode":     "cors",
		"sec-fetch-dest":     "empty",
		"referer":            "https://my.asos.com/",
		"accept-language":    "en-US,en;q=0.9",
		"content-length":     `1770`,
		"content-encoding":   `gzip`,
		"accept-encoding":    `gzip, deflate, br`,
	} */
	/* opts := tls.Options{
		URL:    "https://st.asos.com/SecuredTouch/rest/services/v2/interactions/asos",
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

	if resp.StatusCode != 200 {
		return customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}

	task.Cookies = resp.Request.Jar*/
	return nil
}

func (task *AsosTask) securedTouchPong() error {
	body := base64.StdEncoding.EncodeToString([]byte(`{"pingVersion":"1.3.0p","appId":"asos","appSessionId":"` + task.appSessId + `"}`))
	opts := tls.Options{
		URL:    "https://st-static.asos.com/sdk/pong.js?body=" + body,
		Method: "GET",

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
	task.securedTouch = strings.Split(strings.Split(resp.Body, `window["_securedTouchToken"] = '`)[1], "';")[0]

	return nil
}

func (task *AsosTask) login() error {
	err := task.getLoginpage()
	if err != nil {
		return err
	}
	err = task.getSecuredTouch()
	if err != nil {
		return err
	}
	err = task.getAkamaiScript()
	if err != nil {
		return err
	}
	for !isCookieValid(task.Cookies, "https://my.asos.com") {
		err = task.submitSensorData("https://my.asos.com")
		if err != nil {
			return err
		}
	}

	err = task.postLogin()
	if err != nil {
		return err
	}
	return nil
}
