package consul_common

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type ConsulService struct {
	Id	string
	Name string
	IP   net.IP
	Port uint16
	Tags []string
}


// implement Stringer{}
func (s *ConsulService) String() string {

	var serv_name string

	if "" != s.Id {
		serv_name = fmt.Sprintf("%s[%s]", s.Name, s.Id)
	} else {
		serv_name = s.Name
	}
	return fmt.Sprintf("ConsulService{%s} @%s:%d %s", serv_name, s.IP,s.Port, s.Tags)
}


func ParseTags(input string) ([]string, error) {
	
	if len(input) < 1 {
		return nil, errors.New("empty tags argument")
	}

	rt := make([]string, 0, 1)
	ts := strings.Split(input, ",")
	for _, v := range ts {
		rt = append(rt, v)
	}
	return rt, nil
}

func ParsePort(input string) (uint16,error) {
	
	if n, e := strconv.ParseUint(input, 10, 16); nil == e {
		return uint16(n & 0xFFFF), nil // just in case...
	} else {
		return 0,e
	}
}
