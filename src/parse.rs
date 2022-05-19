use std::error::{Error};
use thiserror::Error;
use std::collections::HashSet;

pub enum QuestionText {
  Text(String),
  Blank(String),
}

pub struct BlankAnswer {
  text: String,
  answer: String,
}

pub enum QuestionAnswer {
  Answer(String),
  Answers(Vec<String>),
  Blanks(Vec<BlankAnswer>),
}

pub struct Question {
  text: Vec<QuestionText>,
  answer: QuestionAnswer,
}

#[derive(Debug, Error)]
enum ParseError {
  #[error("un-closed blank for question {0}")]
  UnClosedBlank(String),
  #[error("missing blank value: {1} for question {0}")]
  MissingBlankValue(String, String),
  #[error("unneded blank value: {1} for question {0}")]
  UnneededBlankValue(String, String),
}

pub fn parse(q: &String, answers: &Vec<String>) -> Result<Question, impl Error> {
  let mut curr = String::new();
  let mut isBlank = false;
  let mut hasBlank = false;
  let mut text: Vec<QuestionText> = Vec::new();

  for c in q.chars() {
    match c {
      '`' => {
        if isBlank {
          text.push(QuestionText::Blank(curr.clone()))
        } else {
          text.push(QuestionText::Text(curr.clone()))
        }
        curr = String::new();
        isBlank = !isBlank;
        hasBlank = true;
      }
      _ => curr.push(c)
    }
  }

  if curr != "" {
    if isBlank {
      text.push(QuestionText::Blank(curr.clone()))
    } else {
      text.push(QuestionText::Text(curr.clone()))
    }
  }

  if isBlank {
    return Err(ParseError::UnClosedBlank(q.clone()))
  }

  // Parse answers
  if hasBlank {
    // Calculate has
    let mut blanks: Vec<BlankAnswer> = Vec::with_capacity(answers.len());
    let mut has = HashSet::new();
    for val in answers {
      let parts: Vec<_> = val.splitn(2, ": ").collect();
      has.insert(parts[0].to_string());
      blanks.push(BlankAnswer{text: parts[0].to_string(), answer: parts[1].to_string()});
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

    return Ok(Question{
      text: text,
      answer: QuestionAnswer::Blanks(blanks),
    })
  }

  if answers.len() == 1 {
    return Ok(Question {
      text: text,
      answer: QuestionAnswer::Answer(answers[0].clone())
    })
  }
  Ok(Question {
    text: text,
    answer: QuestionAnswer::Answers(answers.clone())
  })
}