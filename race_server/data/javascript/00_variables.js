// Topic: Variables and Functions
const appName = "Toofan";
let userScore = 0;

function incrementScore(points) {
    userScore += points;
    console.log(`Score: ${userScore}`);
}

for (let i = 0; i < 3; i++) {
    incrementScore(10);
}

if (userScore >= 30) {
    console.log("You win!");
} else {
    console.log("Try again!");
}
