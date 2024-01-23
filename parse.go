package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type QuestionTextKind int

const (
	QuestionTextKindText QuestionTextKind = iota
	QuestionTextKindBlank
)

type QuestionText struct {
	Kind  QuestionTextKind
	Value string
}

type QuestionAnswer interface{}

type QuestionAnswerSingle struct {
	Answer string
}

type QuestionAnswerMultiple struct {
	Answers []string
}

type QuestionAnswerBlanks struct {
	Answers map[string]string
}

type Question struct {
	Text   []QuestionText
	Answer QuestionAnswer
}

type ParseError struct {
	Line    int
	File    string
	Message string
}

func (p ParseError) Error() string {
	return ErrStyle.Render(fmt.Sprintf("%s:%d: ", p.File, p.Line)) + p.Message
}

func NewParseError(line int, file, msg string) ParseError {
	return ParseError{
		Line:    line,
		File:    file,
		Message: msg,
	}
}

func (m *Model) ParseQuestion(inp string, linenum int, file string) (*Question, error) {
	lines := strings.Split(inp, "\n")

	// Parse
	curr := ""
	inblank := false
	hasblank := false
	text := make([]QuestionText, 0)
	for _, ch := range lines[0] {
		switch ch {
		case '`':
			kind := QuestionTextKindText
			if inblank {
				kind = QuestionTextKindBlank
			}
			text = append(text, QuestionText{
				Kind:  kind,
				Value: curr,
			})
			curr = ""
			inblank = !inblank
			hasblank = true
		default:
			curr += string(ch)
		}
	}
	if inblank {
		return nil, NewParseError(linenum, file, "Unclosed blank")
	}
	text = append(text, QuestionText{
		Kind:  QuestionTextKindText,
		Value: curr,
	})

	// Do answers
	var ans QuestionAnswer
	if len(lines) == 1 {
		return nil, NewParseError(linenum, file, "No answers found")
	}
	if hasblank {
		res := make(map[string]string)
		for i, ans := range lines[1:] {
			parts := strings.SplitN(ans, ":", 2)
			if len(parts) == 1 || len(parts) > 2 {
				return nil, NewParseError(linenum+i, file, "Improper blank answer")
			}
			res[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
		ans = QuestionAnswerBlanks{
			Answers: res,
		}
	} else {
		if len(lines) > 2 {
			ans = QuestionAnswerMultiple{Answers: lines[1:]}
		} else {
			ans = QuestionAnswerSingle{Answer: lines[1]}
		}
	}
	return &Question{
		Text:   text,
		Answer: ans,
	}, nil
}

func (m *Model) ParseFile(path string) ([]Question, tea.Cmd) {
	// Read file
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, m.HandleErr(err)
	}
	qs := strings.SplitN(string(file), "\n\n", 2)
	questions := make([]Question, 0, len(qs))
	for _, q := range qs {
		if strings.HasPrefix(strings.TrimSpace(q), "include") {
			q := strings.Split(q, "\n")
			if len(q) == 1 {
				continue
			}
			for _, v := range q[1:] {
				path := filepath.Join(filepath.Dir(path), v)
				vals, cmd := m.ParseFile(strings.TrimSpace(path))
				if cmd != nil {
					return nil, cmd
				}
				questions = append(questions, vals...)
			}
			continue
		}

		question, err := m.ParseQuestion(q, len(questions)+1, path)
		if err != nil {
			return nil, m.HandleErr(err)
		}
		questions = append(questions, *question)
	}

	return questions, nil
}

func (m *Model) InputParseState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.FilePicker, cmd = m.FilePicker.Update(msg)

	if didSelect, path := m.FilePicker.DidSelectFile(msg); didSelect {
		var pCmd tea.Cmd
		m.Questions, pCmd = m.ParseFile(path)
		if pCmd != nil {
			return m, pCmd
		}

		m.State = ModelStateQuiz
		m.FileID = filepath.Base(path)
		m.UpdateProgress()
		return m, cmd
	}
	return m, cmd
}

func (m *Model) InputParseStateView() string {
	return "Select a file to study.\n\n" + m.FilePicker.View()
}
