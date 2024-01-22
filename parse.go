package main

import (
	"fmt"
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
	Answers map[string]string
}

type Question struct {
	Text   []QuestionText
	Answer QuestionAnswer
}

func (m *Model) ParseQuestion(inp string, linenum int, file string) *Question {
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
		m.HandleErr(fmt.Sprintf("%s:%d: Unclosed blank", file, linenum))
		return nil
	}
	text = append(text, QuestionText{
		Kind:  QuestionTextKindText,
		Value: curr,
	})

	// Do answers
	var ans QuestionAnswer
	if len(lines) == 1 {
		m.HandleErr(fmt.Sprintf("%s:%d: No answers found", file, linenum))
		return nil
	}
	if hasblank {
		res := make(map[string]string)
		for i, ans := range lines[1:] {
			parts := strings.SplitN(ans, ":", 2)
			if len(parts) == 1 || len(parts) > 2 {
				m.HandleErr(fmt.Sprintf("%s:%d: Improper blank answer", file, linenum+i))
				return nil
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
	}
}
