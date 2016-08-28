package commands

import (
	cc "../../common"
	"fmt"
	"github.com/doblenet/go-doblenet/tracer"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"os"
)


var (
	forceAdd	bool
	external	bool
)

// $ consul-node --consul=consul.service.consul:8500 add|create <node_name>

var cmdCreate = &cobra.Command{
	Use:	"register <nodename> [<ip>]",
	Short:	"Register the named node in the Catalog",
	Long:	`Creates a new node entry in the target Consul Catalog.`,
	Aliases: []string{"create", "add"},
	Run:	doRegister,
}

func init() {
	cf := cmdCreate.Flags()
	setupCommonFlags(cf)
	cf.BoolVarP(&external, "external", "", false, "Register an 'external' node")
	cf.BoolVarP(&forceAdd, "force", "f", false, "Force register, overwriting if needed")

}

func doRegister(cmd *cobra.Command, args []string) {

	if len(args) < 2 {
		fmt.Println("FATAL: required arguments missing: <nodename> <address>")
		os.Exit(3)
	}
	
	consul, err := cc.ConsulClient(consulConf)
	if nil != err {
		tracer.FatalErr(err)
	}
	
	catalog := consul.Catalog()
	
	var wo consulapi.WriteOptions
	if ""!=consulConf.Token {
		wo.Token = consulConf.Token
	}
	
	_, err = catalog.Register(&consulapi.CatalogRegistration{
		Node:		args[0],
		Address:	args[1],
//		TaggedAddresses map[string]string
		Datacenter:	consulConf.Datacenter,
//		Service         *AgentService
// 		Check           *AgentCheck
	}, &wo)
	
	if nil!=err {
		tracer.FatalErr(err)
	}
	
	os.Exit(0)
}
