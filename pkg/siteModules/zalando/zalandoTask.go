package zalando

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ProjektPLTS/pkg/colors"
	"github.com/ProjektPLTS/pkg/customErrors"
	taskpkg "github.com/ProjektPLTS/pkg/task"
)

type ZalTask struct {
	*taskpkg.TaskStruct

	akamaiEndpoint string
	configSC       string

	loginUrl     string
	csrfToken    string
	redirect_url string

	cartId string
	timer  bool
	time   time.Time

	checkoutId string
	eTag       string

	finishedPreload bool
}

func (task ZalTask) StartZalando(wg *sync.WaitGroup) {
	defer wg.Done()
	task.Site = "Zalando"

	err := task.preloadApp()
	if err != nil {
		fmt.Println(task.TaskPrefix() + colors.Red("Preload failed!"))
		return
	}

	fmt.Println(task.TaskPrefix() + colors.Cyan("Preload successful!"))
	task.finishedPreload = true
	task.Product = taskpkg.ProductType{
		URL: task.Product.URL,
	}
	err = task.finishCheckoutApp()
	if err != nil {
		fmt.Println(task.TaskPrefix() + colors.Red("Checkout failed!"))
		return
	}
	fmt.Println(task.TaskPrefix() + colors.Magenta("Checkout successful!"))

}

