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
	forceRm	bool = false
)

// $ consul-node --consul=consul.service.consul:8500 rm|remove <node_name>

var cmdRemove = &cobra.Command{
	Use:	"rm <nodename>",
	Short:	"Remove the named node from the Catalog",
	Long:	`Removes the named node from the target Consul Catalog.`,
	Aliases: []string{"remove"},
	Run:	doRemove,
}

func init() {
	cf := cmdRemove.Flags()
	setupCommonFlags(cf)
	cf.BoolVarP(&forceRm, "force", "f", false, "Force the removal")

}

func doRemove(cmd *cobra.Command, args []string) {

	fmt.Println("force=",forceRm)

	if len(args)<1 {
		fmt.Printf("FATAL: No arguments provided")
		os.Exit(1)
	}
	
	consul, err := cc.ConsulClient(consulConf)
	if nil != err {
		tracer.FatalErr(err)
	}
	
	var wo = consulapi.WriteOptions{Datacenter: consulConf.Datacenter}
	if ""!=consulConf.Token {
		wo.Token = consulConf.Token
	}
	
	catalog := consul.Catalog()
	_, err = catalog.Deregister(&consulapi.CatalogDeregistration{
		Node: args[0],
		Datacenter: consulConf.Datacenter,
	}, &wo)
	cc.CheckServerError(err)
	if nil!=err {
		tracer.FatalErr(err)
	}
	
	os.Exit(0)
}
