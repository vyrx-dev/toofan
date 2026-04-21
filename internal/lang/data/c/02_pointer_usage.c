// Topic: Pointer Usage

#include <stdio.h>

void increment(int *x) {
    (*x)++;
}

int main() {
    int num = 5;
    increment(&num);
    printf("%d\n", num);
    return 0;
}