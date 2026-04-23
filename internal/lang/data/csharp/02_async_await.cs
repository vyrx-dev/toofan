// Topic: Async/Await

using System;
using System.Net.Http;
using System.Threading.Tasks;

class Program
{
    static async Task Main()
    {
        using var client = new HttpClient();
        var response = await client.GetStringAsync("https://api.example.com");

        Console.WriteLine(response);
    }
}