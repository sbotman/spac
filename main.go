package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/kardianos/service"
)

const currentVersion = "0.1.0"

var logger service.Logger

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	config := flag.String("config", "", "Full path of configuration file.")
	debug := flag.Bool("debug", false, "Debug mode.")
	flags := flag.String("service", "", "Control the system service.")
	version := flag.Bool("version", false, "Show version.")
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

	c := &Service{
		Config: *config,
		Debug:  *debug,
	}

	s, err := service.New(c, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if *flags != "" {
		err = service.Control(s, *flags)
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

	if err = s.Run(); err != nil {
		logger.Error(err)
	}
}
