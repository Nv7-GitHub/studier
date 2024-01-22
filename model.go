package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelState int

const (
	ModelStateFileInput ModelState = iota
)

type Model struct {
	State ModelState

	ParseSpinner spinner.Model
}

func NewModel() *Model {
	m := &Model{}
	m.ParseSpinner = spinner.New()
	m.ParseSpinner.Spinner = spinner.Dot

	return m
}

func (m *Model) Init() tea.Cmd {

}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.State {
	case ModelStateFileInput:
		return m.InputParseState(msg)

	default:
		return m, nil
	}
}

func (m *Model) View() string {
	
}
