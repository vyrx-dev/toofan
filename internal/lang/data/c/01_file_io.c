// Topic: File Input/Output

#include <stdio.h>

int main() {
    FILE *file = fopen("data.txt", "r");
    char line[256];

    while (fgets(line, sizeof(line), file)) {
        printf("%s", line);
    }

    fclose(file);
    return 0;
}