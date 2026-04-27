struct User {
    name: String,
    age: u8,
}

fn main() {
    let user = User {
        name: String::from("Alice"),
        age: 30,
    };

    println!("{} is {} years old", user.name, user.age);
}
