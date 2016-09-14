package commands

import (
	cc "../../common"
	"github.com/doblenet/go-doblenet/tracer"
 	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"os"
)


var cmdRm = &cobra.Command{
	Use:	"rm <key> [<key>...]",
	Short:	"Remove the specified keys(s)",
	Long:	`Removes the specified key(s) from the target Consul KV store.`,
	Aliases: []string{"delete"},
	Run:	doRm,
}

var cmdRmTree = &cobra.Command{
	Use:	"rmtree <key_prefix>]",
	Short:	"Remove the specified subtree",
	Long:	`Removes the specified subtree from the target Consul KV store.`,
	Aliases: []string{"delete-tree", "deltree"},
	Run:	doRmTree,
}

func init() {
	cf1 := cmdRm.Flags()
	setupCommonFlags(cf1)
	
	cf2 := cmdRmTree.Flags()
	setupCommonFlags(cf2)
}

func doRm(cmd *cobra.Command, args []string) {
	
	if len(args) < 1 {
		tracer.Fatal("Required argument <key> not provided")
		os.Exit(1)
	}
	
	consul, err := cc.ConsulClient(consulConf)
	if nil != err {
		tracer.FatalErr(err)
	}
	
	kv := consul.KV()
	var wo consulapi.WriteOptions
	if ""!=consulConf.Token {
		wo.Token = consulConf.Token
	}
	
	for _,k := range args {
		
		_,e := kv.Delete(k, &wo)
		cc.CheckServerError(e)
		if nil!=e {
			tracer.FatalErr(e)
		}
	
		tracer.TraceV(2,"Removed key",k)
	}
	os.Exit(0)
}

func doRmTree(cmd *cobra.Command, args []string) {
	
	if len(args) != 1 {
		tracer.Fatal("Required argument <prefix> not provided")
		os.Exit(1)
	}
	
	consul, err := cc.ConsulClient(consulConf)
	if nil != err {
		tracer.FatalErr(err)
	}
	
	kv := consul.KV()
	var wo consulapi.WriteOptions
	if ""!=consulConf.Token {
		wo.Token = consulConf.Token
	}
	
	if ""!=kvPrefix {
		tracer.Warn("prefix set: 'rmtree' ignores prefixes for safety")
	}
	
	_,e := kv.DeleteTree(args[0], &wo)
	cc.CheckServerError(e)
	if nil!=e {
		tracer.FatalErr(e)
	}
	
	tracer.TraceV(2,"Removed (sub-)tree", args[0])
	
	os.Exit(0)
}
