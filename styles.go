package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
var CorrectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
var ListAnswerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("31")).Bold(true)
var QuestionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
var BlankStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)

func (m *Model) HandleErr(err error) tea.Cmd {
	// Red
	m.Err = err
	m.State = ModelStateError
	return tea.ClearScreen
}
