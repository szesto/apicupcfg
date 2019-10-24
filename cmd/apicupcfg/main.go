package main

import (
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/szesto/apicupcfg/apicupcfg"
	"log"
)

func main() {
	// template embedding box
	tbox := rice.MustFindBox("../../templates")

	// input: configuration file, output dir, csr subdirectories
	input, outdir, commonCsrSubdir, customCsrSubdir, projectSubdir, validateIp, initConfig, initConfigType := apicupcfg.Input()

	// output files
	output := apicupcfg.OutputFiles(outdir, commonCsrSubdir, customCsrSubdir)

	// install type
	installType := apicupcfg.InstallTypeUknown

	if initConfig {
		installType = apicupcfg.InstallTypeInit

	} else {
		installType = apicupcfg.InstallType(input)
	}

	switch installType {
	case apicupcfg.InstallTypeOva:
		// load subsystem
		subsysvm := apicupcfg.LoadSubsysVm(input)

		if validateIp {
			apicupcfg.ValidateHostIpVm(subsysvm)

		} else {
			// create output directories
			err := apicupcfg.CreateOutputDirectories(outdir, commonCsrSubdir, customCsrSubdir, projectSubdir)
			if err != nil {
				log.Fatal(err)
			}

			// apply templates
			apicupcfg.ApplyTemplateVm(subsysvm, output, tbox)
		}

	case apicupcfg.InstallTypeK8s:
		// load subsystem
		subsysk8s := apicupcfg.LoadSubsysK8s(input)

		if validateIp {
			// not applicable, complain
			fmt.Printf("validateip command line option is not applicable to the %s install type...\n", apicupcfg.InstallTypeK8s)

		} else {
			// create output directories
			err := apicupcfg.CreateOutputDirectories(outdir, commonCsrSubdir, customCsrSubdir, projectSubdir)
			if err != nil {
				log.Fatal(err)
			}

			// apply templates
			apicupcfg.ApplyTemplatesK8s(subsysk8s, output, tbox)
		}

	case apicupcfg.InstallTypeInit:
		apicupcfg.InitConfig(input, initConfigType, tbox)

	default:
		log.Fatalf("unsupported install type %s\n", installType)
	}
}
