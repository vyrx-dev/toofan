// Topic: File Handling

using System;
using System.IO;

class Program
{
    static void Main()
    {
        var content = File.ReadAllText("data.txt");
        Console.WriteLine(content);
    }
}