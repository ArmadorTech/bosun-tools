package commands

import (
	cc "../../common"
	misc "../../misc"
	"fmt"
	"github.com/doblenet/go-doblenet/tracer"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
)

const (
	k_ENV_CONSUL = "CONSUL_HOST"
	k_ENV_DC     = "CONSUL_DC"

	k_CONSUL_URL = "localhost:8500"
	k_CONSUL_DC  = "dc1"
)

var (
	consulConf cc.ClientConfig
	// global vars for commands [from flags]
	verbosity int = 0
)

var RootCmd = &cobra.Command{
	Use:   "consul-node [global_opts] command [cmd_opts]",
	Short: "Operate on Consul's Catalog -- Node section",
	Long: `Consul-Node is a tool designed to ease operating on a consul[by HashiCorp]-operated Catalog.
The tool is intended to replace any and all (raw)HTTP-based interactions, such as cURL-based scripts, with a modern and user-friendly CLI tool`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Running at ", misc.LocalHostname(), "--", misc.LocalIPs())

		fmt.Println("verbosity->", verbosity)
		cc.PrintConfig(&consulConf)
		
		// invoke help...
		cmd.HelpFunc()(cmd, []string{})
		
		os.Exit(0)
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		tracer.SetLevel(verbosity)
		return nil
	},
}

func init() {

	cc.SetupConsulFlags(RootCmd.PersistentFlags(),
		&consulConf, 
		&verbosity,
	)

	RootCmd.AddCommand(cmdCreate)
	RootCmd.AddCommand(cmdRemove)

	RootCmd.AddCommand(cmdLs)

 	RootCmd.AddCommand(cmdInspect)

}

func setupCommonFlags(ff *pflag.FlagSet) {
	ff.StringVar(&consulConf.Token, "token", "", "Provide Consul authorization token")
}
