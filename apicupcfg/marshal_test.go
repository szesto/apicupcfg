package apicupcfg

import (
	bytes2 "bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

type Y struct {
	Z2 string
	B2 string
}

type Z struct {
	A int
	B string
	C string
	E Y
	d int
}

func (z *Z) MarshalJSON() ([]byte, error) {

	var m map[string]interface{} = make(map[string]interface{})

	m["B"] = z.B
	m["E"] = z.E

	return json.Marshal(m)

	//z2 := Z{}
	//z2.B = z.B
	//z2.E.A2 = "hello"
	//z2.E.B2 = "world"
	//return json.Marshal(z2)
}

func (z *Z) MarshalText() (text []byte, err error) {
	s := fmt.Sprintf("{B: %s}", z.B)
	return []byte(s), nil
}

func TestMashal(t *testing.T) {

	z := &Z {
		A: 10,
		B: "b-hello",
		C: "c-world",
		d: 20,
	}

	//enc := json.NewEncoder(os.Stdout)
	//_ = enc.Encode(z)

	jb, _ := json.Marshal(z)
	var out bytes2.Buffer
	_ = json.Indent(&out, jb, "", "\t\t")
	_, _ = out.WriteTo(os.Stdout)
}
