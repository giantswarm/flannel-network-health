package main

import (
	"fmt"
	"github.com/giantswarm/flannel-network-health/flag"
	"github.com/giantswarm/flannel-network-health/server"
	"github.com/giantswarm/flannel-network-health/service"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/microkit/command"
	microserver "github.com/giantswarm/microkit/server"
	"github.com/giantswarm/microkit/transaction"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/microstorage"
	"github.com/giantswarm/microstorage/memory"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	f           *flag.Flag = flag.New()
	description string     = "Flannel-network-health serves as health endpoint for network configuration created by flannel-operator."
	gitCommit   string     = "n/a"
	name        string     = "flannel-network-health"
	source      string     = "https://github.com/giantswarm/flannel-network-health"
)

const (
	MaxRetry       = 100
	ListenPortBase = 21000
)

func main() {
	err := mainWithError()
	if err != nil {
		panic(fmt.Sprintf("%#v\n", microerror.Mask(err)))
	}
}

func mainWithError() error {
	var err error
	// Create a new logger which is used by all packages.
	var newLogger micrologger.Logger
	{
		loggerConfig := micrologger.DefaultConfig()
		loggerConfig.IOWriter = os.Stdout
		newLogger, err = micrologger.New(loggerConfig)
		if err != nil {
			return err
		}
	}

	// We define a server factory to create the custom server once all command
	// line flags are parsed and all microservice configuration is storted out.
	newServerFactory := func(v *viper.Viper) microserver.Server {
		// wait for flannel file to be created
		err = waitForFlannelFile(newLogger)
		if err != nil {
			panic(err)
		}
		// Create a new custom service which implements business logic.
		var newService *service.Service
		{
			serviceConfig := service.DefaultConfig()

			serviceConfig.Flag = f
			serviceConfig.Logger = newLogger

			serviceConfig.Description = description
			serviceConfig.GitCommit = gitCommit
			serviceConfig.Name = name
			serviceConfig.Source = source

			newService, err = service.New(serviceConfig)
			if err != nil {
				panic(err)
			}
		}

		var storage microstorage.Storage
		{
			storage, err = memory.New(memory.DefaultConfig())
			if err != nil {
				panic(err)
			}
		}

		var transactionResponder transaction.Responder
		{
			c := transaction.DefaultResponderConfig()
			c.Logger = newLogger
			c.Storage = storage

			transactionResponder, err = transaction.NewResponder(c)
			if err != nil {
				panic(err)
			}
		}

		var listenAddress string
		{
			listenAddress, err = buildListenAddress()
			if err != nil {
				panic(err)
			}
		}

		// Create a new custom server which bundles our endpoints.
		var newServer microserver.Server
		{
			serverConfig := server.DefaultConfig()

			serverConfig.MicroServerConfig.Logger = newLogger
			serverConfig.MicroServerConfig.ServiceName = name
			serverConfig.MicroServerConfig.TransactionResponder = transactionResponder
			serverConfig.MicroServerConfig.Viper = v
			serverConfig.MicroServerConfig.ListenAddress = listenAddress
			serverConfig.Service = newService

			newServer, err = server.New(serverConfig)
			if err != nil {
				panic(err)
			}
		}

		return newServer
	}

	// Create a new microkit command which manages our custom microservice.
	var newCommand command.Command
	{
		commandConfig := command.DefaultConfig()

		commandConfig.Logger = newLogger
		commandConfig.ServerFactory = newServerFactory

		commandConfig.Description = description
		commandConfig.GitCommit = gitCommit
		commandConfig.Name = name
		commandConfig.Source = source

		newCommand, err = command.New(commandConfig)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	newCommand.CobraCommand().Execute()

	return nil
}

func waitForFlannelFile(newLogger micrologger.Logger) error {
	var flannelFile string = os.Getenv("NETWORK_ENV_FILE_PATH")
	// wait for file creation
	for count := 0; ; count++ {
		// don't wait forever, if file is not created within retry limit, exit with failure
		if count > MaxRetry {
			newLogger.Log(fmt.Sprint("After 100sec flannel file is not created. Exiting"))
			return microerror.New("Failed to read flannel file.")
		}
		// check if file exists
		if _, err := os.Stat(flannelFile); !os.IsNotExist(err) {
			break
		}
		newLogger.Log("debug", fmt.Sprintf("Waiting for file '%s' to be created.", flannelFile))
		time.Sleep(1 * time.Second)
	}
	// all good
	return nil
}

func buildListenAddress() (string, error) {
	vniStr := strings.Split(os.Getenv("NETWORK_FLANNEL_DEVICE"), ".")[1]
	vni, err := strconv.Atoi(vniStr)
	if err != nil {
		return "", microerror.Maskf(err, "Failed to parse VNI from flannel device %s", vniStr)
	}
	listenPort := ListenPortBase + vni
	listenAddress := "http://127.0.0.1:" + strconv.Itoa(listenPort)

	return listenAddress, nil
}
