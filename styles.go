package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
var BlankStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)

func (m *Model) HandleErr(msg string) {
	// Red
	fmt.Println(ErrStyle.Render(msg))
	os.Exit(1)
}
