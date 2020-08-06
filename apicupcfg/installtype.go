package apicupcfg

import (
	"strings"
)

const InstallTypeOva = "ova"
const InstallTypeK8s = "k8s"
const InstallTypeInit = "init"
const InstallTypeUknown = "unknown"

type InstallTypeHeader struct {
	InstallType string
}

func InstallType(configFile string) string {
	h := InstallTypeHeader{InstallType: InstallTypeUknown}
	unmarshalJsonFile(configFile, &h)
	return strings.ToLower(strings.TrimSpace(h.InstallType))
}
