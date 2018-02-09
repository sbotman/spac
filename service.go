package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gorilla/mux"
	"github.com/kardianos/service"
)

const (
	macConfig     = "/opt/spac/config.json"
	linuxConfig   = "/etc/spac/config.json"
	windowsConfig = "C:/etc/spac/config.json"
)

// Config represents the SPAC config.
type Config struct {
	Server string            `json:"server"`
	Port   int               `json:"port"`
	Apps   map[string]string `json:"apps"`
}

// Service represents the service itself.
type Service struct {
	Config string
	Debug  bool
}

// Start implements the service.Interface interface
func (s *Service) Start(svc service.Service) error {
	if s.Config == "" {
		switch runtime.GOOS {
		case "darwin":
			s.Config = macConfig
		case "linux":
			s.Config = linuxConfig
		case "windows":
			s.Config = windowsConfig
		default:
			return fmt.Errorf("Unsupported OS: %s", runtime.GOOS)
		}
	}

	if s.Debug {
		logger.Info("Using configuration file: " + s.Config)
	}

	cfg, err := s.readConfig()
	if err != nil {
		return fmt.Errorf("Error reading config file %q: %v", s.Config, err)
	}

	if cfg.Port == 0 {
		return fmt.Errorf("No port entry in config file")
	}

	go func() {
		r := mux.NewRouter()
		r.HandleFunc("/services/{service}", s.handler(cfg))
		http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Server, cfg.Port), nil)
	}()

	return nil
}

// readConfig reads and parses the config file
func (s *Service) readConfig() (*Config, error) {
	if _, err := os.Stat(s.Config); os.IsNotExist(err) {
		return nil, err
	}

	raw, err := ioutil.ReadFile(s.Config)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = json.Unmarshal(raw, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (s *Service) handler(cfg *Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		bind, ok := cfg.Apps[vars["service"]]
		if !ok {
			http.Error(w, "", 404)
			return
		}

		l, err := net.Listen("tcp4", bind)
		if err != nil {
			fmt.Fprintln(w, "")
			return
		}

		l.Close()
		http.Error(w, "", 404)
	}
}

// Stop implements the service.Interface interface
func (s *Service) Stop(svc service.Service) error {
	return nil
}

func formatOutput(output string) string {
	return fmt.Sprintf("%s [SPAC] %v\n", time.Now().Format("2006-01-02 15:04:05"), output)
}
