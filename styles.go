package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
var BlankStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)

func (m *Model) HandleErr(err error) tea.Cmd {
	// Red
	m.Err = err
	m.State = ModelStateError
	return tea.ClearScreen
}
