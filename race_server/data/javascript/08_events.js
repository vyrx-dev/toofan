// Topic: Event Delegation
document.querySelector(".list").addEventListener("click", (e) => {
    if (e.target.matches(".item")) {
        e.target.classList.toggle("active");
        console.log("toggled:", e.target.textContent);
    }
});
