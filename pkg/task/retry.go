package task

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ProjektPLTS/pkg/colors"
	"github.com/ProjektPLTS/pkg/retry"
)

type Messages struct {
	Success string
	Error   string
}

func (t TaskStruct) Retry(retryableFunc retry.RetryableFunc, messages Messages, retryIf ...retry.RetryIfFunc) error {
	if len(retryIf) > 0 {
		return retry.Do(retryableFunc,
			retry.Attempts(uint(t.Config.Retries)),
			retry.Delay(time.Duration(t.Config.Delay)),
			retry.DelayType(retry.FixedDelay),
			retry.OnRetry(func(n uint, err error) {
				if err != nil {
					if messages.Error != "" {
						fmt.Println(t.TaskPrefix() + colors.Red(messages.Error+".., retrying ["+strconv.FormatUint(uint64(n+1), 10)+"/"+strconv.Itoa(t.Config.Retries)+"]: "+err.Error()))
					}
				} else {
					if messages.Success != "" {
						fmt.Println(t.TaskPrefix() + colors.Green(messages.Success))
					}
				}

			}),
			retry.RetryIf(retryIf[0]),
		)
	}
	return retry.Do(retryableFunc,
		retry.Attempts(uint(t.Config.Retries)),
		retry.Delay(time.Duration(t.Config.Delay)),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			if err != nil {
				if messages.Error != "" {
					fmt.Println(t.TaskPrefix() + colors.Red(messages.Error+".., retrying ["+strconv.FormatUint(uint64(n+1), 10)+"/"+strconv.Itoa(t.Config.Retries)+"]: "+err.Error()))
				}
			} else {
				if messages.Success != "" {
					fmt.Println(t.TaskPrefix() + colors.Green(messages.Success))
				}
			}

		}),
	)
}
