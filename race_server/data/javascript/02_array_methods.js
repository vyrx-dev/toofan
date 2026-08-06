// Topic: Array Methods
const items = [
    { name: "Keyboard", price: 75, inStock: true },
    { name: "Mouse", price: 25, inStock: false },
    { name: "Monitor", price: 300, inStock: true },
    { name: "Cable", price: 10, inStock: true }
];

const available = items
    .filter(item => item.inStock)
    .map(item => item.name)
    .sort();

const total = items
    .filter(item => item.inStock)
    .reduce((sum, item) => sum + item.price, 0);

console.log(available);
console.log(`Total: $${total}`);
