use super::parse::*;
use super::*;

pub fn ask(q: &Question) -> bool {
  match &q.answer {
    QuestionAnswer::Answer(ans) => ask_normal(q.text[0].as_text().unwrap(), &ans),
    QuestionAnswer::Answers(ans) => ask_multiple(q.text[0].as_text().unwrap(), &ans),
    QuestionAnswer::Blanks(blks) => ask_blanks(&q.text, &blks),
  }
}

fn ask_normal(q: &String, answer: &String) -> bool {
  progress();
  true
}

fn ask_multiple(q: &String, answers: &Vec<String>) -> bool {
  progress();
  true
}

fn ask_blanks(q: &Vec<QuestionText>, answers: &Vec<BlankAnswer>) -> bool {
  progress();
  true
}