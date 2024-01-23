package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := NewModel()
	tea.NewProgram(m, tea.WithOutput(os.Stderr)).Run()
}
