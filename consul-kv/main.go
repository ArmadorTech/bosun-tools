package main

import (
	"github.com/doblenet/go-doblenet/cliapp"
	cmd "./commands"
)

const (
	gitCommit = "%GIT_COMMIT%"
)

var	appInfo = cliapp.AppInfo{
		Name: "consul-kv",
		Version: "0.2.0beta",
		GitTag: "",
		Copyright: "(C)2016 DOBLENET Soluciones Tecnológicas, S.L.",
		Email: "consul-kv@labs.doblenet.com",
	}

func main() {

	app := cliapp.New(appInfo)
	app.SetMetadata(map[string]interface{}{
		"extra": "This Consul-KV has super-cow powers!",
	})
	app.SetCommander(cmd.RootCmd)
	
	app.Run()
}
