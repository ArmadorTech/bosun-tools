package consul_common

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
)

func CEtoString(ce consulapi.CoordinateEntry) string {

	v := ce.Coord.Vec
	e := ce.Coord.Error
	a := ce.Coord.Adjustment
	h := ce.Coord.Height
	return fmt.Sprintf("@%s: {%v %g %f %f}", ce.Node, v, e, a, h)
}
