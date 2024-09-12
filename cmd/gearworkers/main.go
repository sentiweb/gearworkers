package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/sentiweb/gearworkers/pkg/config"
	"github.com/sentiweb/gearworkers/pkg/server"
	"github.com/sentiweb/gearworkers/pkg/worker"
	"gopkg.in/yaml.v3"
)

var version string
var commit string
var date string

func loadConfig(file string) *config.AppConfig {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Unable to read config in %s : %s", file, err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("Unable to read config in %s : %s", file, err)
	}
	cfg := config.AppConfig{}
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		log.Fatalf("Unable to parse config in %s : %s", file, err)
	}
	return &cfg
}

func main() {
	fmt.Printf("Gearwokers %s commit:%s (%s)\n", version, commit, date)

	configPtr := flag.String("config", "config.yaml", "Configuration file path")

	flag.Parse()

	cfg := loadConfig(*configPtr)

	manager := worker.NewManager(cfg)

	err := manager.Start()

	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("Started workers using %d goroutines\n", runtime.NumGoroutine())

	if cfg.Server.Addr != "" {
		srv := server.NewHttpServer(cfg)
		go func() {
			err := srv.Start()
			if err != nil {
				log.Printf("Http Server error : %s", err)
			}
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	sig := <-c
	log.Println("receiving signal, stopping services ", sig)

}
