package main

import (
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Server listens for incoming requests and processes them.
type Server struct {
	cfg      *Config
	listener net.Listener
	log      *logrus.Entry
	stopped  chan bool
}

// NewServer creates a new server instance.
func NewServer(cfg *Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	var (
		s = &Server{
			cfg:      cfg,
			listener: l,
			log:      logrus.WithField("context", "server"),
			stopped:  make(chan bool),
		}
		server = http.Server{
			Handler: s,
		}
	)
	go func() {
		defer close(s.stopped)
		defer s.log.Info("server stopped")
		s.log.Info("server started")
		if err := server.Serve(l); err != nil {
			s.log.Error(err)
		}
	}()
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	re, ok := s.cfg.Redirects[r.Host]
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if re.IncludePath {
		w.Header().Set("Location", re.Target+r.URL.Path)
	} else {
		w.Header().Set("Location", re.Target)
	}
	if re.Permanent {
		w.WriteHeader(http.StatusMovedPermanently)
	} else {
		w.WriteHeader(http.StatusFound)
	}
}

// Close shuts down the server.
func (s *Server) Close() {
	s.listener.Close()
	<-s.stopped
}
