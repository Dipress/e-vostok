package main

// SMTPServer struct
type SMTPServer struct {
	Host string
	Port string
}

// ServerName method returns host + port string
func (s *SMTPServer) ServerName() string {
	return s.Host + ":" + s.Port
}
