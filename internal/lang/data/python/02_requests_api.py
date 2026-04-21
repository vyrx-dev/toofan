# Topic: HTTP Requests (API)

import urllib.request
import json

def fetch_data(url: str):
    try:
        with urllib.request.urlopen(url, timeout=5) as response:
            data = response.read()
            return json.loads(data)
    except Exception as e:
        print(f"Request failed: {e}")
        return None

def main():
    url = "https://api.example.com/data"
    data = fetch_data(url)

    if data:
        print("Response received:")
        print(data)

if __name__ == "__main__":
    main()