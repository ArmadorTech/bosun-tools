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
	lsLL, lsLX bool
)

var cmdList = &cobra.Command{
	Use:	"ls [<prefix>] [<pattern>]",
	Short:	"List existing keys under the provided prefix (optionally matching the provided pattern)",
	Long:	`Queries the target Consul KV store and lists the available keys -- optionally matching the provided glob.`,
	Aliases: []string{"list"},
	Run:	doLs,
}

var cmdTree = &cobra.Command{
	Use:	"tree [<prefix>]",
	Short:	"Recursively list keys under the provided prefix",
	Long:	`Queries the target Consul KV store and recursively  the available keys`,
	Run:	doTree,
}


func init() {
	cf1 := cmdList.Flags()
	setupCommonFlags(cf1)
 	cf1.BoolVarP(&lsLL, "long", "l", false, "Request long listing format")
// 	cf1.BoolVarP(&lsLX, "extended", "x", false, "Request extended listing format")

	cf2 := cmdTree.Flags()
	setupCommonFlags(cf2)
	cf2.BoolVarP(&lsLL, "long", "l", false, "Request long listing format")
	cf2.BoolVarP(&lsLX, "extended", "x", false, "Request extended listing format")
	
}

func doLs(cmd *cobra.Command, args []string) {
	
	var prefix []string;
	if len(args) < 1 {
		prefix = []string{"/"}
	} else {
		prefix = args
	}
	
	consul, err := cc.ConsulClient(consulConf)
	if nil != err {
		tracer.FatalErr(err)
	}
	
	kv := consul.KV()
	
	for _,p := range prefix {
		
		w := keyName(p)
		if lsLL {
			fmt.Printf("%s:\n", w)
		}
		result,_,err := kv.Keys(w, kdelim, &consulapi.QueryOptions{
			Datacenter:        consulConf.Datacenter,
			AllowStale:        true,
			RequireConsistent: false,
		})	
		cc.CheckServerError(err)
		if nil!=err {
			tracer.FatalErr(err)
		}

 		var i int
 		var k string
		for i, k = range result {
			if lsLL { 
				fmt.Println(k)
			} else {
				fmt.Printf("%s ",k)
			}
		}
		if lsLL {
			fmt.Print("Total ",i+1)
		}
	}
	fmt.Println()
	os.Exit(0)
}

func doTree(cmd *cobra.Command, args []string) {
	
	if len(args)<1 {
		tracer.Error("Required argument <prefix> not provided")
		cmd.HelpFunc()(cmd, []string{})
		os.Exit(1)
	}
	
	consul, err := cc.ConsulClient(consulConf)
	if nil != err {
		tracer.FatalErr(err)
	}

	kv := consul.KV()
	w := keyName(args[0])
	result,_,err := kv.List(w, &consulapi.QueryOptions{
		Datacenter:        consulConf.Datacenter,
		AllowStale:        true,
		RequireConsistent: false,
	})
	cc.CheckServerError(err)
	if nil!=err {
		tracer.FatalErr(err)
	}
	
	for _,p := range result {
		
		if !lsLL && !lsLX {
			fmt.Println(p.Key)
		} else {
			fmt.Println(cc.KVPairtoString(p, lsLL,lsLX))
		}
	}
	
	os.Exit(0)
}
