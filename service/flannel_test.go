package service

import (
	"testing"

	"github.com/giantswarm/flannel-network-health/flag/service/network"
	"github.com/giantswarm/microerror"
)

func Test_Flannel_ParseIP(t *testing.T) {
	tests := []struct {
		config             func(flannelFile []byte) (network.Network,error)
		flannelFileContent []byte
		expectedConfig     network.Network
		expectedErr        error
	}{
		// test 1
		{
			config: func(flannelFile []byte) (network.Network,error) {
				conf := DefaultConfig()
				conf.parseIPs(flannelFile)
				return conf.Flag.Service.NetworkConfig, nil
			},
			expectedConfig: network.Network{FlannelIP: "172.23.3.65", BridgeIP: "172.23.3.66"},
			flannelFileContent: []byte(`FLANNEL_NETWORK=172.23.3.0/24
FLANNEL_SUBNET=172.23.3.65/30
FLANNEL_MTU=1450
FLANNEL_IPMASQ=false`),
expectedErr:nil,
		},
		// test 2
		{
			config: func(flannelFile []byte) (network.Network,error) {
				conf := DefaultConfig()
				conf.parseIPs(flannelFile)
				return conf.Flag.Service.NetworkConfig, nil
			},
			expectedConfig: network.Network{FlannelIP: "198.168.0.0", BridgeIP: "192.168.0.1"},
			flannelFileContent: []byte(`FLANNEL_NETWORK=192.168.0.0/24
FLANNEL_SUBNET=198.168.0.0/30
FLANNEL_MTU=1450
FLANNEL_IPMASQ=false`),
			expectedErr:nil,
		},
		// test 3 - missing FLANNEL_SUBNET
		{
			config: func(flannelFile []byte) (network.Network,error) {
				conf := DefaultConfig()
				conf.parseIPs(flannelFile)
				return conf.Flag.Service.NetworkConfig, nil
			},
			expectedConfig: network.Network{FlannelIP: "198.168.0.0", BridgeIP: "192.168.0.1"},
			flannelFileContent: []byte(`FLANNEL_NETWORK=192.168.0.0/24
FLANNEL_MTU=1450
FLANNEL_IPMASQ=false`),
			expectedErr:invalidFlannelConfiguration,
		},
		// test 4 - invalid subnet in flannel file
		{
			config: func(flannelFile []byte) (network.Network,error) {
				conf := DefaultConfig()
				conf.parseIPs(flannelFile)
				return conf.Flag.Service.NetworkConfig, nil
			},
			expectedConfig: network.Network{FlannelIP: "198.168.0.0", BridgeIP: "192.168.0.1"},
			flannelFileContent: []byte(`FLANNEL_NETWORK=192.168.0.0/24
FLANNEL_SUBNET=x.68.c.0/30
FLANNEL_MTU=1450
FLANNEL_IPMASQ=false`),
			expectedErr:errorParsingFLannelSubnet,
		},
	}

	for index, test := range tests {
		networkConfig, err := test.config(test.flannelFileContent)

		if microerror.Cause(err) != microerror.Cause(test.expectedErr) {
			t.Fatalf("%v: unexcepted error, expected %v but got %v", index, test.expectedErr, err)
		}
		if networkConfig.FlannelIP != test.expectedConfig.FlannelIP {
			t.Fatalf("%v: Incorrent ip, expected %v but got %v.", index, test.expectedConfig.FlannelIP, networkConfig.FlannelIP)
		}
		if networkConfig.BridgeIP != test.expectedConfig.BridgeIP {
			t.Fatalf("%v: Incorrent ip, expected %v but got %v.", index, test.expectedConfig.BridgeIP, networkConfig.BridgeIP)
		}
	}
}
