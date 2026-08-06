// Topic: Closures and Iterator Chaining
fn main() {
    let numbers = vec![1, 2, 3, 4, 5, 6, 7];
    
    let even_squares: Vec<i32> = numbers
        .into_iter()
        .filter(|&x| x % 2 == 0)
        .map(|x| x * x)
        .collect();

    println!("Even squares: {:?}", even_squares);
}
