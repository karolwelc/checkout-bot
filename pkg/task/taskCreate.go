package task

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"net/url"
)

const (
	profileFileName string = "./profiles.csv"
	configFileName  string = "./config.json"
)

func loadProfile(profileName string, records [][]string) (Profile, AddressType) {
	var address AddressType
	var cardInfo CardInfo
	var profile Profile
	for i := 1; i < len(records); i++ {
		if records[i][0] == profileName {
			address.FirstName = records[i][1]
			address.LastName = records[i][2]
			address.Street1 = records[i][4]
			address.Street2 = records[i][5]
			address.City = records[i][6]
			address.ZipCode = records[i][7]
			address.Country = records[i][8]
			address.PhoneNum = records[i][9]

			cardInfo.Type = records[i][11]
			cardInfo.CardNumber = records[i][12]
			cardInfo.ExpMonth = records[i][13]
			cardInfo.ExpYear = records[i][14]
			cardInfo.CVC = records[i][15]

			profile.ProfileName = records[i][0]
			profile.Webhook = records[i][10]
			profile.Card = cardInfo
		}
	}
	return profile, address
}

func CreateTasks(store string, taskFile string) []TaskStruct {
	taskFileName := store + "/" + taskFile

	TasksF, err := os.Open(taskFileName)
	if err != nil {
		log.Fatal("Unable to read task file ", err)
	}
	defer TasksF.Close()
	csvReaderTask := csv.NewReader(TasksF)
	taskRecords, err := csvReaderTask.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse task file as CSV : ", err)
	}

	ProfileF, err := os.Open(profileFileName)
	if err != nil {
		log.Fatal("Unable to read profiles file ", err)
	}
	defer ProfileF.Close()
	csvReaderProfile := csv.NewReader(ProfileF)
	profileRecords, err := csvReaderProfile.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse profile file as CSV : ", err)
	}

	tasks, err := createTasks(taskRecords, profileRecords, store)
	if err != nil {
		log.Fatal("Unable to create tasks ", err)
	}
	return tasks
}

func createTasks(records [][]string, profileRecords [][]string, site string) ([]TaskStruct, error) {
	var config Config
	var userAgent userAgent
	//parse config file
	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("Unable to read config file ", err)
	}
	defer configFile.Close()
	configBody, err := io.ReadAll(configFile)
	if err != nil {
		log.Fatal("Unable to read config file ", err)
	}
	if err := json.Unmarshal(configBody, &config); err != nil {
		log.Fatal("Unable to parse config file ", err)
	}
	//parse ua file
	uaFile, err := os.Open("./src/userAgents.json")
	if err != nil {
		log.Fatal("Unable to read userAgent file ", err)
	}
	defer uaFile.Close()
	uaBody, err := io.ReadAll(uaFile)
	if err != nil {
		log.Fatal("Unable to read userAgent file ", err)
	}
	if err := json.Unmarshal(uaBody, &userAgent); err != nil {
		log.Fatal("Unable to parse userAgent file ", err)
	}

	var tasks []TaskStruct

	maxTask := fmt.Sprint((len(records) - 1))
	var (
		profileI         int
		billingProfileI  int
		modeI            int
		checkoutProductI int
		sizeI            int
		preloadProductI  int
		paymentMethodI   int
		emailI           int
		passowrdI        int
		discountI        int
		keepAccAddressI  int
		proxyListI       int
	)

	for i, r := range records[0] {
		switch strings.ToLower(r) {
		case "profile":
			profileI = i
		case "billing profile":
			billingProfileI = i
		case "mode":
			modeI = i
		case "url/pid/variant":
			checkoutProductI = i
		case "size":
			sizeI = i
		case "preload url/pid/var":
			preloadProductI = i
		case "payment method":
			paymentMethodI = i
		case "account email", "email":
			emailI = i
		case "account password":
			passowrdI = i
		case "discount code":
			discountI = i
		case "keep account address":
			keepAccAddressI = i
		case "proxylist":
			proxyListI = i
		}

	}

	for i := 1; i < len(records); i++ {
		var task TaskStruct
		task.Site = site
		profile, address := loadProfile(records[i][profileI], profileRecords)
		task.Profile = profile
		task.Profile.BillingProfileName = records[i][billingProfileI]
		task.Address.Delivery = address
		if records[i][billingProfileI] != "" {
			_, task.Address.Billing = loadProfile(records[i][profileI], profileRecords)
		} else {
			task.Address.Billing = address
		}

		task.Mode = records[i][modeI]
		task.Product.URL = records[i][checkoutProductI]
		task.WantedSize = records[i][sizeI]
		if strings.Contains(records[i][sizeI], "R") { //todo: be changed maybe - konik
			task.RandomSizeIfOOS = true
			task.WantedSize = strings.Replace(records[i][sizeI], "R", "", -1)
		}
		task.DummyProductURL = records[i][preloadProductI]
		task.PaymentMethod = records[i][paymentMethodI]
		task.Email = records[i][emailI]
		task.Password = records[i][passowrdI]
		task.Discount = records[i][discountI]
		keepAddress := strings.ToLower(records[i][keepAccAddressI])
		if keepAddress == "yes" || keepAddress == "true" {
			task.KeepAccountAddress = true
		}
		task.ID = fmt.Sprint(i)
		task.IDint = i
		proxyListName := records[i][proxyListI]
		if proxyListName != "" {
			proxiesByte, err := os.ReadFile("./proxies/" + proxyListName)
			if err != nil {
				return tasks, err
			}
			if len(proxiesByte) != 0 {
				proxyListParsed, err := parseProxyList(proxiesByte)
				if err != nil {
					return tasks, err
				}
				task.ProxyList = proxyListParsed
			}

		}
		task.SetClientSetting()
		for _, v := range userAgent.UserAgents {
			if v.Name == "chrome105" {
				task.UserAgent = v.UserAgent
				task.Sec_ch_ua = v.SecChUa
				task.Sec_ch_ua_platform = v.SecChUaPlatform
			}

		}
		task.MaxTask = maxTask
		task.MaxTaskint = len(records) - 1
		task.Config = config
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func parseProxyList(file []byte) ([]string, error) {
	var proxyList []string
	proxiesStr := string(file)
	proxiesSplit := strings.Split(proxiesStr, "\n")
	for i := 0; i < len(proxiesSplit); i++ {
		parsedProxy, err := parseProxy(proxiesSplit[i])
		if err != nil {
			return []string{}, err
		}
		proxyList = append(proxyList, parsedProxy)
	}
	return proxyList, nil
}

func parseProxy(proxy string) (string, error) {
	proxy = strings.Replace(proxy, "\r", "", -1)
	proxySplit := strings.Split(proxy, ":")
	if len(proxySplit) == 4 {
		proxy := "http://" + proxySplit[2] + ":" + proxySplit[3] + "@" + proxySplit[0] + ":" + proxySplit[1]
		_, err := url.Parse(proxy)
		if err != nil {
			return "", err
		}
		return proxy, nil
	} else if len(proxySplit) == 2 {
		proxy := "http://" + proxySplit[0] + ":" + proxySplit[1]
		_, err := url.Parse(proxy)
		if err != nil {
			return "", err
		}
		return proxy, nil

	} else {
		return "", errors.New("error wrong proxy format")
	}

}
