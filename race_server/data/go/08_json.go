//go:build ignore

// Topic: JSON and Structs
type Config struct {
    Port    int      `json:"port"`
    Debug   bool     `json:"debug"`
    Origins []string `json:"origins"`
}

func loadConfig(path string) (*Config, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var cfg Config
    if err := json.NewDecoder(f).Decode(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
