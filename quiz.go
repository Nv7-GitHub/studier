package main

import (
	"fmt"
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
				return m, tea.Batch(cmd, m.IncrQuestion(true))
			} else {
				m.State = ModelStateQuestionResult
				m.IncorrectAnswer = ans
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
					return m, tea.Batch(cmd, m.IncrQuestion(true))
				}
			} else {
				m.State = ModelStateQuestionResult
				m.IncorrectAnswer = ans
			}

		case QuestionAnswerBlanks:
			if CompareAnswer(ans, a.Answers[a.Order[m.BlankIndex]]) {
				m.BlankAnswers[a.Order[m.BlankIndex]] = ans
				m.BlankIndex++
				if m.BlankIndex == len(a.Order) {
					return m, tea.Batch(cmd, m.IncrQuestion(true))
				}
			} else {
				m.State = ModelStateQuestionResult
				m.IncorrectAnswer = ans
			}
		}
	}

	return m, cmd
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

func (m *Model) ResultStateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.QuestionInput, cmd = m.QuestionInput.Update(msg)
	var ncmd tea.Cmd
	m.QuestionViewport, ncmd = m.QuestionViewport.Update(msg)
	cmd = tea.Batch(cmd, ncmd)

	if k, ok := msg.(tea.KeyMsg); ok && k.Type == tea.KeyEnter {
		ans := m.QuestionInput.Value()
		m.QuestionInput.Reset()

		switch a := m.Questions[m.Question].Answer.(type) {
		case QuestionAnswerSingle:
			m.State = ModelStateQuiz
			if strings.HasPrefix(strings.ToLower(ans), "y") {
				return m, tea.Batch(cmd, m.IncrQuestion(true))
			} else {
				return m, tea.Batch(cmd, m.IncrQuestion(false))
			}

		case QuestionAnswerMultiple:
			m.State = ModelStateQuiz
			m.IncrQuestion(false)

		case QuestionAnswerBlanks:
			m.State = ModelStateQuiz
			if strings.HasPrefix(strings.ToLower(ans), "y") {
				m.BlankAnswers[a.Order[m.BlankIndex]] = a.Answers[a.Order[m.BlankIndex]]
				m.BlankIndex++
				if m.BlankIndex == len(a.Order) {
					return m, tea.Batch(cmd, m.IncrQuestion(true))
				}
			} else {
				return m, tea.Batch(cmd, m.IncrQuestion(false))
			}
		}
	}

	return m, cmd
}

func (m *Model) ResultStateView() string {
	cont := &strings.Builder{}
	switch a := m.Questions[m.Question].Answer.(type) {
	case QuestionAnswerSingle:
		fmt.Fprintf(cont, "%s%s\n\n%s%s\n\n%s", ErrStyle.Render("Your answer: "), m.IncorrectAnswer, CorrectStyle.Render("Correct answer: "), a.Answer, MessageStyle.Render("Typo? [y/n]"))

	case QuestionAnswerMultiple:
		fmt.Fprintf(cont, "%s%s\n\n%s", ErrStyle.Render("Your answer: "), m.IncorrectAnswer, CorrectStyle.Render("Correct answers: "))
		for _, v := range a.Answers {
			alreadyfound := false
			for _, w := range m.MultipleAnswerProgress {
				if CompareAnswer(v, w) {
					alreadyfound = true
					break
				}
			}
			if alreadyfound {
				continue
			}
			fmt.Fprintf(cont, "%s\n", ListAnswerStyle.Render(v))
		}
		cont.WriteString(MessageStyle.Render("\nPress ENTER to continue..."))

	case QuestionAnswerBlanks:
		fmt.Fprintf(cont, "%s%s\n\n%s%s\n\n%s", ErrStyle.Render("Your answer: "), m.IncorrectAnswer, CorrectStyle.Render("Correct answer: "), a.Answers[a.Order[m.BlankIndex]], MessageStyle.Render("Typo? [y/n]"))
	}

	m.QuestionViewport.SetContent(m.RenderQuestion() + "\n\n" + m.QuestionInput.View() + "\n\n" + cont.String())
	return m.QuestionProgress.View() + "\n\n" + m.QuestionViewport.View()
}
