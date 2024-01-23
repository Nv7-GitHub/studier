package main

import (
	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelState int

const (
	ModelStateFileInput ModelState = iota
	ModelStateQuiz
	ModelStateQuitting
)

type Model struct {
	State ModelState

	FilePicker filepicker.Model
	Questions  []Question
}

func ChangeStateCmd() tea.Msg {
	return nil
}

func NewModel() *Model {
	m := &Model{}

	// Input parser
	m.FilePicker = filepicker.New()

	return m
}

func (m *Model) Init() tea.Cmd {
	return m.FilePicker.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.State = ModelStateQuitting
			return m, tea.Quit
		}
	}

	switch m.State {
	case ModelStateFileInput:
		return m.InputParseState(msg)

	case ModelStateQuiz:
		return m.QuizStateUpdate(msg)

	default:
		return m, nil
	}
}

func (m *Model) View() string {
	switch m.State {
	case ModelStateFileInput:
		return m.InputParseStateView()

	case ModelStateQuiz:
		return m.QuizStateView()

	default:
		return "Quitting..."
	}
}
