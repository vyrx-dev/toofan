//go:build ignore

// Topic: Closures and Defer
func withRetry(n int, fn func() error) error {
    var last error
    for i := 0; i < n; i++ {
        if err := fn(); err != nil {
            last = err
            continue
        }
        return nil
    }
    return fmt.Errorf("failed after %d attempts: %w", n, last)
}
