package commands

import (
	_ "errors"
 	"fmt"
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

	
	os.Exit(0)
}
