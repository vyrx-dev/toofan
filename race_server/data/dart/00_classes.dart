// Topic: Classes and Null Safety
class User {
  final String name;
  final int age;
  String? bio;

  User({required this.name, required this.age, this.bio});

  String greet() {
    final desc = bio ?? 'no bio';
    return 'Hi, I am $name ($age) - $desc';
  }
}

void main() {
  final user = User(name: 'Alice', age: 25);
  print(user.greet());

  final admin = User(
    name: 'Bob',
    age: 30,
    bio: 'System admin',
  );
  print(admin.greet());
}
