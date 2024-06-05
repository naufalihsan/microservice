package inmem

import (
	"context"
)

type Registry struct{}

func NewRegistry() (*Registry, error) {
	return &Registry{}, nil
}

func (r *Registry) Register(ctx context.Context, instanceId, serviceName, hostPort string) error {
	return nil
}

func (r *Registry) Deregister(ctx context.Context, instanceId, serviceName string) error {
	return nil
}

func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	return []string{}, nil
}

func (r *Registry) HealthCheck(instanceId, serviceName string) error {
	return nil
}
