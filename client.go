package main

import (
        "net"
        "strconv"
	"fmt"
	"runtime"
	"sync"
	"time"
        "net/http"
	"github.com/kardianos/service"
)

const (
        macConfig     = "/opt/spac/config.json"
	linuxConfig   = "/etc/spac/config.json"
	windowsConfig = "C:/etc/spac/config.json"
)

type App struct {
	Name string `json:"name"`	
	Port string `json:"port"`	
}

type Config struct {
	Port   string `json:"port"`
	Server string `json:"server"`
        Apps   []App
}

// Client represents a client
type Client struct {
	sync.WaitGroup
	outputCh chan<- string
	stopCh   chan struct{}
        Server string
        Port  string
        Config string
	Debug bool
}

// Start implements the service.Interface interface
func (c *Client) Start(s service.Service) error {
	var config string

        if c.Config != "" {
		config = c.Config
        } else {
		switch runtime.GOOS {
        	case "darwin":
        	        config = macConfig
		case "linux":
			config = linuxConfig
		case "windows":
			config = windowsConfig
		default:
			return fmt.Errorf("Unsupported OS: %s", runtime.GOOS)
		}
        }
        if c.Debug {
	        logger.Info("Using configuration file: " + config)
	}
	// Since preparing can take some time, we continue in another
	// goroutine and return here to indicate we are started
	go c.prepare(s, config)

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := strconv.ParseUint(r.URL.Path[1:], 10, 16)

	if err != nil {
        	http.Error(w, "", 404)
	} else {
		ln, err := net.Listen("tcp4", ":" + r.URL.Path[1:])
		if err != nil {
			fmt.Fprintln(w, "")
		} else {
			ln.Close()
			http.Error(w, "", 404)
		}
	}
}

func (c *Client) prepare(s service.Service, config string) {
	// Make sure we stop the service whenever we are done
	defer func() {
		time.Sleep(1 * time.Second)

		if err := s.Stop(); err != nil {
			logger.Errorf("Unable to stop service: %v", err)
		}
	}()

	cfg, err := readConfig(config)
	if err != nil {
		logger.Error(err)
		return
	}

	if cfg == nil {
		logger.Warning("No config file found, stopping the service.")
		return
	}

        if cfg.Port == "" {
                logger.Warning("No port entry in config file found, stopping the service.")
                return
        }

	c.stopCh = make(chan struct{})

        http.HandleFunc("/", handler)
        http.ListenAndServe(cfg.Server + ":" + cfg.Port, nil)

	return
}


// Stop implements the service.Interface interface
func (c *Client) Stop(s service.Service) error {
	return nil
}

func formatOutput(output string) string {
	return fmt.Sprintf("%s [SPAC] %v\n", time.Now().Format("2006-01-02 15:04:05"), output)
}
