package app

import (
	"github.com/spf13/cobra"
	"pingmesh-agent/cmd/pingmesh-agent/app/options"
)

func NewPingmeshAgentCommand(stopCh <-chan struct{}) *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Short: "Launch pingmesh-agent",
		Long:  "Launch pingmesh-agent",
		RunE: func(c *cobra.Command, args []string) error {
			if err := runCommand(opts, stopCh); err != nil {
				return err
			}
			return nil
		},
	}
	//opts.Flags(cmd)
	return cmd
}

func runCommand(o *options.Options, ch <-chan struct{}) error {
	//if o.ShowVersion {
	//	fmt.Println(version.VersionInfo())
	//	os.Exit(0)
	//}
	config, err := o.PingmeshAgentConfig()
	if err != nil {
		return err
	}

	pm, err := config.Complete()
	if err != nil {
		return err
	}

	//err = ms.AddHealthChecks(healthz.NamedCheck("healthz", ms.CheckHealth))
	//if err != nil {
	//	return err
	//}
	return pm.RunUntil(ch)
}

