// Topic: Struct Usage

#include <stdio.h>

struct Person {
    char name[50];
    int age;
};

int main() {
    struct Person p = {"Alice", 25};
    printf("%s %d\n", p.name, p.age);
    return 0;
}