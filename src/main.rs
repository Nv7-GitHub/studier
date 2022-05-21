#[macro_use]
extern crate lazy_static;

mod cmd;
mod parse;
mod ask;
use std::{error, fs};
use std::sync::Mutex;

fn parse(filename: &String) -> Result<Vec<parse::Question>, Box<dyn error::Error>> {
    cmd::clear();
    println!("Parsing {}...", cmd::correct(filename));

    // Read
    let cont = fs::read_to_string(filename)?;
    let qvals: Vec<_> = cont.split("\n\n").collect();
    let mut qs: Vec<parse::Question> = Vec::with_capacity(qvals.len());
     
    // Go through
    for q in qvals {
        let mut lines = q.split("\n");
        let question = lines.next().unwrap();
        if question == "include" {
            for val in lines {
                let mut vals = parse(&val.to_string())?;
                qs.append(&mut vals);
            }
        } else {
            let lines = lines.map(|s| s.to_string()).collect();
            qs.push(parse::parse(&question.to_string(), &lines)?);
        }
    }

    Ok(qs)
}

lazy_static! {
    static ref PROGRESS_VAL: Mutex<f32> = Mutex::new(0.0);
}

fn progress() {
    cmd::clear();
    println!("{:.2}%\n", PROGRESS_VAL.lock().unwrap());
}

fn main() {
    cmd::clear();
    eprint!("{}", cmd::primary(&"Questions File: ".to_string()));
    let qs = match parse(&cmd::input()) {
        Err(e) => cmd::error(&e.to_string()),
        Ok(v) => v,
    };
    cmd::clear();

    // Ask
    let mut done_cnt = 0;
    let mut finished = vec![false; qs.len()];
    while done_cnt < finished.len() {
        for (i, val) in finished.iter_mut().enumerate() {
            if !*val {
                if ask::ask(&qs[i]) {
                    *val = true;
                    done_cnt += 1;
                    *PROGRESS_VAL.lock().unwrap() = done_cnt as f32 / qs.len() as f32 * 100.0;
                }
            }
        }
    }

    cmd::clear();
    println!("{}", cmd::correct(&"Done!".to_string()));
}
