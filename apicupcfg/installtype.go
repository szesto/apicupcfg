package apicupcfg

import (
	"strings"
)

const InstallTypeOva = "ova"
const InstallTypeK8s = "k8s"
const InstallTypeUknown = "unknown"

type InstallTypeHeader struct {
	InstallType string
}

func InstallType(configFile string) string {
	h := InstallTypeHeader{InstallType: InstallTypeUknown}
	unmarshallJsonFile(configFile, &h)
	return strings.ToLower(strings.TrimSpace(h.InstallType))
}
