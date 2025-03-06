package task

import (
	"fmt"

	"github.com/ProjektPLTS/pkg/colors"
	"github.com/fsnotify/fsnotify"
)

func (task *TaskStruct) MonitorChanges(taskFile string, proxyFile string, c chan []TaskStruct) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err) //remove in production
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op == fsnotify.Write || event.Op == fsnotify.Remove {
					fmt.Println(Prefix() + colors.LightBlue("file "+event.Name+" changed"))
					if event.Op == fsnotify.Remove {
						err = watcher.Add(event.Name)
						if err != nil {
							fmt.Println(err) //remove in production
						}
					}
					tasks := CreateTasks(task.Site, taskFile)
					for i := 1; i <= task.MaxTaskint; i++ {
						c <- tasks
					}

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()
	err = watcher.Add("./" + task.Site + "/" + taskFile)
	if err != nil {
		fmt.Println(err) //remove in production
	}
	err = watcher.Add("proxies/" + proxyFile)
	if err != nil {
		fmt.Println(err) //remove in production
	}

	<-done

}

func (task *TaskStruct) WaitForChange(c chan []TaskStruct) {
	tasks := <-c
	newTask := tasks[task.IDint-1]
	task.DynamicTask = newTask.DynamicTask
	task.SetClientSetting()
}
