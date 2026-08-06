//go:build ignore

// Topic: Interfaces
type Writer interface {
    Write(data []byte) (int, error)
    Close() error
}

type FileWriter struct {
    path string
    buf  bytes.Buffer
}

func (f *FileWriter) Write(data []byte) (int, error) {
    return f.buf.Write(data)
}

func (f *FileWriter) Close() error {
    return os.WriteFile(f.path, f.buf.Bytes(), 0644)
}
