package main

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) QuizStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m *Model) QuizStateView() string {
	return "QUIZ"
}
