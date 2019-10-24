package apicupcfg

import (
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"log"
)

type InitConfigObj struct {
	Mode string
	Version string
	Namespace string
}

func InitConfig(configFile string, installType string, tbox *rice.Box) {

	if ! (installType == InstallTypeOva || installType == InstallTypeK8s) {
		log.Fatalf("unsupported install type %s\n", installType)
	}

	// if config file already exists, complain
	exist, err := isFileExist(configFile)
	if err != nil {
		log.Fatalf("init-config... %v\n", err)
	}

	if exist {
		log.Fatalf("init-config... subsystem config file %s exists...\n", configFile)
	}

	configObj := InitConfigObj{}

	if installType == InstallTypeOva {
		writeTemplate(parseTemplate(tbox, tpdir(tbox) + "subsys-config-ova.tmpl"), configFile, configObj)
		fmt.Printf("initialized ova subsystem config file %s\n", configFile)

	} else if installType == InstallTypeK8s {
		writeTemplate(parseTemplate(tbox, tpdir(tbox) + "subsys-config-k8s.tmpl"), configFile, configObj)
		fmt.Printf("initialized k8s subsystem config file %s\n", configFile)
	}
}
