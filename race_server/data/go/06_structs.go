//go:build ignore

// Topic: Struct Methods
type Server struct {
    host string
    port int
    tls  bool
}

func NewServer(host string, port int) *Server {
    return &Server{host: host, port: port}
}

func (s *Server) Addr() string {
    scheme := "http"
    if s.tls {
        scheme = "https"
    }
    return fmt.Sprintf("%s://%s:%d", scheme, s.host, s.port)
}
