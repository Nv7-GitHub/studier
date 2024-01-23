package main

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
)

const ProgressFile = "progress.json"

func (m *Model) QuizStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.QuestionInput, cmd = m.QuestionInput.Update(msg)
	var ncmd tea.Cmd
	m.QuestionViewport, ncmd = m.QuestionViewport.Update(msg)
	cmd = tea.Batch(cmd, ncmd)
	if k, ok := msg.(tea.KeyMsg); ok && k.Type == tea.KeyEnter {
		ans := m.QuestionInput.Value()
		m.QuestionInput.Reset()

		// Process ans
		switch a := m.Questions[m.Question].Answer.(type) {
		case QuestionAnswerSingle:
			if CompareAnswer(ans, a.Answer) {
				return m, tea.Batch(cmd, m.IncrQuestion())
			} else {
				m.State = ModelStateQuestionResult
			}

		case QuestionAnswerMultiple:
			cont := false
			for _, v := range a.Answers {
				if CompareAnswer(ans, v) {
					cont = true
					break
				}
			}
			if cont {
				m.MultipleAnswerProgress = append(m.MultipleAnswerProgress, ans)
				if len(m.MultipleAnswerProgress) == len(a.Answers) {
					return m, tea.Batch(cmd, m.IncrQuestion())
				}
			} else {
				m.State = ModelStateQuestionResult
			}

		case QuestionAnswerBlanks:
			if CompareAnswer(ans, a.Answers[a.Order[m.BlankIndex]]) {
				m.BlankAnswers[a.Order[m.BlankIndex]] = ans
				m.BlankIndex++
				if m.BlankIndex == len(a.Order) {
					return m, tea.Batch(cmd, m.IncrQuestion())
				}
			} else {
				m.State = ModelStateQuestionResult
			}
		}
	}

	return m, cmd
}

func (m *Model) IncrQuestion() tea.Cmd {
	m.MultipleAnswerProgress = make([]string, 0)
	m.BlankIndex = 0
	m.BlankAnswers = make(map[string]string)
	m.Finished[m.Question] = struct{}{}
	cmd := m.QuestionProgress.SetPercent(float64(len(m.Finished)) / float64(len(m.Questions)))

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

func (m *Model) RenderQuestion() string {
	out := &strings.Builder{}
	for _, t := range m.Questions[m.Question].Text {
		switch t.Kind {
		case QuestionTextKindText:
			out.WriteString(QuestionStyle.Render(t.Value))

		case QuestionTextKindBlank:
			if _, exists := m.BlankAnswers[t.Value]; exists {
				out.WriteString(BlankStyle.Render(m.BlankAnswers[t.Value]))
			} else {
				out.WriteString(BlankStyle.Render(t.Value))
			}
		}
	}
	return wordwrap.String(out.String(), m.QuestionViewport.Width-4)
}

func CompareAnswer(a string, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

func (m *Model) QuizStateView() string {
	m.QuestionViewport.SetContent(m.RenderQuestion() + "\n\n" + m.QuestionInput.View() + "\n\n" + ListAnswerStyle.Render(strings.Join(m.MultipleAnswerProgress, "\n")))
	return m.QuestionProgress.View() + "\n\n" + m.QuestionViewport.View()
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
