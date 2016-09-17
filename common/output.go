package consul_common

import (
	"bytes"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"encoding/json"
	"io"
	"strings"
)

func CEtoString(ce consulapi.CoordinateEntry) string {

	v := ce.Coord.Vec
	e := ce.Coord.Error
	a := ce.Coord.Adjustment
	h := ce.Coord.Height
	return fmt.Sprintf("@%s: {%v %g %f %f}", ce.Node, v, e, a, h)
}

func Map2String(o map[string]string) string {
	
	if 0 == len(o) {
		return ""
	}
	
	var t []string
	for k,v:= range o {
		t=append(t,fmt.Sprintf("%s:%s", k,v))
	}
	return "{"+strings.Join(t,", ")+"}"
}


func PrintPropMap(w io.Writer, x map[string]map[string]interface{}, lx bool) {

	if !lx {
		io.WriteString(w,propMap2str(&x))
	} else {
		b := propMap2BufJSON(&x)
		b.WriteTo(w)
	}
}
	
func propMap2str(x *map[string]map[string]interface{}) string {
	
	var s,t string
	for k,v := range *x {

		s += fmt.Sprintf("%s::\n", k)

		t = ""
		for pn,pv := range v {
			t += fmt.Sprintf("\t%s: %v\n", pn,pv)
		}
		s+=t
	}
	return s
}

func propMap2BufJSON(x *map[string]map[string]interface{}) bytes.Buffer {
 	var buf bytes.Buffer
 	
 	d,_ := json.Marshal(*x)
 	json.Indent(&buf, d, "", "  ")
	return buf
}
