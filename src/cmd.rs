use std::{fmt, process::exit};

pub fn clear() {
  eprint!("\x1b[H\x1b[2J");
}

pub fn wrong(txt: &String) -> String {
  format!("\x1b[31;1m{}\x1b[0m", txt)
}

pub fn correct(txt: &String) -> String {
  format!("\x1b[32;1m{}\x1b[0m", txt)
}

pub fn primary(txt: &String) -> String {
  format!("\x1b[34;1m{}\x1b[0m", txt)
}

pub fn secondary(txt: &String) -> String {
  format!("\x1b[33;1m{}\x1b[0m", txt)
}

pub fn ternary(txt: &String) -> String {
  format!("\x1b[35;1m{}\x1b[0m", txt)
}

pub fn input() -> String {
  let mut line = String::new();
  std::io::stdin().read_line(&mut line).unwrap();
  line.pop();  // removes newline
  line
}

pub fn error(msg: &String) -> ! {
  println!("{}", wrong(&msg.to_string()));
  exit(1);
}