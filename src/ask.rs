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

fn print_list_text(q: &String, answers: &Vec<String>) {
  progress();
  println!("{} {}", cmd::primary(q), cmd::secondary(&format!("[{} answers]", answers.len())));
}

fn ask_multiple(q: &String, answers: &Vec<String>) -> bool {
  let mut ans: Vec<String> = Vec::new();
  let mut finished: Vec<bool> = vec![false; answers.len()];
  while ans.len() < answers.len() {
    print_list_text(q, answers);
    for answer in ans.iter() {
      println!("{}", answer);
    }
    eprint!("\n{} ", cmd::secondary(&"Answer:".to_string()));

    // Check if input is correct
    let inp = cmd::input();
    let mut ind = -1;
    for (i, val) in answers.iter().enumerate() {
      if val.to_lowercase() == inp {
        ind = i as i32;
        break;
      }
    }
    if ind == -1 {
      // Wrong
      print_list_text(q, answers);
      println!("\n{}", cmd::correct(&"Remaining Answers:".to_string()));
      for (i, v) in answers.iter().enumerate() {
        if !finished[i] {
          println!("{}", v);
        }
      }
      println!("\n{} {}", cmd::wrong(&"Your Answer:".to_string()), inp);
      eprint!("\nPress ENTER to continue...");
      cmd::input();
      return false;
    }

    finished[ind as usize] = true;
    ans.push(inp);
  }
  true
}

fn ask_blanks(q: &Vec<QuestionText>, answers: &Vec<BlankAnswer>) -> bool {
  progress();
  true
}