// Topic: Multithreading

#include <iostream>
#include <thread>
using namespace std;

void task() {
    cout << "Running in thread\n";
}

int main() {
    thread t(task);
    t.join();
}