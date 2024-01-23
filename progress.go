package main

import (
	"encoding/json"
	"errors"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) IncrQuestion(correct bool) tea.Cmd {
	m.MultipleAnswerProgress = make([]string, 0)
	m.BlankIndex = 0
	m.BlankAnswers = make(map[string]string)
	var cmd tea.Cmd
	if correct {
		m.Finished[m.Question] = struct{}{}
		cmd = m.QuestionProgress.SetPercent(float64(len(m.Finished)) / float64(len(m.Questions)))
	}

	if len(m.Finished) == len(m.Questions) {
		m.State = ModelStateQuitting

		// Reset progress
		dat, err := os.ReadFile(ProgressFile)
		if err != nil {
			panic(err)
		}
		progress := make(map[string]map[int]struct{})
		err = json.Unmarshal(dat, &progress)
		if err != nil {
			panic(err)
		}
		delete(progress, m.FileID)
		dat, err = json.Marshal(progress)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(ProgressFile, dat, os.ModePerm)
		if err != nil {
			panic(err)
		}

		return tea.Quit
	}

	for i := m.Question + 1; i < len(m.Questions); i++ {
		if _, exists := m.Finished[i]; !exists {
			m.Question = i
			m.UpdateProgress()
			return cmd
		}
	}
	// Have to start over
	for i := 0; i < len(m.Questions); i++ {
		if _, exists := m.Finished[i]; !exists {
			m.Question = i
			m.UpdateProgress()
			return cmd
		}
	}
	panic("couldn't find next question")
}

func (m *Model) UpdateProgress() {
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
		for v := range progress[m.FileID] {
			m.Finished[v] = struct{}{}
		}
		_, exists := progress[m.FileID][m.Question]
		if exists {
			for i := m.Question; i < len(m.Questions); i++ {
				if _, exists := progress[m.FileID][i]; !exists {
					m.Question = i
					break
				}
			}
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
