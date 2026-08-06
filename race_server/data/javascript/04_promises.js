// Topic: Promises
function loadData(url) {
    return new Promise((resolve, reject) => {
        setTimeout(() => {
            if (url) {
                resolve({ data: "loaded" });
            } else {
                reject(new Error("No URL"));
            }
        }, 500);
    });
}

loadData("/api/data")
    .then(res => {
        console.log(res.data);
        return loadData("/api/next");
    })
    .catch(err => {
        console.error(err.message);
    })
    .finally(() => {
        console.log("Done");
    });
