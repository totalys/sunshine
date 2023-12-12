package externalserver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Server a mock external server
type Server struct {
	s *httptest.Server
}

// Create a mock HTTP Server that will return a response with HTTP code and body.
func MockServer(code int, body string) *httptest.Server {
	serv := mockServerForQuery("", code, body)
	return serv.s
}

func mockServerForQuery(query string, code int, body string) *Server {
	server := &Server{}

	server.s = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprintln(w, body)
	}))

	return server
}
