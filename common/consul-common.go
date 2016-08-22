package consul_common

import (
	"errors"
	"fmt"
	"net"
)

// type Node struct {
//         Name    string
//         Address string
// }


func ResolveNode(nn string) (net.IP,error) {
	
	return net.IP{},errors.New("ENOTIMPL")
}
