package interfaceHealthz


import (
	"context"
	"fmt"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/microendpoint/service/healthz"
	"github.com/giantswarm/flannel-network-health/service/healthz/interfaceHealthz/interface"
	"github.com/vishvananda/netlink"
	"github.com/giantswarm/flannel-network-health/service/healthz/interfaceHealthz/key"
)

const (
	// Description describes which functionality this health check implements.
	Description = "Ensure network interface is present and has proper network configuration."
	// Name is the identifier of the health check. This can be used for emitting
	// metrics.
	Name = "interfaceHealthz"
)

// Config represents the configuration used to create a healthz service.
type Config struct {
	// Dependencies.
	NetworkInterface _interface.NetworkInterface
	Logger    micrologger.Logger
}

// DefaultConfig provides a default configuration to create a new healthz service
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		NetworkInterface: _interface.NetworkInterface{},
		Logger:    nil,
	}
}

// Service implements the healthz service interface.
type Service struct {
	// Dependencies.
	networkInterface _interface.NetworkInterface
	logger    micrologger.Logger

	// Settings.
	timeout time.Duration
}

// New creates a new configured healthz service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.NetworkInterface.IP == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.NetworkInterface.IP must not be empty string")
	}
	if config.NetworkInterface.Name == "" {
		return nil, microerror.Maskf(invalidConfigError, "config.NetworkInterface.Name must not be empty string")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}


	newService := &Service{
		// Dependencies.
		networkInterface: config.NetworkInterface,
		logger:    config.Logger,
	}

	return newService, nil
}

// GetHealthz implements the health check for network interface.
func (s *Service) GetHealthz(ctx context.Context) (healthz.Response, error) {
	message := fmt.Sprintf("Healthcheck for interface %s has been successful. Interface is present and configured with ip %s.",s.networkInterface.Name,s.networkInterface.IP)

	failed, message := s.healthCheck(message)

	response := healthz.Response{
		Description: Description,
		Failed:      failed,
		Message:     message,
		Name:        Name,
	}

	return response, nil
}
// implementation fo the interface healthz logic
func (s *Service) healthCheck(message string) (bool, string){
	// load interface
	bridge, err := netlink.LinkByName(s.networkInterface.Name)
	if err != nil {
		message = fmt.Sprintf("Cant find interface %s. %s", s.networkInterface.Name, err)
		return true, message
	}
	// check ip on interface
	ipList, err := netlink.AddrList(bridge, netlink.FAMILY_V4)
	if err != nil || len(ipList) == 0 {
		message = fmt.Sprintf("Missing ip %s in the bridge configuration.", s.networkInterface.IP)
		return true, message
	}
	// compare ip on interface
	if len(ipList) > 0 &&  key.GetInterfaceIP(ipList) != s.networkInterface.IP {
		message = fmt.Sprintf("Wrong ip on interface %s. Expected %s, but found %s.", s.networkInterface.Name, s.networkInterface.IP, key.GetInterfaceIP(ipList))
		return true, message
	}
	// all good
	return false, message
}