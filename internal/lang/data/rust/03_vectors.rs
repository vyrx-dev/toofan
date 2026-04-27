// Topic: Vectors and Iterators
fn main() {
    let mut vec = vec![1, 2, 3, 4];
    vec.push(5);
    vec.pop();

    for num in &vec {
        println!("Number: {}", num);
    }

    let doubled: Vec<i32> = vec.iter()
        .map(|x| x * 2)
        .collect();

    println!("Doubled: {:?}", doubled);
}
