package concolor

import (
	"fmt"
)

var colorFormat = "\x1b[%dm%s\x1b[0m"

func colorize(colorCode int, s string) string {
	return fmt.Sprintf(colorFormat, colorCode, s)
}

type Color int

const (
	ColorWhite         = 0
	ColorRed           = 31
	ColorGreen         = 32
	ColorYellow        = 33
	ColorBlue          = 34
	ColorPurple        = 35
	ColorLightGreen    = 36
	ColorGray          = 37
	ColorRedBackground = 41
)

func Blue(s string) string {
	return colorize(ColorBlue, s)
}

func LightGreen(s string) string {
	return colorize(ColorLightGreen, s)
}

func Purple(s string) string {
	return colorize(ColorPurple, s)
}

func White(s string) string {
	return colorize(ColorWhite, s)
}

func Gray(s string) string {
	return colorize(ColorGray, s)
}

func Red(s string) string {
	return colorize(ColorRed, s)
}

func RedBackground(s string) string {
	return colorize(ColorRedBackground, s)
}

func Green(s string) string {
	return colorize(ColorGreen, s)
}

func Yellow(s string) string {
	return colorize(ColorYellow, s)
}
