// Topic: Structs and Impl Blocks
struct User {
    name: String,
    age: u8,
}

impl User {
    fn new(name: &str, age: u8) -> Self {
        User {
            name: name.to_string(),
            age,
        }
    }
}

fn main() {
    let user = User::new("Alice", 30);
    println!("{} is {} years old", user.name, user.age);
}
