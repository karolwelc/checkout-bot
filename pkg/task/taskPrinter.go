package task

import (
	"time"

	"github.com/ProjektPLTS/pkg/colors"
)

var (
	line = " | "
)

func Prefix() string {
	dt := time.Now()
	dtf := dt.Format("15:04:05.00")
	return colors.White(dtf) + colors.Cyan(line)
}
func (t TaskStruct) TaskPrefix() string {
	ids := t.ID
	maxTask := t.MaxTask
	spacingCounter := len(maxTask) - len(ids)
	var spacing string
	for i := 0; i < spacingCounter; i++ {
		spacing += " "
	}
	dt := time.Now()
	dtf := dt.Format("15:04:05.00")
	return colors.White(dtf) + colors.Cyan(line) + colors.White(t.Site) + colors.Cyan(line) + colors.White("Task "+ids+spacing) + colors.LightMagenta(" > ")
}
