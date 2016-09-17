package commands

import (
	cc "../../common"
	"../../misc"
	"fmt"
	"github.com/doblenet/go-doblenet/tracer"
 	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"os"
)

var (
	getJSON	bool
	kvFlags	int64
)

var cmdGet = &cobra.Command{
	Use:	"get <key> [<key2> ...]",
	Short:	"Get value(s) for the specified key (or keys)",
	Long:	`Queries the target Consul KV store and retrieves the value at that key`,
	Run:	doGet,
}

var cmdSet = &cobra.Command{
	Use:	"set <key> <value>",
	Short:	"Set <key> to the specified value",
	Long:	`Sets a key at the target Consul KV store to the provided value.`,
	Run:	doSet,
}

var cmdSetMulti = &cobra.Command{
	Use:	"multiset [{<key> <value>} ...]",
	Short:	"Set <key> to the specified <value>, for every pair given",
	Long:	`Sets keys at the target Consul KV store to the provided values.`,
	Aliases: []string{"set-multi"},
	Run:	doSetMulti,
}


func init() {
	cf1 := cmdGet.Flags()
	setupCommonFlags(cf1)
	cf1.BoolVarP(&getJSON, "pretty-json", "h", false, "Format result as JSON (vs 'plain object')")
	
	cf2 := cmdSet.Flags()
	setupCommonFlags(cf2)
	cf2.Int64VarP(&kvFlags,"flags","f",0,"Flags for the value to be set")
	
	cf3 := cmdSetMulti.Flags()
	setupCommonFlags(cf3)
	cf3.Int64VarP(&kvFlags,"flags","f",0,"Flags for the value to be set")
}

func doGet(cmd *cobra.Command, args []string) {
	
	if len(args)<1 {
		tracer.Error("Required argument <key> not provided")
		cmd.HelpFunc()(cmd, []string{})
		os.Exit(1)
	}
	
	consul, err := cc.ConsulClient(consulConf)
	if nil != err {
		tracer.FatalErr(err)
	}
	
	kv := consul.KV()
	
	for _,x := range args {
		
		w := keyName(x)
		res,_,err := kv.Get(w, &consulapi.QueryOptions{
 			Datacenter:        consulConf.Datacenter,
			AllowStale:        true,
			RequireConsistent: false,
		})
		cc.CheckServerError(err)
		if nil!=err {
			tracer.FatalErr(err)
		}
		
		if nil==res {
			tracer.Warn("No result returned")
			continue
		}
		
		if !getJSON {
			fmt.Println(string(res.Value))
		} else {
			misc.PrettyJSON(os.Stdout,res.Value)
		}
	}
	
	os.Exit(0)
}

func doSet(cmd *cobra.Command, args []string) {
	
	if len(args)<2 {
		tracer.Error("Required arguments <key> <value> not provided")
		cmd.HelpFunc()(cmd, []string{})
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
	
	keyname := keyName(args[0])
	
	// XXX: TODO: add more functionality
	var p = consulapi.KVPair{Key: keyname, Value: []byte(args[1])}
	if 0 != kvFlags {
		p.Flags = uint64(kvFlags)
	}

	_, e := kv.Put(&p,&wo)
	cc.CheckServerError(e)
	if nil!= e {
		tracer.FatalErr(e)
	}
	
	os.Exit(0)
}

func doSetMulti(cmd *cobra.Command, args []string) {

	if 0 != len(args)%2 {
		tracer.Fatal("Provided arguments need to be even (key-value pairs)")
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

	var v consulapi.KVPair
	for i:=0 ; i <= len(args)/2; i+=2 {
	
		// reset the KVPair
		v = consulapi.KVPair{}
		// ..next iteration data (including flags)
		v.Key = keyName(args[i])
		v.Value = []byte(args[i+1])
		if 0 != kvFlags {
			v.Flags = uint64(kvFlags)
		}
		
		_,e := kv.Put(&v,&wo)
		cc.CheckServerError(e)
		if nil!= e {
			tracer.FatalErr(e)
		}
		
		tracer.TraceV(2,"Putting KV:", cc.KVPair2String(&v))
	}

	os.Exit(0)
}
