package log

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

const (
	check = "✔️"
	cross = "❌"
	green = lipgloss.Color("#38c172")
	red   = lipgloss.Color("#e3342f")
)

var (
	redStyle   = lipgloss.NewStyle().Foreground(red)
	greenStyle = lipgloss.NewStyle().Foreground(green)
)

func Info(msg string) {
	fmt.Println(greenStyle.Render(check) + " " + msg)
}

func Fatal(err error) {
	fmt.Println(redStyle.Render(cross) + " " + err.Error())
	os.Exit(2)
}
