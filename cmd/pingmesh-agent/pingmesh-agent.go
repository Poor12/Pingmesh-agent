package main

import (
	"k8s.io/apimachinery/pkg/util/wait"
	"os"
	"pingmesh-agent/cmd/pingmesh-agent/app"
	"runtime"
	"k8s.io/component-base/logs"
)

func main(){
	logs.InitLogs()
	defer logs.FlushLogs()

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	cmd := app.NewPingmeshAgentCommand(wait.NeverStop)
	//cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}

}