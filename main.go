package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/ProjektPLTS/pkg/colors"
	"github.com/ProjektPLTS/pkg/inputs"
	"github.com/ProjektPLTS/pkg/menu"
	"github.com/ProjektPLTS/pkg/task"
	"github.com/ProjektPLTS/pkg/taskConsumer"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	for {
		menu.MainMenu()
		fmt.Print("\n\n")
		menu.SiteSelection()
		fmt.Print("\n")

		var selectedSite string
		switch inputs.InputWithText("Select site: ") {
		case "1":
			selectedSite = "Zalando"
		case "2":
			selectedSite = "Asos"
		default:
			fmt.Println(task.Prefix() + colors.Red("Invalid selection"))
			continue
		}
		fmt.Print("\n")
		taskFiles := menu.TaskSelection(selectedSite)
		fmt.Print("\n")

		taskIndex, err := strconv.Atoi(inputs.InputWithText("Select task file: "))
		if taskIndex-1 > len(taskFiles) {
			fmt.Println(task.Prefix() + colors.Red("Invalid selection"))
			continue
		}
		if err != nil {
			fmt.Println(err)
		}
		tasks := task.CreateTasks(selectedSite, taskFiles[taskIndex-1])

		taskConsumer.StartTasks(tasks)
	}
}
