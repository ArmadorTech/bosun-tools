package commands

import (
	cc "../../common"
	// 	misc "../../misc"
	"fmt"
	"github.com/doblenet/go-doblenet/tracer"
	"github.com/spf13/cobra"
	consulapi "github.com/hashicorp/consul/api"
	"os"
)

// $ consul-sc --consul=consul.service.consul:8500 datacenter list

var cmdDatacenter = &cobra.Command{
	Use:   "datacenter",
	Short: "Datacenter-related information",
	Run:   datacenterRun,
}

var cmdDCLsLL bool
var cmdDCls = &cobra.Command{
	Use:     "ls",
	Short:   "List known datacenters",
	Aliases: []string{"list"},
	Run:     datacenterLsRun,
}

func init() {
	cf := cmdDatacenter.Flags()
	setupCommonFlags(cf)
	
	cmdDCls.Flags().BoolVarP(&cmdDCLsLL, "long", "l", false, "Request long listing (with extra information)")
	cmdDatacenter.AddCommand(cmdDCls)
}

func datacenterRun(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(0)
}

func datacenterLsRun(cmd *cobra.Command, args []string) {

	client, err := cc.ConsulClient(consulConf)
	if nil != err {
		tracer.FatalErr(err)
	}

	switch cmdDCLsLL {
		case false:
			if err := datacenterList(client); nil!=err {
				tracer.FatalErr(err)
			}
		case true:
			if err := datacenterXList(client); nil!=err {
				tracer.FatalErr(err)
			}
	}

	os.Exit(0)
}

func datacenterList(client *consulapi.Client) error {
	
	obj := client.Catalog()	
	
	result, err := obj.Datacenters()
	if nil != err {
		return err
	}
	
	for _, x := range result {
		fmt.Println(x)
	}
	
	return nil
}

func datacenterXList(client *consulapi.Client) error {
	
	obj := client.Coordinate()	
	
	result, err := obj.Datacenters()
	if nil != err {
		return err
	}
	
	for _, x := range result {
		
		fmt.Printf("%s\t", x.Datacenter)
		for i, c:= range x.Coordinates {
			if i>1 { fmt.Print("\t\t") }
			fmt.Println(cc.CEtoString(c))
		}
		
	}
	
	return nil
}