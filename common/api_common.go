package consul_common

import (
	consulapi "github.com/hashicorp/consul/api"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	SSLenabled	= 0x0100
)

type ClientConfig struct {
	URL			string
	Datacenter	string

	directAPI	bool

	Token		string	// ACL token

// 	ConnFlags	uint
// 	SSLkey		string
// 	SSLcert		string
// 	SSLcacert	string
// 	TLSnoverify	bool

	AuthUser	string
	AuthPass	string
}




//! ENVIRONMENT (implemented by Consul api directly:
//	- CONSUL_HTTP_ADDR
//	- CONSUL_HTTP_TOKEN
//	- CONSUL_HTTP_AUTH -> user:password
//	- CONSUL_HTTP_SSL -> bool
//	- CONSUL_HTTP_SSL_VERIFY -> bool
func ConsulClient(params ClientConfig) (*consulapi.Client,error) {

	var err error
	
	config := consulapi.DefaultConfig()
	config.Address = params.URL
	
	if "" != params.Token {
		config.Token = params.Token
	}

	// setup HTTP basic auth
	config.HttpAuth, err = httpAuth(&params)
	if nil!=err {
		return nil,err
	}

	consul, e := consulapi.NewClient(config)
	if nil != e {
		return nil, e
	}

	// if passthrough requested, reinitialize the client
	if params.directAPI {
		
		// ....pointing to the cluster leader instead
		s,e := consul.Status().Leader()
		if nil!=e {
			return nil,e
		}
		
		config.Address = strings.Replace(s,":8300",":8500",1)
		consul,err = consulapi.NewClient(config)
		if nil != err {
			return nil,err
		}
	}

	return consul,nil
}

func httpAuth(params *ClientConfig) (*consulapi.HttpBasicAuth,error) {
	
	if ""==params.AuthUser && ""==params.AuthPass {
		return &consulapi.HttpBasicAuth{}, nil
	}
	if ""==params.AuthUser || ""==params.AuthPass {
		return nil, errors.New("AUTH: Both username and password need to be provided")
	}
	
	a := consulapi.HttpBasicAuth{
		Username: params.AuthUser,
		Password: params.AuthPass,
	}
	return &a,nil
}


func CheckServerError(e error) {
	if nil == e {
		return
	}
	
	if consulapi.IsServerError(e) {
		fmt.Fprintf(os.Stderr,"ERROR: Server error '%s' --- please try again later", e.Error())
		os.Exit(3)
	}
}
