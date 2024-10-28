package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Registery interface {
	// Register a service with the registry.
	Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error
	// Deregister a service with the registry.
	Deregister(ctx context.Context, instanceID string, serviceName string) error
	// Discover a service with the registry.
	Discover(ctx context.Context, serviceName string) ([]string, error)
	// HealthCheck a service with the registry.
	HealthCheck(instanceID string, serviceName string) error
	// GetIdleInterval for service healthcheck.
	GetIdleInterval() time.Duration
}

// ErrNotFound is returned when no service addresses are found.
var ErrNotFound = errors.New("no service addresses found")

// ErrNotFound is returned when no service instance are found.
var ErrNotFoundInstance = errors.New("no service instance found")

// ErrNotFound is returned when no service found.
var ErrNotFoundService = errors.New("no service found")

// GenerateInstanceID generates a psuedo-random service instance identifier, using a service name. Suffixed by dash and number
func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
