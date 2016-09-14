package commands

import (
	cc "../../common"
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
	
	kvPrefix	string
	kdelim		string
)

var RootCmd = &cobra.Command{
	Use:   "consul-kv [global_opts] command [cmd_opts]",
	Short: "Operate on Consul's KV interface",
	Long: `Consul-kv is a tool designed to ease operating on a consul[by HashiCorp] Key-Value distributed store.
The tool is intended to replace any and all (raw)HTTP-based interactions, such as cURL-based scripts, with a modern and user-friendly CLI tool`,
	Run: func(cmd *cobra.Command, args []string) {
		
		// invoke help...
		cmd.HelpFunc()(cmd, []string{})
		
		os.Exit(0)
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		tracer.SetLevel(verbosity)		
		cc.NormalizePrefix(&kvPrefix, kdelim)
		return nil
	},
}

func init() {

	cc.SetupConsulFlags(RootCmd.PersistentFlags(),
		&consulConf, 
		&verbosity,
	)
	
	pf := RootCmd.PersistentFlags()
	pf.StringVarP(&kvPrefix, "prefix", "p", "", "Prefix for KV operations")
	pf.StringVarP(&kdelim, "delimiter", "d", "/", "Specify a key separator")
	
	RootCmd.AddCommand(cmdGet)
	RootCmd.AddCommand(cmdSet)
	RootCmd.AddCommand(cmdSetMulti)

	RootCmd.AddCommand(cmdList)
	RootCmd.AddCommand(cmdTree)

 	RootCmd.AddCommand(cmdRm)
	RootCmd.AddCommand(cmdRmTree)
}

func setupCommonFlags(ff *pflag.FlagSet) {
	ff.StringVar(&consulConf.Token, "token", "", "Provide Consul authorization token")
}

func keyName(x string) string {
	
	if ""==kvPrefix {
		return x
	}

	return kvPrefix+kdelim+x;
}