func (task *ZalTask) preloadApp() error {
	var err error
	fmt.Println(task.TaskPrefix() + colors.DarkYellow("Logging in..."))
	task.Retry(func() error {
		err = task.login()
		if err != nil {
			return err
		}
		return nil
	}, taskpkg.Messages{
		Error:   "Error logging in",
		Success: "Successfully logged in.",
	},
		func(err error) bool {
			switch err {
			case customErrors.Unauthorized:
				fmt.Println(task.TaskPrefix() + colors.Red("Invalid login credentials!"))
				return false
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
		var products []string
		if task.cartId == "" {
			products, err = task.getCartedProducts()
			if err != nil {
				return err
			}
		}
		err = task.deleteProducts(products)
		if err != nil {
			return err
		}
		return nil
	}, taskpkg.Messages{
		Error:   "Error deleting items from cart",
		Success: "Emptied cart successfully",
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

	//need custom retry func here to be able to use sess struct values
	var retries int = 1
	for retries <= task.Config.Retries {
		err = task.addToCart()
		if err != nil {
			if retries == task.Config.Retries {
				return err
			}
			retries++
			if err == customErrors.Forbidden {
				fmt.Println(task.TaskPrefix() + colors.Red("Rotating proxy..."))
				task.ChangeProxy()
				task.submitSensorData("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country))
				continue
			}
			fmt.Println(task.TaskPrefix() + colors.Red("Error adding product ") + colors.Cyan(task.Product.URL) + colors.Red(" in size ") + colors.Cyan(task.WantedSize) + colors.Red(" to cart..., retrying ["+strconv.Itoa(retries)+"/"+strconv.Itoa(task.Config.Retries)+"]: "+err.Error()))
		} else {
			fmt.Println(task.TaskPrefix() + colors.Green("Product ") + colors.Cyan(task.Product.Name) + colors.Green(" in size ") + colors.Cyan(task.Product.SelectedSize) + colors.Green(" added to cart successfully."))
			break
		}
	}
	fmt.Println(task.TaskPrefix() + colors.DarkYellow("Getting checkout session..."))

	var nextStepResp nextStepResp
	task.Retry(func() error {
		nextStepResp, err = task.nextStepApp()
		if err != nil {
			return err
		}
		return nil
	}, taskpkg.Messages{Error: "Error getting checkout info."},
		func(err error) bool {
			switch err {
			case customErrors.Forbidden:
				fmt.Println(task.TaskPrefix() + colors.Red("Rotating proxy..."))
				task.ChangeProxy()
				task.submitSensorData("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country))
			}
			return true
		})
	if err != nil {
		return err
	}
	//controlling flow, depending on next step response
	switch nextStepResp.URL {
	case "/checkout/address":
		var defaultBilling, defaultShipping bool
		//if user wants to keep the address saved in account
		if task.KeepAccountAddress {
			task.Retry(func() error {
				err = task.defaultAddressExists()
				if err != nil {
					return err
				}
				return nil
			}, taskpkg.Messages{Error: "Error getting addresses."},
				func(err error) bool {
					switch err {
					case customErrors.Forbidden:
						fmt.Println(task.TaskPrefix() + colors.Red("Rotating proxy..."))
						task.ChangeProxy()
						task.submitSensorData("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country))
					}
					return true
				})
			if err != nil {
				return err
			}
			if task.Address.Billing.City == "" {
				defaultBilling = true
			}
			if task.Address.Delivery.City == "" {
				defaultShipping = true
			}
			if defaultBilling && defaultShipping {
				fmt.Println(task.TaskPrefix() + colors.Yellow("No default address set, adding one..."))
			} else if defaultBilling {
				fmt.Println(task.TaskPrefix() + colors.Yellow("No default billing address set, adding one..."))
			} else if defaultShipping {
				fmt.Println(task.TaskPrefix() + colors.Yellow("No default shipping address set, adding one..."))
			}
		}

		for i := 0; i < 1; i++ {
			err = task.submitSensorData("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country))
			if err != nil {
				return err
			}
		}

		task.Retry(func() error {
			if task.Profile.BillingProfileName == "" || (!defaultBilling && !defaultShipping) {
				//if billing same as delivery
				err = task.submitAddressWeb(task.Address.Delivery, true, true)
				if err != nil {
					return err
				}
			} else {
				//if billing different from delivery
				err = task.submitAddressWeb(task.Address.Delivery, true, false)
				if err != nil {
					return err
				}
				err = task.submitAddressWeb(task.Address.Billing, false, true)
				if err != nil {
					return err
				}
			}
			return nil
		}, taskpkg.Messages{Error: "Error submitting address."},
			func(err error) bool {
				switch err {
				case customErrors.Forbidden:
					fmt.Println(task.TaskPrefix() + colors.Red("Rotating proxy..."))
					task.ChangeProxy()
					task.submitSensorData("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country))
				}
				return true
			})
		if err != nil {
			return err
		}
		task.Retry(func() error {
			err = task.submitPaymentWeb()
			if err != nil {
				return err
			}
			return nil
		}, taskpkg.Messages{Error: "Error submitting payment method!"},
			func(err error) bool {
				switch err {
				case customErrors.Forbidden:
					fmt.Println(task.TaskPrefix() + colors.Red("Rotating proxy..."))
					task.ChangeProxy()
					task.submitSensorData("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country))
				}
				return true
			})
		if err != nil {
			return err
		}

	case "/checkout/payment":
		task.Retry(func() error {
			err = task.submitPaymentWeb()
			if err != nil {
				return err
			}
			return nil
		}, taskpkg.Messages{Error: "Error submitting payment method!"},
			func(err error) bool {
				switch err {
				case customErrors.Forbidden:
					fmt.Println(task.TaskPrefix() + colors.Red("Rotating proxy..."))
					task.ChangeProxy()
					task.submitSensorData("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country))
				}
				return true
			})
		if err != nil {
			return err
		}
	}
	fmt.Println(task.TaskPrefix() + colors.Green("Got session"))
	if task.Discount != "" {
		task.Retry(func() error {
			err = task.submitDiscountCode()
			if err != nil {
				return err
			}
			return nil
		}, taskpkg.Messages{Error: "Error submitting discount code " + task.Discount + " !", Success: "Successfully submitted discount code!"},
			func(err error) bool {
				switch err {
				case customErrors.Forbidden:
					fmt.Println(task.TaskPrefix() + colors.Red("Rotating proxy..."))
					task.ChangeProxy()
					task.submitSensorData("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country))
				}
				return true
			})
	}
	task.Retry(func() error {
		err = task.getCheckoutWeb()
		if err != nil {
			return err
		}
		return nil
	}, taskpkg.Messages{Error: "Error getting checkout page!"},
		func(err error) bool {
			switch err {
			case customErrors.Forbidden:
				fmt.Println(task.TaskPrefix() + colors.Red("Rotating proxy..."))
				task.ChangeProxy()
				task.submitSensorData("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country))
			}
			return true
		})
	if err != nil {
		return err
	}
	task.Retry(func() error {
		err = task.deleteProducts([]string{task.Product.SelectedSku})
		if err != nil {
			return err
		}
		return nil
	}, taskpkg.Messages{
		Error:   "Error deleting dummy item",
		Success: "Deleted dummy item successfully",
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

func (task *ZalTask) finishCheckoutApp() error {
	var err error
	task.Retry(func() error {
		err = task.getProductInfo()
		if err != nil {
			return err
		}
		err = task.chooseSize()
		if err != nil {
			return err
		}
		return nil
	}, taskpkg.Messages{
		Error: "Error getting / Selecting size",
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
	if task.timer {
		localTimer := task.time.In(time.Local)
		fmt.Println(task.TaskPrefix() + colors.DarkYellow("Product with timer, waiting for ") + localTimer.Format("15:04:05"))

		for time.Now().Unix()+2 < localTimer.Unix() {
			time.Sleep(time.Second * 1)
		}
		fmt.Println(task.TaskPrefix() + colors.Cyan("Waking up..."))
	}

	var retries int = 1
	for retries <= task.Config.Retries {
		err = task.atcApp()
		if err != nil {
			if retries == task.Config.Retries {
				return err
			}
			retries++
			if err == customErrors.Forbidden {
				fmt.Println(task.TaskPrefix() + colors.Red("Rotating proxy..."))
				task.ChangeProxy()
				task.submitSensorData("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country))
				continue
			}
			fmt.Println(task.TaskPrefix() + colors.Red("Error adding product ") + colors.Cyan(task.Product.URL) + colors.Red(" in size ") + colors.Cyan(task.WantedSize) + colors.Red(" to cart..., retrying ["+strconv.Itoa(retries)+"/"+strconv.Itoa(task.Config.Retries)+"]: "+err.Error()))
		} else {
			fmt.Println(task.TaskPrefix() + colors.Green("Product ") + colors.Cyan(task.Product.Name) + colors.Green(" in size ") + colors.Cyan(task.Product.SelectedSize) + colors.Green(" added to cart successfully."))
			break
		}
	}

	task.Retry(func() error {
		err = task.postCheckoutWeb()
		if err != nil {
			return err
		}
		return nil
	}, taskpkg.Messages{Error: "Error checking out!"},
		func(err error) bool {
			switch err {
			case customErrors.Forbidden:
				fmt.Println(task.TaskPrefix() + colors.Red("Rotating proxy..."))
				task.ChangeProxy()
				task.submitSensorData("https://www.zalando." + strings.ToLower(task.Address.Delivery.Country))
			}
			return true
		})
	if err != nil {
		return err
	}
	return nil
}
