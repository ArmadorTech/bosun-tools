package consul_common

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
)


// type KVPair struct {
//     Key         string
//     CreateIndex uint64
//     ModifyIndex uint64
//     LockIndex   uint64
//     Flags       uint64
//     Value       []byte
//     Session     string
// }

func KVPairLL(p *consulapi.KVPair, ll, lx bool) string {

	if !ll && !lx { 
		return fmt.Sprintf("%s = %s", p.Key, string(p.Value))
	}
	
	if !lx {
		return fmt.Sprintf("%s\t\tflags=%d", p.Key, p.Flags)
	}
	
	var v = fmt.Sprintf("%s\t[Create=%d, Modify=%d", p.Key, 
					p.CreateIndex,p.ModifyIndex)
	
	if 0!=p.LockIndex {
		v += fmt.Sprintf("; Lock=%d; Session=%s",
					p.LockIndex, p.Session)
	}
	v += "]";
	
	return v
}

func KVPair2String(p *consulapi.KVPair) string {

	v := fmt.Sprintf("%s\t\tflags=%d", p.Key, p.Flags)
	
	v += fmt.Sprintf("\t[Create=%d, Modify=%d", 
					p.CreateIndex,p.ModifyIndex)
	
	if 0!=p.LockIndex {
		v += fmt.Sprintf("; Lock=%d; Session=%s",
					p.LockIndex, p.Session)
	}
	v += "]";
	
	return v
}