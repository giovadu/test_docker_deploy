package app

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Router       *mux.Router
	Addr         string
	TLSConfig    *tls.Config
	TLSNextProto map[string]func(*http.Server, *tls.Conn, http.Handler)
}

func (s *Server) Initialize(r *mux.Router) {
	s.Router = r
	s.TLSConfig = &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	s.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
}

func (s *Server) Run() {
	http.Handle("/", Cors.Handler(s.Router))

	addrr := ":2020"

	log.Println("Server running on", addrr)

	http.ListenAndServe(addrr, nil)
}

func NewServer() Server {
	return Server{}
}
