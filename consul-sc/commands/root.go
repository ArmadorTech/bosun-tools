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
	Use:   "consul-sc [global_opts] command [cmd_opts]",
	Short: "Operate on Consul's Service Catalog",
	Long: `Consul-SC is a tool designed to ease operating on a consul[by HashiCorp]-operated Service Catalog.
The tool is intended to replace any and all (raw)HTTP-based interactions, such as cURL-based scripts, with a modern and user-friendly CLI tool`,
	Run: func(cmd *cobra.Command, args []string) {

		// Do Stuff Here
		fmt.Println("Running at ", misc.LocalHostname(), " -- ", misc.LocalIPs())

		fmt.Println("verbosity->", verbosity)

		fmt.Println("direct->", consulConf.AgentAPI)
		fmt.Println("consul->", consulConf.URL)
		fmt.Println("datacenter->", consulConf.Datacenter)

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

	pf := RootCmd.PersistentFlags()
	pf.CountVarP(&verbosity, "verbose", "v", "Enable (and/or increase) verbosity level")
	pf.BoolVarP(&consulConf.AgentAPI, "passthrough", "z", false, "Enable direct operation against a Consul server, bypassing the (normally local) agent")

	pf.StringVar(&consulConf.URL, "consul", "", "Consul HTTP API endpoint to use (default: $CONSUL_HOST)")
	pf.StringVar(&consulConf.Datacenter, "dc", "", "Consul datacenter identifier to use")

	var url, dc string
	if url = os.Getenv(k_ENV_CONSUL); "" == url {
		url = k_CONSUL_URL
	}
	if dc = os.Getenv(k_ENV_DC); "" == dc {
		dc = k_CONSUL_DC
	}
	consulConf.URL = url
	consulConf.Datacenter = dc

	RootCmd.AddCommand(cmdRegister)

	RootCmd.AddCommand(cmdLs)

	RootCmd.AddCommand(cmdInspect)

	RootCmd.AddCommand(cmdQuery)

	RootCmd.AddCommand(cmdMaintenance)

	RootCmd.AddCommand(cmdDatacenter)
}

func setupCommonFlags(ff *pflag.FlagSet) {
	ff.StringVar(&consulConf.Token, "token", "", "Provide Consul authorization token")
}
