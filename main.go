package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := NewModel()
	tea.NewProgram(m, tea.WithOutput(os.Stderr)).Run()
	if m.Err != nil {
		fmt.Println(m.Err)
		os.Exit(1)
	}
	if m.Done {
		fmt.Println(CorrectStyle.Render("Done!"))
	}
}
