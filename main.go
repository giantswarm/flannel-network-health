package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/giantswarm/flannel-network-health/server"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/microerror"
	"github.com/pkg/errors"
)

var (
	description string = "Flannel-network-health serves as health endpoint for network configuration created by flannel-operator."
	gitCommit   string = "n/a"
	name        string = "flannel-network-health"
	source      string = "https://github.com/giantswarm/flannel-network-health"
)

const (
	ListenOn = ":8081"
	MaxRetry = 100
)

func main() {
	err := mainWithError()
	if err != nil {
		panic(fmt.Sprintf("%#v\n", microerror.Mask(err)))
	}
}

func mainWithError() (error){
	// for architect
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "--help") {
		println("flannel network health version 0.1")
		return nil
	}


	var err error
	// Create a new logger which is used by all packages.
	var logger micrologger.Logger
	{
		loggerConfig := micrologger.DefaultConfig()
		loggerConfig.IOWriter = os.Stdout
		logger, err = micrologger.New(loggerConfig)
		if err != nil {
			return err
		}
	}

	var flannelFile string = os.Getenv("NETWORK_ENV_FILE_PATH")
	// wait for file creation
	for count := 0; ; count++ {
		// don't wait forever, if file is not created within retry limit, exit with failure
		if count > MaxRetry {
			logger.Log(fmt.Print("After 100sec flannel file is not created. Failure"))
			return errors.New("Failed to read flannel file.")
		}
		// check if file exists
		if _, err := os.Stat(flannelFile); !os.IsNotExist(err) {
			break
		}
		logger.Log(fmt.Printf("Waiting for file %s to be created.", flannelFile))
		time.Sleep(1 * time.Second)
	}

	s := server.DefaultConfig()
	s.Logger = logger
	if !s.LoadConfig() {
		// failed to load config exiting
		return errors.New("Failed to load config from env.")
	}

	// start blocking http server
	http.HandleFunc("/bridge-healthz", s.CheckBridgeInterface)
	http.HandleFunc("/flannel-healthz", s.CheckFlannelInterface)
	err = http.ListenAndServe(ListenOn, nil)

	return err
}
