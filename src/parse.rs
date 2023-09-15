use enum_as_inner::EnumAsInner;
use std::collections::HashSet;
use std::error::Error;
use std::fmt::Write;
use thiserror::Error;

#[derive(EnumAsInner)]
pub enum QuestionText {
    Text(String),
    Blank(String),
}

pub struct BlankAnswer {
    pub text: String,
    pub answer: String,
}

pub enum QuestionAnswer {
    Answer(String),
    Answers(Vec<String>),
    Blanks(Vec<BlankAnswer>),
}

pub struct Question {
    pub text: Vec<QuestionText>,
    pub answer: QuestionAnswer,
}

impl Question {
    pub fn rawtext(&self) -> String {
        let mut out = String::new();
        for v in self.text.iter() {
            match v {
                QuestionText::Blank(t) | QuestionText::Text(t) => out.push_str(t.as_str()),
            }
        }
        out.push('\n');
        match &self.answer {
            QuestionAnswer::Answer(v) => {
                out.push_str(v.as_str());
                out.push('\n')
            }
            QuestionAnswer::Answers(vals) => {
                for v in vals {
                    out.push_str(v.as_str());
                    out.push('\n');
                }
            }
            QuestionAnswer::Blanks(ans) => {
                for v in ans {
                    write!(out, "{}: {}", v.text, v.answer).unwrap()
                }
            }
        }

        out
    }
}

#[derive(Debug, Error)]
enum ParseError {
    #[error("un-closed blank for question {0}")]
    UnClosedBlank(String),
    #[error("missing blank value: {1} for question {0}")]
    MissingBlankValue(String, String),
    #[error("unneded blank value: {1} for question {0}")]
    UnneededBlankValue(String, String),
    #[error("improper blank for question {0}")]
    ImproperBlank(String),
}

pub fn parse(q: &String, answers: &Vec<String>) -> Result<Question, impl Error> {
    let mut curr = String::new();
    let mut is_blank = false;
    let mut has_blank = false;
    let mut text: Vec<QuestionText> = Vec::new();

    for c in q.chars() {
        match c {
            '`' => {
                if is_blank {
                    text.push(QuestionText::Blank(curr.clone()))
                } else {
                    text.push(QuestionText::Text(curr.clone()))
                }
                curr = String::new();
                is_blank = !is_blank;
                has_blank = true;
            }
            _ => curr.push(c),
        }
    }

    if curr != "" {
        if is_blank {
            text.push(QuestionText::Blank(curr.clone()))
        } else {
            text.push(QuestionText::Text(curr.clone()))
        }
    }

    if is_blank {
        return Err(ParseError::UnClosedBlank(q.clone()));
    }

    // Parse answers
    if has_blank {
        // Calculate has
        let mut blanks: Vec<BlankAnswer> = Vec::with_capacity(answers.len());
        let mut has = HashSet::new();
        for val in answers {
            let parts: Vec<_> = val.splitn(2, ": ").collect();
            if parts.len() < 2 {
                return Err(ParseError::ImproperBlank(q.clone()));
            }
            has.insert(parts[0].to_string());
            blanks.push(BlankAnswer {
                text: parts[0].to_string(),
                answer: parts[1].to_string(),
            });
        }

        // Calculate needed
        let mut needed = HashSet::new();
        for val in text.iter() {
            if let QuestionText::Blank(label) = val {
                needed.insert(label.clone());
            }
        }

        // Check for errors
        for val in needed.iter() {
            if !has.contains(val) {
                return Err(ParseError::MissingBlankValue(q.clone(), val.clone()));
            }
        }
        for val in has.iter() {
            if !needed.contains(val) {
                return Err(ParseError::UnneededBlankValue(q.clone(), val.clone()));
            }
        }

        return Ok(Question {
            text: text,
            answer: QuestionAnswer::Blanks(blanks),
        });
    }

    if answers.len() == 1 {
        return Ok(Question {
            text: text,
            answer: QuestionAnswer::Answer(answers[0].clone()),
        });
    }
    Ok(Question {
        text: text,
        answer: QuestionAnswer::Answers(answers.clone()),
    })
}
