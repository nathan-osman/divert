package main

import (
	"os"
	"os/signal"
	"syscall"
)

// Redirect stores behavior for a specific redirect.
type Redirect struct {
	Target      string
	Permanent   bool
	IncludePath bool
}

// Config stores configuration for the application.
type Config struct {
	Addr      string
	Redirects map[string]*Redirect
}

func main() {
	// The application configuration is supplied via the command line.
	// The arguments are supplied in a certain order - first the parameters
	// for a redirect (target, permanent, path, etc) and then the domain.

	var (
		target      string
		permanent   bool
		includePath bool
		cfg         = &Config{
			Addr:      ":80",
			Redirects: map[string]*Redirect{},
		}
	)
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--addr", "-a":
			i++
			cfg.Addr = os.Args[i]
		case "--target", "-t":
			i++
			target = os.Args[i]
		case "--permanent", "-p":
			permanent = true
		case "--include-path", "-i":
			includePath = true
		default:
			if len(target) == 0 {
				panic("target must be specified")
			}
			cfg.Redirects[os.Args[i]] = &Redirect{
				Target:      target,
				Permanent:   permanent,
				IncludePath: includePath,
			}
			target = ""
			permanent = false
			includePath = false
		}
	}

	s, err := NewServer(cfg)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	// Wait for SIGINT or SIGTERM
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
