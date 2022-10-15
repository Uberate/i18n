fn main() {
    use std::env::Args;
    let args: Args = std::env::args();

    println!("{:?}", args);
    for arg in args {
        println!("{:?}", arg);
    }
}