package commands

import (
	cc "../../common"
	"fmt"
	"github.com/doblenet/go-doblenet/tracer"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"os"
)



// $ consul-sc --consul=consul.service.consul:8500 ls -name=glob | -tags tagset
var cmdLsLL bool = false

var cmdLs = &cobra.Command{
	Use:   "ls",
	Aliases: []string{"list"},
	Short: "List available services (optionally matching the provided spec)",
	Long: `Queries the target Consul Service Catalog and lists the available entries that satisfy the provided predicates, all if no filtering predicates given.
	
Valid predspecs are:
  -name <glob>
  -port NUM
  -node <name>
  -tags <tagset>`,
	ValidArgs: []string{"-name", "-port", "-node", "-tags"},
	Run:       lsRun,
}

func init() {
	cf := cmdLs.Flags()
	setupCommonFlags(cf)
	cf.BoolVarP(&cmdLsLL, "long", "l", false, "Request long listing (with tags)")
}

func lsRun(cmd *cobra.Command, args []string) {

	client, err := cc.ConsulClient(consulConf)
	if nil != err {
		tracer.FatalErr(err)
	}

	catalog := client.Catalog()
	// 	if nil != err {
	// 		tracer.FatalErr(err)
	// 	}

	result, _, err := catalog.Services(&consulapi.QueryOptions{
		Datacenter:        consulConf.Datacenter,
		AllowStale:        true,
		RequireConsistent: false,
	})
	if nil != err {
		tracer.FatalErr(err)
	}

	for k,v := range result {
		
		fmt.Println(lsFormat(k, v, cmdLsLL))
	}
	
	sp := cc.ServicePredicate{}
	


	os.Exit(0)
}

func lsFormat(name string, tags []string, ll bool) string {
	
	if ll {
		return fmt.Sprintf("%s\t%s", name, tags)
	} else {
		return name
	}
}

func lsFilterEntry(filters []cc.ServicePredicate, name string, tags []string) bool {

	for i,p := range filters {
	
// 		switch p.PredCat {
// 			
// 		}
		
	}
	return false
}
