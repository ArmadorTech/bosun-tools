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
	lsLL,lsLX	bool
)

// $ consul-node --consul=consul.service.consul:8500 rm|remove <node_name>

var cmdLs = &cobra.Command{
	Use:	"ls [<nodename glob>]",
	Short:	"List available nodes (optionally matching the provided pattern)",
	Long:	`Queries the target Consul Catalog and lists the available nodes optionally matching the provided glob.`,
	Aliases: []string{"list"},
	Run:	doLS,
}

func init() {
	cf := cmdLs.Flags()
	setupCommonFlags(cf)
	cf.BoolVarP(&lsLL, "long", "l", false, "Request long listing format")
	cf.BoolVarP(&lsLX, "extended", "x", false, "Request extended listing format")
}

func doLS(cmd *cobra.Command, args []string) {

	consul, err := cc.ConsulClient(consulConf)
	if nil != err {
		tracer.FatalErr(err)
	}
	
	catalog := consul.Catalog()
	result, _, err := catalog.Nodes(&consulapi.QueryOptions{
		Datacenter:        consulConf.Datacenter,
		AllowStale:        true,
		RequireConsistent: false,
	})
	if nil!=err {
		tracer.FatalErr(err)
	}
	
	for _,x := range result {
	
		fmt.Println(lsFormat(x,lsLL,lsLX))
		
	}
	os.Exit(0)
}

func lsFormat(n *consulapi.Node, ll, lx bool) string {
	switch {
		case true==ll:
			return fmt.Sprintf("%s\t%s", n.Node, n.Address)
		case true==lx:
// 			// wait till API v0.7
// 			return fmt.Sprintf("%s\t%s %v", n.Node, n.Address, n.TaggedAddresses)
			return fmt.Sprintf("%s\t%s", n.Node, n.Address)
		default:
			return n.Node
	}
}
