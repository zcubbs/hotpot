package style

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

type Color string

const (
	Black   Color = "0"
	Red     Color = "1"
	Green   Color = "2"
	Yellow  Color = "3"
	Blue    Color = "4"
	Magenta Color = "5"
	Cyan    Color = "6"
	White   Color = "7"
)

func PrintColoredHeader(text string, foreground Color, background Color) {
	var style = lipgloss.NewStyle().
		Bold(true).
		MarginTop(1).
		PaddingLeft(1).
		Foreground(lipgloss.Color(foreground)).
		Background(lipgloss.Color(background)).
		Align(lipgloss.Left).
		Width(40)

	fmt.Println(style.Render(text))
}

func PrintColoredSuccess(text string) {
	var style = lipgloss.NewStyle().
		Align(lipgloss.Left).
		Background(lipgloss.Color("86")).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 1).
		Width(6)

	fmt.Println(style.Render("OK") + " " + text)
}

func PrintSuccess(text string) {
	var style = lipgloss.NewStyle().
		Align(lipgloss.Left).
		Foreground(lipgloss.Color("#04B575")).
		Padding(0, 1).
		Width(6)

	fmt.Println(style.Render("ok") + " " + text)
}

func PrintColoredError(text string) {
	var style = lipgloss.NewStyle().
		Align(lipgloss.Left).
		Background(lipgloss.Color("9")).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 1).
		Width(6)

	fmt.Println(style.Render("ERR") + " " + text)
}

func PrintError(text string) {
	var style = lipgloss.NewStyle().
		Align(lipgloss.Left).
		Foreground(lipgloss.Color("#FF3B30")). // red
		Padding(0, 1).
		Width(6)

	fmt.Println(style.Render("err") + " " + text)
}

func PrintColoredWarning(text string) {
	var style = lipgloss.NewStyle().
		Align(lipgloss.Left).
		Background(lipgloss.Color("11")).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 1).
		Width(6)

	fmt.Println(style.Render("WARN") + " " + text)
}

func PrintColoredInfo(text string) {
	var style = lipgloss.NewStyle().
		Align(lipgloss.Left).
		Background(lipgloss.Color("14")).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 1).
		Width(6)

	fmt.Println(style.Render("INFO") + " " + text)
}

func PrintInfo(text string) {
	var style = lipgloss.NewStyle().
		Align(lipgloss.Left).
		Foreground(lipgloss.Color("15")). // white
		Padding(0, 1).
		Width(6)

	fmt.Println(style.Render("info") + " " + text)
}

func PrintColoredDebug(text string) {
	var style = lipgloss.NewStyle().
		Align(lipgloss.Left).
		Background(lipgloss.Color("8")).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 1).
		Width(6)

	fmt.Println(style.Render("DEBG") + " " + text)
}
