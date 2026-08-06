// Topic: Destructuring and Spread
const config = {
    host: "localhost",
    port: 3000,
    debug: true
};

const { host, port, ...rest } = config;
console.log(`${host}:${port}`);

const defaults = { theme: "dark", lang: "en" };
const settings = { ...defaults, ...rest };

const nums = [1, 2, 3];
const more = [...nums, 4, 5];

function sum(...args) {
    return args.reduce((a, b) => a + b, 0);
}

console.log(sum(...more));
