// Topic: Collections and Iteration
void main() {
  final scores = <String, int>{
    'alice': 92,
    'bob': 85,
    'eve': 97,
  };

  scores.forEach((name, score) {
    print('$name: $score');
  });

  final passed = scores.entries
      .where((e) => e.value >= 90)
      .map((e) => e.key)
      .toList();

  print('Passed: $passed');

  final nums = [3, 1, 4, 1, 5, 9];
  nums.sort();
  final unique = nums.toSet().toList();
  print(unique);
}
