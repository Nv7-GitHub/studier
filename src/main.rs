mod cmd;
mod parse;
use std::{error, fs};

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

fn main() {
    cmd::clear();
    eprint!("{}", cmd::primary(&"Questions File: ".to_string()));
    let qs = match parse(&cmd::input()) {
        Err(e) => cmd::error(&e.to_string()),
        Ok(v) => v,
    };
    cmd::clear();
}
