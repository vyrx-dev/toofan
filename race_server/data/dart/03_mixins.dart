// Topic: Mixins
mixin Printable on Object {
    void prettyPrint() {
        print(toString());
    }
}

class Task with Printable {
    final String title;
    final bool done;

    Task(this.title, {this.done = false});

    @override
    String toString() => '[${ done ? "x" : " " }] $title';
}
