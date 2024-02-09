package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	Order   []string
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
	return fmt.Sprintf("%s:%d: ", p.File, p.Line) + p.Message
}

func NewParseError(line int, file, msg string) ParseError {
	return ParseError{
		Line:    line,
		File:    file,
		Message: msg,
	}
}

func ParseQuestion(inp string, linenum int, file string) (*Question, error) {
	lines := strings.Split(inp, "\n")

	// Parse
	curr := ""
	inblank := false
	hasblank := false
	text := make([]QuestionText, 0)
	blankOrder := make([]string, 0)
	for _, ch := range lines[0] {
		switch ch {
		case '`':
			kind := QuestionTextKindText
			if inblank {
				kind = QuestionTextKindBlank
				blankOrder = append(blankOrder, curr)
			}
			text = append(text, QuestionText{
				Kind:  kind,
				Value: curr,
			})
			curr = ""
			inblank = !inblank
			hasblank = true
		case ' ':
			if inblank {
				curr += string(ch)
			} else {
				text = append(text, QuestionText{
					Kind:  QuestionTextKindText,
					Value: curr + " ",
				})
				curr = ""
			}
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
			Order:   blankOrder,
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

func ParseFile(path string) []Question {
	// Read file
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	qs := strings.Split(strings.ReplaceAll(string(file), "\r", ""), "\n\n")
	questions := make([]Question, 0, len(qs))
	for _, q := range qs {
		if strings.HasPrefix(strings.TrimSpace(q), "include") {
			q := strings.Split(q, "\n")
			if len(q) == 1 {
				continue
			}
			for _, v := range q[1:] {
				path := filepath.Join(filepath.Dir(path), v)
				vals := ParseFile(strings.TrimSpace(path))
				questions = append(questions, vals...)
			}
			continue
		}

		question, err := ParseQuestion(q, len(questions)+1, path)
		if err != nil {
			panic(err)
		}
		questions = append(questions, *question)
	}

	return questions
}
