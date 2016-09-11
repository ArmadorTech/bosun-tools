package commands

import (
	cc "../../common"
	"github.com/doblenet/go-doblenet/tracer"
// 	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"os"
)

var (
	maintenanceReason string
)

// Sets maintenance mode for the agent that we connect to
// $ consul-node --consul=consul.service.consul:8500 maintenance ['enable'|'disable'] [--reason=TEXT]
// $ consul-node --consul=consul.service.consul:8500 resume


var cmdMaintenance = &cobra.Command{
	Use:       "maintenance [enable|disable]",
	Short:     "Put the target node in maintenance mode or resume normal activity",
	Long:      `Allows placing the node into "maintenance mode" or "resuming" normal activity. During maintenance mode, all services registered at this node will be marked as unavailable and will not be present in DNS or API responses. Maintenance mode is persistent: survives an agent restart.`,
	Aliases:   []string{"resume"},
	ValidArgs: []string{"enable", "disable"},
	Run:       maintenanceRun,
}

func init() {
	cf := cmdMaintenance.Flags()
	setupCommonFlags(cf)
	cf.StringVar(&maintenanceReason, "reason", "", "Text explaining the reason for placing the node into maintenance mode")
}

func maintenanceRun(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		tracer.Fatal("required argument <ID> not provided")
		os.Exit(3)
	}

	consul,err := cc.ConsulClient(consulConf)
	cc.CheckServerError(err)
	if nil!=err {
		tracer.FatalErr(err)
	}
	
	agent := consul.Agent()
	
	switch cmd.CalledName() {
	case "maintenance":
		if len(args) < 2 {
			tracer.Warn("argument missing; Asuming 'maintenance enable'")
		} else {
			switch args[1] {
			case "enable":
				tracer.Trace(2,"ENABLE")
				if e := agent.EnableNodeMaintenance(maintenanceReason); nil!=e {
					tracer.FatalErr(e)
				}
				
			case "disable":
				tracer.Trace(2,"DISABLE")
				if e := agent.DisableNodeMaintenance(); nil!=e {
					tracer.FatalErr(e)
				}
			}
		}
		
	case "resume":
		tracer.Trace(2,"RESUME")
		if e := agent.DisableNodeMaintenance(); nil!=e {
			tracer.FatalErr(e)
		}
	}

	os.Exit(0)
}
