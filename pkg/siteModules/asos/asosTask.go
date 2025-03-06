package asos

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ProjektPLTS/pkg/colors"
	"github.com/ProjektPLTS/pkg/customErrors"
	taskpkg "github.com/ProjektPLTS/pkg/task"
)

var stores = []string{"de", "es", "it", "fr", "ru", "row", "roe", "nl", "se", "pl", "dk"}

// Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36
type AsosTask struct {
	*taskpkg.TaskStruct

	bearerToken string
	bagId       string

	xsrf string

	loginUrl     string
	securedTouch string
	appSessId    string
	checkSum     string
	deviceID     string
	instanceuuid string

	storeCode string

	akamaiEndpoint string
	configSC       string
}

func (task *AsosTask) StartAsos(wg *sync.WaitGroup) {
	if contains(stores, strings.ToLower(task.Address.Delivery.Country)) {
		task.storeCode = task.Address.Delivery.Country
	} else {
		task.storeCode = "ROE"
	}
	defer wg.Done()
	err := task.Preload()
	if err != nil {
		fmt.Println(err)
	}
}

func (task *AsosTask) Preload() error {
	var err error
	task.Retry(func() error {
		err = task.anonLogin()
		if err != nil {
			return err
		}
		return nil
	}, taskpkg.Messages{
		Error: "Error Getting session",
	},
		func(err error) bool {
			switch err {
			case customErrors.Forbidden:
				fmt.Println(task.TaskPrefix() + colors.Red("Rotating proxy..."))
				task.ChangeProxy()
			}
			return true
		})
	if err != nil {
		return err
	}
	task.Retry(func() error {
		err = task.getAkamaiScriptAnon()
		if err != nil {
			return err
		}

		return nil
	}, taskpkg.Messages{
		Error: "Error Getting akamai script",
	},
		func(err error) bool {
			switch err {
			case customErrors.Forbidden:
				fmt.Println(task.TaskPrefix() + colors.Red("Rotating proxy..."))
				task.ChangeProxy()
			}
			return true
		})
	if err != nil {
		return err
	}
	task.Retry(func() error {
		err = task.getBag()
		if err != nil {
			return err
		}
		return nil
	}, taskpkg.Messages{
		Error: "Error Getting bag",
	},
		func(err error) bool {
			switch err {
			case customErrors.Forbidden:
				fmt.Println(task.TaskPrefix() + colors.Red("Rotating proxy..."))
				task.ChangeProxy()
			}
			return true
		})
	if err != nil {
		return err
	}
	task.Retry(func() error {

		err := task.submitSensorData("https://www.asos.com")
		if err != nil {
			return err
		}

		err = task.atc()
		if err != nil {
			return err
		}
		return nil
	}, taskpkg.Messages{
		Error: "Error adding to cart",
	},
		func(err error) bool {
			switch err {
			case customErrors.Forbidden:
				fmt.Println(task.TaskPrefix() + colors.Red("Rotating proxy..."))
				task.ChangeProxy()
			}
			return true
		})
	if err != nil {
		return err
	}
	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
