// Topic: Async Error Handling

async function fetchWithTimeout<T>(
    url: string,
    ms: number
): Promise<T> {
    const ctrl = new AbortController();
    const id = setTimeout(() => ctrl.abort(), ms);
    try {
        const res = await fetch(url, {
            signal: ctrl.signal,
        });
        return await res.json();
    } finally {
        clearTimeout(id);
    }
}