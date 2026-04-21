// Topic: File Input/Output
 
#include <fstream>
#include <iostream>
#include <string>
using namespace std;

int main() {
    ifstream file("data.txt");
    string line;

    while (getline(file, line)) {
        cout << line << endl;
    }
}