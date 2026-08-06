// Topic: Error Handling
async function fetchJSON(url) {
    try {
        const res = await fetch(url);
        if (!res.ok) {
            throw new Error(`HTTP ${res.status}`);
        }
        return await res.json();
    } catch (err) {
        console.error("fetch failed:", err.message);
        return null;
    }
}
