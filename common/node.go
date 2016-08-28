package consul_common

import (
	"errors"
	"net"
)



func ParseNode(x string, predicate func(nodename string) bool) (string, bool) {
	if len(x) < 2 || '@'!=x[1] {
		// not a nodename ...
		return "", false
	}
	nodename := x[1:]
	if nil!=predicate {
		return nodename, predicate(nodename)
	}
	return nodename, true
}

func ResolveNode(nn string) (net.IP,error) {
	
	return net.IP{},errors.New("ENOTIMPL")
}
