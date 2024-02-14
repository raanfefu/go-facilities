package web

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/raanfefu/go-facilities/pkg/common"
)

func (s *impl) Init() {
	if s.Params.Mode == nil {
		log.Println("Error not configuration loaded")
		os.Exit(-1)
	}
	mode := *s.Params.Mode
	s.Router = mux.NewRouter()
	switch mode {
	case common.TLS:
		if s.Params.Port == 0 {
			s.Params.Port = 443
		}
		certificate := s.Params.Certs
		s.Server = &http.Server{
			Handler:   s.Router,
			Addr:      fmt.Sprintf(":%v", s.Params.Port),
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{certificate}},
		}
	case common.HTTP:
		if s.Params.Port == 0 {
			s.Params.Port = 80
		}
		s.Server = &http.Server{
			Handler: s.Router,
			Addr:    fmt.Sprintf(":%v", s.Params.Port),
		}
	}
	log.Println("Initializing server... Done ✓")
}

func (s *impl) Start() {
	if s.Params.Mode == nil {
		log.Println("Error not configuration loaded")
		os.Exit(-1)
	}
	mode := *s.Params.Mode
	if mode == common.TLS {
		go func() {
			if err := s.Server.ListenAndServeTLS("", ""); err != nil {
				fmt.Printf("Failed to listen and serve webhook server: %v\n", err)
				os.Exit(-1)
			}
		}()
	} else {
		go func() {
			if err := s.Server.ListenAndServe(); err != nil {
				fmt.Printf("Failed to listen and serve webhook server: %v\n", err)
				os.Exit(-1)
			}
		}()
	}
	log.Printf("Starting web server 📡... Done ✓")
	log.Printf("Listen on %s://0.0.0.0:%v/", s.Params.Mode.String(), s.Params.Port)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	fmt.Println("Got shutdown signal, shutting down webhook server gracefully...\n")
	s.Server.Shutdown(context.Background())
}

func (s *impl) GetRouter() *mux.Router {
	return s.Router
}
