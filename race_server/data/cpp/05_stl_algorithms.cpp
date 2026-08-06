// Topic: STL Algorithms

#include <vector>
#include <algorithm>
#include <iostream>
using namespace std;

int main() {
    vector<int> nums = {4, 2, 3, 1};

    sort(nums.begin(), nums.end());

    for (int n : nums) {
        cout << n << " ";
    }
}