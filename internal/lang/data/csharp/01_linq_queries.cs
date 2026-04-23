// Topic: LINQ Queries

using System;
using System.Linq;

class Program
{
    static void Main()
    {
        var numbers = new[] { 1, 2, 3, 4 };

        var squares = numbers.Select(n => n * n);

        foreach (var n in squares)
        {
            Console.WriteLine(n);
        }
    }
}