//go:build ignore

// Topic: Channels and Select
ch := make(chan string, 1)
timeout := time.After(2 * time.Second)

go func() {
    time.Sleep(time.Second)
    ch <- "done"
}()

select {
case msg := <-ch:
    fmt.Println("received:", msg)
case <-timeout:
    fmt.Println("timed out")
}
