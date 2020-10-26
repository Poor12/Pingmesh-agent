package pingmesh_agent

import (
	"time"
)

type Config struct {
	//Rest             *rest.Config
	Resolution time.Duration
}

func (c Config) Complete() (*PingmeshAgent, error) {
	store := NewStorage()
	pt := NewPingtool(store)
	return &PingmeshAgent{
		pingTool:   pt,
		storage:	store,
		resolution: c.Resolution,
	}, nil
}

