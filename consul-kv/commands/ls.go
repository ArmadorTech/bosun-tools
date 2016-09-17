package commands

import (
	cc "../../common"
	"fmt"
	"github.com/doblenet/go-doblenet/tracer"
 	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"io"
	"os"
	"text/tabwriter"
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
	
	var e error
	n := len(prefix)
	for _,p := range prefix {
		
		w := keyName(p)
		
		// Output header for every entry (multi-ls)
		if n>1 {
			fmt.Printf("%s:\n", w)
		}

		if !lsLL {
			e = doListS(os.Stdout, kv,w)
		} else {
			e = doListL(os.Stdout, kv,w)
		}
		cc.CheckServerError(e)
		if nil!=e {
			tracer.FatalErr(e)
		}
	}
	os.Exit(0)
}

func doListS(w io.Writer, kv *consulapi.KV, p string) error {

	// Simple key listing (just use "Keys")
	result,_,err := kv.Keys(p, kdelim, &consulapi.QueryOptions{
		Datacenter:        consulConf.Datacenter,
		AllowStale:        true,
		RequireConsistent: false,
	})
	if nil != err {
		return err
	}

	for _,k := range result {
		fmt.Fprintf(w, "%s ",k)
	}
	return nil
}

func doListL(w io.Writer, kv *consulapi.KV, p string) error {

	// Extended key listing: "List" needed
	result,_,err := kv.List(p, &consulapi.QueryOptions{
		Datacenter:        consulConf.Datacenter,
		AllowStale:        true,
		RequireConsistent: false,
	})
	if nil != err {
		return err
	}

	t := tabwriter.NewWriter(w, 3,4,1,' ',0)
	for _,v := range result {
		fmt.Fprintln(t, cc.KVPair2String(v))
	}
	t.Flush()
	return nil
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
	
	t := tabwriter.NewWriter(os.Stdout, 7,4,1,' ',0)
	for _,v := range result {
		
		if !lsLL && !lsLX {
			fmt.Println(v.Key)
		} else {
			fmt.Fprintln(t,cc.KVPair2String(v))
		}
	}
	t.Flush()
	os.Exit(0)
}
