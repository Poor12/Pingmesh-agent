package pingmesh_agent

import (
	"context"
	"k8s.io/klog/v2"
	"sync"
	"time"
)

type PingmeshAgent struct {
	resolution    time.Duration
	pingTool	  *PingTool
	storage       *Storage
	healthMu      sync.RWMutex
	lastTickStart time.Time
	lastOk        bool
}

func (pm *PingmeshAgent) RunUntil(stopCh <-chan struct{}) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return pm.runPing(ctx)
}

func (pm *PingmeshAgent) runPing(ctx context.Context) error {
	ticker := time.NewTicker(pm.resolution)
	defer ticker.Stop()
	pm.ping(ctx, time.Now())

	for {
		select {
		case startTime := <-ticker.C:
			pm.ping(ctx, startTime)
		case <-ctx.Done():
			return nil
		}
	}
}

func (pm *PingmeshAgent) ping(ctx context.Context, startTime time.Time) {
	pm.healthMu.Lock()
	pm.lastTickStart = startTime
	pm.healthMu.Unlock()

	healthyTick := true

	ctx, cancelTimeout := context.WithTimeout(ctx, pm.resolution)
	defer cancelTimeout()

	klog.V(6).Infof("Beginning cycle, starting to ping nodes...")
	err := pm.pingTool.Ping(ctx)
	if err != nil{
		healthyTick = false
		klog.Errorf("Ping failed....")
	}

	pm.healthMu.Lock()
	pm.lastOk = healthyTick
	pm.healthMu.Unlock()
}
