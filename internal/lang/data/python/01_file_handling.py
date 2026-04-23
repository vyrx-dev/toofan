# Topic: File Handling

with open("data.txt", "r") as file:
    for line in file:
        print(line.strip())