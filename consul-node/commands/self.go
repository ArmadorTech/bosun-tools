package commands

import (
	cc "../../common"
	"fmt"
	"github.com/doblenet/go-doblenet/tracer"
// 	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"os"
)

var (
 	infoLX	bool
)
// $ consul-node --consul=consul.service.consul:8500 self
// $ consul-node --consul=consul.service.consul:8500 info [-x]

var cmdSelf = &cobra.Command{
	Use:	"self",
	Short:	"Returns the target's nodename",
	Long:	`Queries the target Consul agent and returns its hostname`,
	Aliases: []string{"who"},
	Run:	doSelf,
}
var cmdInfo = &cobra.Command{
	Use:	"info",
	Short:	"Returns info about the target agent",
	Long:	`Queries the target Consul agent and returns all available information`,
	Run:	doSelf,
}

func init() {
	cf := cmdSelf.Flags()
	setupCommonFlags(cf)
	cf2 := cmdInfo.Flags()
	setupCommonFlags(cf2)
 	cf.BoolVarP(&infoLX, "extended", "x", false, "Request extended listing")
}

func doSelf(cmd *cobra.Command, args []string) {

	consul, err := cc.ConsulClient(consulConf)
	if nil != err {
		tracer.FatalErr(err)
	}

	agent := consul.Agent()
	
	switch cmd.CalledName() {
		case "who": fallthrough
		case "self":
			name, e := agent.NodeName()
			cc.CheckServerError(e)
			if nil!=e {
				tracer.FatalErr(e)
			}
	
			fmt.Println(name)
		
		case "info":
	
			result, e := agent.Self()
			cc.CheckServerError(e)
			if nil!=e {
				tracer.FatalErr(e)
			}
			cc.PrintPropMap(os.Stdin, result, false)
	}
	
	os.Exit(0)
}
