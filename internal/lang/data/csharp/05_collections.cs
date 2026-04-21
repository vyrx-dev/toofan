// Topic: Collections

using System;
using System.Collections.Generic;

class Program
{
    static void Main()
    {
        var list = new List<int> { 1, 2, 3 };

        foreach (var n in list)
        {
            Console.WriteLine(n * n);
        }
    }
}