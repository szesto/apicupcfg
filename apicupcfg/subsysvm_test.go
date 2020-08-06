package apicupcfg

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func writeout(s interface{}) {
	jb, _ := json.Marshal(s)
	var out bytes.Buffer
	_ = json.Indent(&out, jb, "", "\t")
	_, _ = out.WriteTo(os.Stdout)
}

func Test_GenSubsysVm(t *testing.T) {

	s := &SubsysVm{}

	v10 := true
	verbose := false
	s.genDefaults(v10, verbose)

	s.Management.gen(v10, verbose)
	s.Analytics.gen(v10, verbose)
	s.Portal.gen(v10, verbose)
	s.Gateway.gen(v10, verbose)

	writeout(s)
}
