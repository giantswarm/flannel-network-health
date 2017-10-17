package service

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidFlannelFile = microerror.New("Error reading flannel file")

func IsInvalidFlannelFile(err error) bool {
	return microerror.Cause(err) == invalidFlannelFile
}

var invalidFlannelConfiguration = microerror.New("Unable to find FLANNEL_SUBNET in flannel file")

func IsInvalidFlannelConfiguration(err error) bool {
	return microerror.Cause(err) == invalidFlannelConfiguration
}

var errorParsingFLannelSubnet = microerror.New("Unable to find FLANNEL_SUBNET in flannel file")

func IsErrorParsingFLannelSubnet(err error) bool {
	return microerror.Cause(err) == errorParsingFLannelSubnet
}
