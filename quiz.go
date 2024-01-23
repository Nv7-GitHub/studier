package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const ProgressFile = "progress.json"

func (m *Model) QuizStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.QuestionInput, cmd = m.QuestionInput.Update(msg)
	if k, ok := msg.(tea.KeyMsg); ok && k.Type == tea.KeyEnter {
		ans := m.QuestionInput.Value()
		m.QuestionInput.Reset()

		// Process ans
		fmt.Println(ans)
	}

	return m, cmd
}

func (m *Model) QuizStateView() string {
	m.QuestionViewport.SetContent("Question\n\n" + m.QuestionInput.View() + "\n\nAnswer\n\nCool!")
	return m.QuestionProgress.View() + "\n\n" + m.QuestionViewport.View()
}

func (m *Model) UpdateProgress() {
	m.QuestionProgress.SetPercent(float64(len(m.Finished)) / float64(len(m.Questions)))

	// Read progress
	if _, err := os.Stat(ProgressFile); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(ProgressFile, []byte("{}"), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	dat, err := os.ReadFile(ProgressFile)
	if err != nil {
		panic(err)
	}
	progress := make(map[string]map[int]struct{})
	err = json.Unmarshal(dat, &progress)
	if err != nil {
		panic(err)
	}

	// Add progress
	_, exists := progress[m.FileID]
	if exists && len(m.Finished) < len(progress[m.FileID]) {
		for v := range m.Finished {
			progress[m.FileID][v] = struct{}{}
		}
	}
	progress[m.FileID] = m.Finished

	// Save
	dat, err = json.Marshal(progress)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(ProgressFile, dat, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
