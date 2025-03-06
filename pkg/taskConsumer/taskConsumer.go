package taskConsumer

import (
	"sync"

	"github.com/ProjektPLTS/pkg/colors"
	"github.com/ProjektPLTS/pkg/inputs"
	"github.com/ProjektPLTS/pkg/siteModules/asos"
	"github.com/ProjektPLTS/pkg/siteModules/zalando"
	"github.com/ProjektPLTS/pkg/task"
)

//some explanation:
//have to use pointer when assigning task var, to be able to make changes properly

func StartTasks(tasks []task.TaskStruct) {
	var wg sync.WaitGroup
	c := make(chan []task.TaskStruct)
	for i := 0; i < len(tasks); i++ {
		task := tasks[i]
		if task.Site == "Zalando" {
			if task.IDint == 1 {
				go task.MonitorChanges("tasks.csv", "proxies.txt", c)
			}
			go func() {
				for {
					task.WaitForChange(c)
				}
			}()
			zalTask := zalando.ZalTask{TaskStruct: &task}
			wg.Add(1)
			go zalTask.StartZalando(&wg)
		}
		if task.Site == "Asos" {
			if task.IDint == 1 {
				go task.MonitorChanges("tasks.csv", "proxies.txt", c)
			}
			go func() {
				for {
					task.WaitForChange(c)
				}
			}()
			asosTask := asos.AsosTask{TaskStruct: &task}
			wg.Add(1)
			go asosTask.StartAsos(&wg)
		}
	}

	wg.Wait()
	inputs.InputWithText(colors.Yellow("Press enter to continue\n"))
}
