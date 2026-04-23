# Topic: Data Processing

data = [
    {"name": "Alice", "age": 25},
    {"name": "Bob", "age": 30}
]

names = [person["name"] for person in data]

print(names)