package misc

import (
	"bytes"
	"fmt"
	"encoding/json"
	"io"
)

var (
	jsonPrefix = ""
	jsonIndent= "    "
)



func PrettyJSON(w io.Writer, d []byte) {
	//b, err := json.Marshal(roads)
	var buf bytes.Buffer
	
	json.Indent(&buf, d, jsonPrefix, jsonIndent)
	buf.WriteTo(w)
	fmt.Fprintf(w,"\n")
}

func OutputJSON(w io.Writer, d []byte) {
	var buf bytes.Buffer
	json.Compact(&buf, d)
	buf.WriteTo(w)
}
