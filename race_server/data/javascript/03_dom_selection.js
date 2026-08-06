// Topic: DOM and Events
const form = document.querySelector('#login-form');
const input = document.getElementById('username');
const msg = document.querySelector('.status');

form.addEventListener('submit', (e) => {
    e.preventDefault();
    const value = input.value.trim();

    if (!value) {
        msg.textContent = "Username required";
        msg.classList.add('error');
        return;
    }

    msg.textContent = `Welcome, ${value}!`;
    msg.classList.remove('error');
});
