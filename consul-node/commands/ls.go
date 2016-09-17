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

// $ consul-node --consul=consul.service.consul:8500 ls [-l][-x]

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
	
	var e error
	if lsLX {
		e = nodeXList(consul, lsLL)
	} else {
		e = nodeList(consul, lsLL)
	}
	cc.CheckServerError(e)
	if nil != err {
		tracer.FatalErr(e)
	}
	
	os.Exit(0)
}

func nodeList(client *consulapi.Client, ll bool) error {
	
	catalog := client.Catalog()
	result, _, err := catalog.Nodes(&consulapi.QueryOptions{
		Datacenter:        consulConf.Datacenter,
		AllowStale:        true,
		RequireConsistent: false,
	})
	if nil!=err {
		return err
	}
	
	for _,x := range result {
	
		if !ll {
			fmt.Println(x.Node)
		} else {
			// XXX: API v0.7
			fmt.Printf("%s\t%s %s\n", x.Node, x.Address,
							cc.Map2String(x.TaggedAddresses))
		}
	}
	return nil
}

func nodeXList(client *consulapi.Client, ll bool) error {
	
	coord := client.Coordinate()	
	result, _, err := coord.Nodes(&consulapi.QueryOptions{
		Datacenter:        consulConf.Datacenter,
		AllowStale:        true,
		RequireConsistent: false,
	})
	if nil != err {
		return err
	}
	
	for _, x := range result {
		
		fmt.Printf("%s\t", x.Node)
		if ll {	// doesn't make much sense without ...
			fmt.Println(cc.CEtoString(*x))
		}
	}
	return nil
}
