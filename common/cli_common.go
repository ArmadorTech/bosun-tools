package consul_common

import (
	"github.com/spf13/pflag"
	"errors"
	"fmt"
	"os"
)

const (
	k_ENV_CONSUL = "CONSUL_HOST"
	k_ENV_DC     = "CONSUL_DC"

	k_CONSUL_URL = "localhost:8500"
	k_CONSUL_DC  = "dc1"
)

// ConnFlags	uint
// 	SSLkey		string
// 	SSLcert		string
// 	SSLcacert	string
// 	SSLnoverify	bool
// 


func SetupConsulFlags(ff *pflag.FlagSet, consulConf *ClientConfig, vv *int) error {
	
	ff.CountVarP(vv, "verbose", "v", "Enable (and/or increase) verbosity level")
	ff.BoolVarP(&consulConf.directAPI, "passthrough", "z", false, "Enable direct operation against a Consul server, bypassing the (normally local) agent")

	ff.StringVar(&consulConf.URL, "consul", "", "Consul HTTP API endpoint to use (default: $CONSUL_HOST)")
	ff.StringVar(&consulConf.Datacenter, "dc", "", "Consul datacenter identifier to use")

	var url, dc string
	if url = os.Getenv(k_ENV_CONSUL); "" == url {
		url = k_CONSUL_URL
	}
	if dc = os.Getenv(k_ENV_DC); "" == dc {
		dc = k_CONSUL_DC
	}
	consulConf.URL = url
	consulConf.Datacenter = dc
	
	// Extended flags: HTTP basic auth & SSL
	ff.StringVar(&consulConf.AuthUser, "user", "u", "Consul (HTTP basic auth) Username")
	ff.StringVar(&consulConf.AuthPass, "pass", "", "Consul (HTTP basic auth) Password")
	
	
	
	
	return nil
}

func ValidateConfig(cf *ClientConfig) error {
	
	if ""!=cf.SSLkey && ""!=cf.SSLcert {
		cf.ConnFlags = SSLenabled
	}
	
	if ""==cf.SSLcacert && !cf.TLSnoverify {
		return errors.New("No CAcert provided, and strict TLS verification requested")
	}
	
	return nil
}


func PrintConfig(consulConf *ClientConfig) {
	
	fmt.Println("direct->", consulConf.directAPI)
	fmt.Println("consul->", consulConf.URL)
	fmt.Println("datacenter->", consulConf.Datacenter)	
}
