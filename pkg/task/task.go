package task

import (
	"math/rand"

	tls "github.com/ProjektPLTS/pkg/tls"
	"github.com/ProjektPLTS/pkg/tls/http/cookiejar"
)

type AddressType struct {
	PhoneNum string

	ID string

	FirstName string
	LastName  string

	Street1 string
	Street2 string

	City    string
	ZipCode string

	Country string
}

type ProductType struct {
	URL string
	Id  string

	Name  string
	Image string

	Sizes []string
	Skus  []string
	Stock []string

	SelectedSku  string
	SelectedSize string
}

/*
Dynamic task is fields that will be changed each time we save task.csv or proxy.txt file
For example we do not want to change (erase) cookie jar so we do not touch that in changeMonitor.go file
*/
type DynamicTask struct {
	Mode            string
	DummyProductURL string
	Product         ProductType
	WantedSize      string
	RandomSizeIfOOS bool
	Email           string
	Password        string
	PaymentMethod   string
	Discount        string

	Address struct {
		Delivery AddressType
		Billing  AddressType
	}
	KeepAccountAddress bool

	Profile   Profile
	ProxyList []string
}

type TaskStruct struct {
	DynamicTask

	ID    string
	IDint int

	MaxTask    string
	MaxTaskint int

	Site string

	ClientSettings tls.ClientSettings

	//HEADERS
	UserAgent          string
	Sec_ch_ua          string
	Sec_ch_ua_platform string

	Cookies *cookiejar.Jar

	Config Config

	proxyIndex int
}

func (t *TaskStruct) SetClientSetting() {
	var cSettings tls.ClientSettings
	cSettings.AddDefaults()
	//change this
	//shuffle the proxy list so each task starts with a different proxy if a same proxy list
	if len(t.ProxyList) != 0 {
		rand.Shuffle(len(t.ProxyList), func(i, j int) { t.ProxyList[i], t.ProxyList[j] = t.ProxyList[j], t.ProxyList[i] })
		selectedProxy := t.ProxyList[t.proxyIndex]
		cSettings.Proxy = selectedProxy
		cSettings.SkipCertChecks = true
	}
	t.ClientSettings = cSettings
}

func (t *TaskStruct) ChangeProxy() {
	t.proxyIndex += 1
	if t.proxyIndex >= len(t.ProxyList) {
		t.proxyIndex = 0
	}
	t.SetClientSetting()
}

type Options struct {
	URL    string
	Method string
}

type Profile struct {
	ProfileName        string
	BillingProfileName string

	Webhook string

	Card CardInfo
}

type CardInfo struct {
	Type       string
	CardNumber string
	ExpMonth   string
	ExpYear    string
	CVC        string
}

type Config struct {
	LicenseKey   string `json:"license_key"`
	Delay        int    `json:"delay"`
	MonitorDelay int    `json:"monitor_delay"`
	Retries      int    `json:"retry_amount"`
}

type userAgent struct {
	UserAgents []struct {
		Name            string `json:"name"`
		UserAgent       string `json:"userAgent"`
		SecChUa         string `json:"sec_ch_ua"`
		SecChUaPlatform string `json:"sec_ch_ua_platform"`
	} `json:"userAgents"`
}
