package colors

import (
	"github.com/gookit/color"
)

func Red(data string) string {
	red := color.FgLightRed.Render
	return red(data)
}
func LightRed(data string) string {
	LightRed := color.LightRed.Render
	return LightRed(data)
}
func Green(data string) string {
	green := color.FgLightGreen.Render
	return green(data)
}
func Yellow(data string) string {
	yellow := color.FgLightYellow.Render
	return yellow(data)
}
func DarkYellow(data string) string {
	yellow := color.FgYellow.Darken().Render
	return yellow(data)
}
func White(data string) string {
	white := color.FgLightWhite.Render
	return white(data)
}
func Cyan(data string) string {
	cyan := color.FgLightCyan.Render
	return cyan(data)
}

func Magenta(data string) string {
	Magenta := color.Magenta.Render
	return Magenta(data)
}
func LightMagenta(data string) string {
	Magenta := color.FgLightMagenta.Render
	return Magenta(data)
}

func LightBlue(data string) string {
	lBlue := color.FgLightBlue.Render
	return lBlue(data)
}
