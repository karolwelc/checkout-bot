package akamai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ProjektPLTS/pkg/customErrors"
)

type hawkPayloadData struct {
	Site      string `json:"site"`
	Abck      string `json:"abck"`
	Events    string `json:"events"`
	UserAgent string `json:"user_agent"`
	Bm_sz     string `json:"bm_sz"`
	//Config    string `json:"config"`
}

const (
	apiKey = "7f66d849-6091-4481-b00e-c6a7d6f67947"
)

func GenerateSensorData(targetUrl string, _abck string, bm_sz string, scriptSource string) (string, error) {
	userAgent, err := getHawkUserAgent()
	if err != nil {
		return "", err
	}
	//config, err := getConfig(scriptSource)
	if err != nil {
		return "", err
	}
	payload := hawkPayloadData{
		Site:      targetUrl,
		Abck:      _abck,
		Events:    "1,1",
		UserAgent: userAgent,
		//Config:    config,
	}

	if bm_sz != "" {
		payload.Bm_sz = bm_sz
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://ak01-eu.hwkapi.com/akamai/generate", bytes.NewBuffer(payloadJson))
	if err != nil {
		return "", err
	}
	req.Header = http.Header{
		"content-type": {"application/json"},
		"X-Api-Key":    {apiKey},
		"X-Sec":        {"new"},
	}
	url, _ := url.Parse("http://localhost:8888")
	client := http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(url)},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.Split(string(respBody), "****")[0], nil
}

func getHawkUserAgent() (string, error) {
	req, err := http.NewRequest("GET", "https://ak01-eu.hwkapi.com/akamai/ua", nil)
	if err != nil {
		return "", err
	}
	req.Header = http.Header{
		"Accept-Encoding": {"gzip, deflate"},
		"X-Api-Key":       {apiKey},
		"X-Sec":           {"new"},
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", customErrors.StatusCodeNotExpectedError(resp.StatusCode)

	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}

func getConfig(configSource string) (string, error) {
	str := base64.StdEncoding.EncodeToString([]byte(configSource))
	jsonData := strings.NewReader(`{"body":"` + str + `"}`)

	req, err := http.NewRequest("POST", "https://ak-ppsaua.hwkapi.com/006180d12cf7", jsonData)
	if err != nil {
		return "", err
	}
	req.Header = http.Header{
		"content-type": {"application/json"},
		//"Accept-Encoding": {"gzip, deflate"},
		"X-Api-Key": {apiKey},
		"X-Sec":     {"new"},
	}
	url, _ := url.Parse("http://localhost:8888")
	client := http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(url)},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", customErrors.StatusCodeNotExpectedError(resp.StatusCode)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}
