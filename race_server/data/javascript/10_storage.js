// Topic: Local Storage
const save = (key, value) => {
    localStorage.setItem(key, JSON.stringify(value));
};

const load = (key, fallback = null) => {
    const raw = localStorage.getItem(key);
    if (raw === null) return fallback;
    try {
        return JSON.parse(raw);
    } catch {
        return fallback;
    }
};
