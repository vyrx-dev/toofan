// Topic: Arrow Functions and Closures
const greet = (name) => `Hello, ${name}!`;

const multiply = (a, b) => a * b;

const counter = () => {
    let count = 0;
    return {
        inc: () => ++count,
        dec: () => --count,
        get: () => count
    };
};

const c = counter();
c.inc();
c.inc();
c.dec();
console.log(c.get());

const nums = [3, 1, 4, 1, 5];
const sorted = nums.sort((a, b) => a - b);
console.log(sorted);
