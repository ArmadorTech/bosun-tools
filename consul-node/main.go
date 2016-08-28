package main

import (
	"github.com/doblenet/go-doblenet/cliapp"
	cmd "./commands"
)

const (
	gitCommit = "%GIT_COMMIT%"
)

var	appInfo = cliapp.AppInfo{
		Name: "consul-node",
		Version: "0.1.0",
		GitTag: "",
		Copyright: "(C)2016 DOBLENET Soluciones Tecnológicas, S.L.",
		Email: "consul-node@labs.doblenet.com",
	}

func main() {

	app := cliapp.New(appInfo)
	app.SetCommander(cmd.RootCmd)
	
	app.Run()
}
