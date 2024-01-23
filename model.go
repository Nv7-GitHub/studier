package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ModelState int

const (
	ModelStateFileInput ModelState = iota
	ModelStateQuiz
	ModelStateQuitting
	ModelStateError
)

type Model struct {
	Err error

	State ModelState

	FilePicker filepicker.Model
	Questions  []Question
	FileID     string

	Question         int
	QuestionInput    textinput.Model
	Finished         map[int]struct{}
	QuestionProgress progress.Model
	QuestionViewport viewport.Model
}

func NewModel() *Model {
	m := &Model{}

	m.FilePicker = filepicker.New()
	m.QuestionProgress = progress.New(progress.WithDefaultGradient())
	m.QuestionInput = textinput.New()
	m.QuestionInput.Placeholder = "Answer"
	m.QuestionInput.Focus()
	m.QuestionViewport = viewport.New(m.QuestionProgress.Width, 0)
	m.QuestionViewport.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	m.Finished = make(map[int]struct{})

	return m
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.FilePicker.Init(), textinput.Blink)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.State = ModelStateQuitting
			return m, tea.Quit
		}

	case progress.FrameMsg:
		progressModel, cmd := m.QuestionProgress.Update(msg)
		m.QuestionProgress = progressModel.(progress.Model)
		if cmd != nil {
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.QuestionProgress.Width = msg.Width
		m.QuestionViewport.Width = msg.Width
		m.QuestionViewport.Height = msg.Height - 4
		m.QuestionInput.Width = msg.Width - 8
		m.QuestionInput.Placeholder = "Answer" + strings.Repeat(" ", msg.Width-13)
	}

	switch m.State {
	case ModelStateError:
		return m, tea.Quit

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

	case ModelStateError:
		return "" // Display error after

	default:
		return "Quitting..."
	}
}
