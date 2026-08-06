// Topic: Modern C++ Features

#include <iostream>
using namespace std;

auto square(auto x) {
    return x * x;
}

int main() {
    cout << square(5) << endl;
}