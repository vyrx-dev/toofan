//go:build ignore

// Topic: Error Wrapping
func readConfig(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("config %s: %w", path, err)
    }
    if len(data) == 0 {
        return nil, errors.New("empty config")
    }
    return data, nil
}
