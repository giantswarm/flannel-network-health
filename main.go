package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/flannel-network-health/server"
)

var (
	description string = "Flannel-network-health serves as health endpoint for network configuration created by flannel-operator."
	gitCommit   string = "n/a"
	name        string = "flannel-network-health"
	source      string = "https://github.com/giantswarm/flannel-network-health"
)

const ListenOn  = ":8081"

func main() {
	// for architect
	if os.Args[1] == "version" {
                println("flannel network health version 0.1")
                return
        }
        if os.Args[1] == "--help" {
                println("flannel network health version 0.1")
                return
        }


	var err error
	// Create a new logger which is used by all packages.
	var logger micrologger.Logger
	{
		loggerConfig := micrologger.DefaultConfig()
		loggerConfig.IOWriter = os.Stdout
		logger, err = micrologger.New(loggerConfig)
		if err != nil {
			panic(err)
		}
	}

	var flannelFile string = os.Getenv("NETWORK_ENV_FILE_PATH")
	// wait for file creation
	for {
		if _, err := os.Stat(flannelFile); !os.IsNotExist(err) {
			break
		}
		logger.Log(fmt.Printf("Waiting for file %s to be created.", flannelFile))
		time.Sleep(1 * time.Second)
	}

	s := server.DefaultConfig()
	s.Logger = logger
	if ! s.LoadConfig() {
		// failed to load config exiting
		os.Exit(1)
	}

	// start blocking http server
	http.HandleFunc("/bridge-healthz", s.CheckBridgeInterface)
	http.HandleFunc("/flannel-healthz", s.CheckFlannelInterface)
	http.ListenAndServe(ListenOn, nil)
}


