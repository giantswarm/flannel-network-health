package service

import (
	"github.com/giantswarm/microerror"
)

var invalidConfigError = microerror.New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidFlannelFileError = microerror.New("invalid flannel file")

func IsInvalidFlannelFile(err error) bool {
	return microerror.Cause(err) == invalidFlannelFileError
}

var invalidFlannelConfigurationError = microerror.New("Unable to find FLANNEL_SUBNET in flannel file")

func IsInvalidFlannelConfiguration(err error) bool {
	return microerror.Cause(err) == invalidFlannelConfigurationError
}

var parsingFlannelSubnetError = microerror.New("parsing flannel file")

func IsParsingFlannelSubnet(err error) bool {
	return microerror.Cause(err) == parsingFlannelSubnetError
}
