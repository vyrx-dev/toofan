// Topic: Dynamic Memory Allocation

#include <stdio.h>
#include <stdlib.h>

int main() {
    int *arr = (int *)malloc(5 * sizeof(int));

    for (int i = 0; i < 5; i++) {
        arr[i] = i * i;
        printf("%d ", arr[i]);
    }

    free(arr);
    return 0;
}