package options

import (
	pingmesh_agent "pingmesh-agent/pkg/pingmesh-agent"
	"time"
)

type Options struct {
	resolution time.Duration
}

func NewOptions() *Options {
	o := &Options{
		resolution: 60 * time.Second,
	}
	return o
}

func (o Options) PingmeshAgentConfig() (*pingmesh_agent.Config, error) {
	return &pingmesh_agent.Config{
		Resolution: o.resolution,
	}, nil
}
