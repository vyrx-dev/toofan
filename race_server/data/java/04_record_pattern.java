// Topic: Records and Sealed Classes

record User(String name, String email, int age) {
    User {
        if (age < 0) throw new IllegalArgumentException("age < 0");
    }
}

sealed interface Shape permits Circle, Rectangle {}
record Circle(double radius) implements Shape {}
record Rectangle(double width, double height) implements Shape {}

double area = switch (shape) {
    case Circle c -> Math.PI * c.radius() * c.radius();
    case Rectangle r -> r.width() * r.height();
};