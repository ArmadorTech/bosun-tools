package commands

import (
	cc "../../common"
	"../../misc"
	"fmt"
	"github.com/doblenet/go-doblenet/tracer"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"os"
	"encoding/json"
)

// $ consul-sc --consul=consul.service.consul:8500 inspect <ID>

var cmdInspectJSON bool
var cmdInspect = &cobra.Command{
	Use:   "inspect <ID>",
	Short: "Return detailed information on a service, given its ID or name",
	Long:  `Queries the target Consul Service Catalog and lists the available entries`,
	Run:   inspectRun,
}

func init() {
	cf := cmdInspect.Flags()
	setupCommonFlags(cf)
	cf.BoolVarP(&cmdInspectJSON, "pretty-json", "x", false, "Format result as JSON (vs 'plain object')")
}

func inspectRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("FATAL: required argument <ID> not provided")
		os.Exit(3)
	}

	serv_name := args[0]
	
	client, err := cc.ConsulClient(consulConf)
	if nil != err {
		tracer.FatalErr(err)
	}

	// XXX: TODO: implement "-z"
	catalog := client.Catalog()

	result, _, err := catalog.Service(serv_name, "", &consulapi.QueryOptions{
		Datacenter:        consulConf.Datacenter,
		AllowStale:        true,
		RequireConsistent: false,
	})
	if nil != err {
		tracer.FatalErr(err)
	}

	if 0 == len(result) {
		fmt.Println("No service found matching the provided name:", serv_name)
		os.Exit(0)
	}
	
	for _,v := range result {
		
		o, _ := json.Marshal(*v)
		if cmdInspectJSON {
			misc.PrettyJSON(os.Stdout, o)
		} else {
			misc.OutputJSON(os.Stdout, o)
		}
	}
	
	os.Exit(0)
}
