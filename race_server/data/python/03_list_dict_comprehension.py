# Topic: List & Dictionary Comprehension

nums = [1, 2, 3, 4]

squares = [n**2 for n in nums]

square_map = {n: n**2 for n in nums}

print(squares)
print(square_map)