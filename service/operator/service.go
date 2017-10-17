package operator

import (
"sync"
"time"

"github.com/giantswarm/microerror"
"github.com/giantswarm/micrologger"
)

// Config represents the configuration used to create a new service.
type Config struct {
	// Dependencies.
	Logger            micrologger.Logger
}

// DefaultConfig provides a default configuration to create a new service by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:            nil,
	}
}

// Service implements the operator service.
type Service struct {
	// Dependencies.
	logger            micrologger.Logger

	// Internals.
	bootOnce       sync.Once
	mutex          sync.Mutex
}

// New creates a new configured service.
func New(config Config) (*Service, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	newService := &Service{
		// Dependencies.
		logger:            config.Logger,

		// Internals
		bootOnce:       sync.Once{},
		mutex:          sync.Mutex{},
	}

	return newService, nil
}

// Boot starts the service and implements the watch for the cluster TPR.
func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		// dummy  operator wait
		s.logger.Log("Dummy operator boot.")
		for ;; {
			time.Sleep(1*time.Hour)
		}

	})
}
