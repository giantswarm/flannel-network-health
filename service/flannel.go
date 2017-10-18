package service

import (
	"fmt"
	"github.com/giantswarm/microerror"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
)

func (c *Config) LoadFlannelConfig() error {
	// fetch configuration from OS env
	confFile, err := c.fetchConfFromOS()
	if err != nil {
		return microerror.Mask(err)
	}

	// parse config and generate IP for interfaces
	err = c.parseIPs(confFile)
	if err != nil {
		return microerror.Mask(err)
	}
	// debug output
	c.Logger.Log("debug", fmt.Sprintf("Loaded Config: %+v", c.Flag.Service.NetworkConfig))
	return nil
}

// fetch ENVs values and read flannel file
func (c *Config) fetchConfFromOS() ([]byte, error) {
	// load NIC interfaces from ENV
	c.Flag.Service.NetworkConfig.BridgeInterface = os.Getenv("NETWORK_BRIDGE_NAME")
	c.Flag.Service.NetworkConfig.FlannelInterface = os.Getenv("NETWORK_FLANNEL_DEVICE")
	// read flannel file
	filename := os.Getenv("NETWORK_ENV_FILE_PATH")
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, microerror.Maskf(invalidFlannelFile, "%s", filename)
	}

	return fileContent, nil
}

// parse flannel configuration file and generate ips for interface
func (c *Config) parseIPs(confFile []byte) error {

	// get FLANNEL_SUBNET from flannel file via regexp
	r, _ := regexp.Compile("FLANNEL_SUBNET=[0-9]+.[0-9]+.[0-9]+.[0-9]+/[0-9]+")
	flannelLine := r.Find(confFile)
	// check if regexp returned non-empty line
	if len(flannelLine) < 5 {
		return microerror.Mask(invalidFlannelConfiguration)
	}

	// parse flannel subnet
	flannelSubnetStr := strings.Split(string(flannelLine), "=")[1]
	flannelIP, _, err := net.ParseCIDR(flannelSubnetStr)
	if err != nil {
		return microerror.Maskf(errorParsingFLannelSubnet, "%v", err)
	}
	// force ipv4 for later trick
	flannelIP = flannelIP.To4()

	// get bridge ip
	c.Flag.Service.NetworkConfig.BridgeIP = flannelIP.String()
	// get flannel ip,, which is just one number smaller than bridge hence the [3]++ trick
	flannelIP[3]--
	c.Flag.Service.NetworkConfig.FlannelIP = flannelIP.String()

	return nil
}
