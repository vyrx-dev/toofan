// Topic: Async/Await and Fetch
async function fetchUser(id) {
    const res = await fetch(`/api/users/${id}`);
    if (!res.ok) {
        throw new Error(`HTTP ${res.status}`);
    }
    return await res.json();
}

async function main() {
    try {
        const user = await fetchUser(123);
        console.log(user.name);
    } catch (err) {
        console.error("Failed:", err.message);
    }
}

main();
