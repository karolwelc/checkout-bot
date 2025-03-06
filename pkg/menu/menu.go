package menu

import (
	"fmt"
	"os"

	"github.com/ProjektPLTS/pkg/colors"
)

var (
	LogoColor = colors.Yellow
	Version   = "0.0.1"
)

func MainMenu() {
	fmt.Println(LogoColor(`____  _____________________________________`))
	fmt.Println(LogoColor(`__  |/ ___  _______  /___    ____  ___  __ \`))
	fmt.Println(LogoColor(`__    /__  __/  __  / __  /| |__  / _  / / /`))
	fmt.Println(LogoColor(`_    | _  /___  _  /___  ___ __/ /  / /_/ /`))
	fmt.Println(LogoColor(`/_/|_| /_____/  /_____/_/  |_/___/  \____/ ` + Version))
}

func SiteSelection() {
	fmt.Println("[1] Zalando")
	fmt.Println("[2] Asos")
}

func TaskSelection(site string) (tasks []string) {
	taskFilePath := "./" + site
	dir, _ := os.ReadDir(taskFilePath)
	for i, f := range dir {
		fmt.Println("[" + fmt.Sprint(i+1) + "] " + f.Name())
		tasks = append(tasks, f.Name())
	}
	return tasks
}
