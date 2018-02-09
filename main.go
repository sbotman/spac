package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
	"github.com/kardianos/service"
)

const currentVersion = "0.1.0"

var logger service.Logger

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	version := flag.Bool("version", false, "Show version.")
	debug   := flag.Bool("debug", false, "Debug mode.")
	config  := flag.String("config", "", "Full path of configuration file.")
	flag.Parse()

	if *version {
		fmt.Println("Version: " + currentVersion)
		return
	}

	svcConfig := &service.Config{
		Name:        "spac",
		DisplayName: "Simple Port Application Checker",
		Description: "Enables https calls to checks applications by active / open ports",
	}

	c := &Client{}
        
        if *config != "" {
		strConfig := *config
		c.Config = strConfig
        }

	if *debug {
		c.Debug = true
	}

	s, err := service.New(c, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

// readConfig reads and parses the config file
func readConfig(config string) (*Config, error) {
	for i := 0; i < 60; i++ {
		if _, err := os.Stat(config); os.IsNotExist(err) {
			time.Sleep(5 * time.Second)
			continue
		}

		raw, err := ioutil.ReadFile(config)
		if err != nil {
			return nil, fmt.Errorf("Failed to read config file %q: %v", config, err)
		}

		cfg := &Config{}
		err = json.Unmarshal(raw, cfg)
		if err != nil {
			return nil, fmt.Errorf("Error parsing config file %q: %v", config, err)
		}

		return cfg, nil
	}

	return nil, nil
}
