package main

import (
	cmd "./commands"
	"./misc"
	"fmt"
	"os"
)

func main() {

	fmt.Println("Running at ", misc.LocalHostname(), "-- ", misc.LocalIPs())

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}
