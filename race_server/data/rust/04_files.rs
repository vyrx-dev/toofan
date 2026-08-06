// Topic: File I/O and Result
use std::fs::File;
use std::io::{self, Write};

fn main() -> io::Result<()> {
    let mut file = File::create("output.txt")?;
    writeln!(file, "Hello, Rust!")?;

    Ok(())
}
