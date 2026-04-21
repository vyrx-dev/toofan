// Topic: Dependency Injection

using System;

interface IMessageService
{
    void Send(string message);
}

class ConsoleMessageService : IMessageService
{
    public void Send(string message)
    {
        Console.WriteLine(message);
    }
}

class Program
{
    static void Main()
    {
        IMessageService service = new ConsoleMessageService();
        service.Send("Hello from DI");
    }
}