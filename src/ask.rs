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
  println!("{}\n", cmd::primary(q));
  eprint!("{}", cmd::secondary(&"Answer: ".to_string()));
  let v = cmd::input();
  if v.to_lowercase() == answer.to_lowercase() {
    return true;
  }

  // Wrong
  progress();
  println!("{}\n", cmd::primary(q));
  println!("{}{}", cmd::correct(&"Correct Answer: ".to_string()), answer);
  eprint!("{}", cmd::wrong(&"Typo? [y/n]: ".to_string()));

  cmd::input() == "y"
}

fn ask_multiple(q: &String, answers: &Vec<String>) -> bool {
  progress();
  true
}

fn ask_blanks(q: &Vec<QuestionText>, answers: &Vec<BlankAnswer>) -> bool {
  progress();
  true
}