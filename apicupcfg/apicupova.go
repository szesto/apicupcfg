package apicupcfg

import (
	"bytes"
	"encoding/json"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"os"
	"text/template"
)

type SubsysVm struct {
	InstallTypeHeader

	Version string
	Tag string
	V10 bool

	UseVersion bool // use version for the apic executable
	//Passive bool // passive site deployment, import crypto from active site, replace with byok in certs

	// defaults
	Mode string // dev|standard, v10: Production|Nonproduction
	SshPublicKeyFile string
	SearchDomains[] string
	VmFirstBoot VmFirstBoot
	CloudInit CloudInit

	Certs Certs

	Management MgtSubsysVm
	Analytics  AltSubsysVm
	Portal PtlSubsysVm
	Gateway GwySubsysVm

	//OsEnv

	// config file from which this config was loaded
	configFileName string
}

func ValidateHostIpVm(subsys *SubsysVm) {

	isManagement := len(subsys.Management.SubsysName) > 0
	isAnalytics := len(subsys.Analytics.SubsysName) > 0
	isPortal := len(subsys.Portal.SubsysName) > 0
	isGateway := len(subsys.Gateway.SubsysName) > 0

	if isManagement {
		fmt.Printf("\n--- ip check for the management subsystem \"%s\"\n", subsys.Management.SubsysName)
		for _, hostvm := range subsys.Management.VmFirstBoot.Hosts {
			fmt.Printf("\n-host: %s\n", hostvm.Name)
			hostvm.validateIp()
		}
	}

	if isAnalytics {
		fmt.Printf("\n--- ip check for the analytics subsystem \"%s\"\n", subsys.Analytics.SubsysName)
		for _, hostvm := range subsys.Analytics.VmFirstBoot.Hosts {
			fmt.Printf("\n-host: %s\n", hostvm.Name)
			hostvm.validateIp()
		}
	}

	if isPortal {
		fmt.Printf("\n--- ip check for the portal subsystem \"%s\"\n", subsys.Portal.SubsysName)
		for _, hostvm := range subsys.Portal.VmFirstBoot.Hosts {
			fmt.Printf("\n-host: %s\n", hostvm.Name)
			hostvm.validateIp()
		}
	}

	if isGateway {
		fmt.Printf("\n--- ip check for the gateway subsystem \"%s\"\n", subsys.Gateway.SubsysName)
		for _, hostvm := range subsys.Gateway.Hosts {
			fmt.Printf("\n-host: %s\n", hostvm.Name)
			hostvm.validateIp()
		}
	}
}

func GenSubsysVm(jsonConfigFile string, v10 bool, verbose bool) error {
	// gen subsys
	subsys := &SubsysVm{}
	subsys.gen(v10, verbose)

	// json encode subsys
	b, err := json.Marshal(subsys)
	if err != nil {
		return err
	}

	var out bytes.Buffer
	_ = json.Indent(&out, b, "", "\t\t")
	_, _ = out.WriteTo(os.Stdout)

	return nil
}

func (subsys *SubsysVm) gen(v10, verbose bool) {
	subsys.genDefaults(v10, verbose)

	//subsys.Certs.gen();

	subsys.Management.gen(v10, verbose);
	subsys.Analytics.gen(v10, verbose);
	subsys.Portal.gen(v10, verbose);
	subsys.Gateway.gen(v10, verbose);
}

func (subsys *SubsysVm) genDefaults(v10, verbose bool) {
	subsys.InstallType = InstallTypeOva

	if v10 {
		subsys.Version = "linux_lts_v10"
	} else {
		subsys.Version = "linux_lts_v2018.4.1.12"
	}

	subsys.Tag = "tag"
	subsys.V10 = v10
	subsys.UseVersion = true

	if v10 {
		subsys.Mode = "Production"
	} else {
		subsys.Mode = "standard"
	}

	subsys.SshPublicKeyFile = "/path/to/public-key-file"
	subsys.SearchDomains = []string {"my.domain.com", "domain.com"}

	subsys.VmFirstBoot.gen(v10, verbose, false, false)
	subsys.CloudInit.gen(v10, verbose)
}

func LoadSubsysVm(jsonConfigFile string) *SubsysVm {

	// unmarshal values file
	subsys := &SubsysVm{}
	unmarshalJsonFile(jsonConfigFile, &subsys)

	// save input config file
	subsys.configFileName = jsonConfigFile

	isManagement := len(subsys.Management.SubsysName) > 0
	isAnalytics := len(subsys.Analytics.SubsysName) > 0
	isPortal := len(subsys.Portal.SubsysName) > 0
	isGateway := len(subsys.Gateway.SubsysName) > 0

	// copy defaults
	if isManagement {
		subsys.Management.copyDefaults(*subsys)
	}

	if isAnalytics {
		subsys.Analytics.copyDefaults(*subsys)
	}

	if isPortal {
		subsys.Portal.copyDefaults(*subsys)
	}

	if isGateway {
		// @todo: copy defaults
	}

	// copy certs defaults
	byok := false
	subsys.Certs.Passive = byok // subsys.Passive

	return subsys
}

func ApplyTemplateVm(subsys *SubsysVm, outfiles map[string]string, subsysOnly, certsOnly, datapowerOnly bool, tbox *rice.Box) {

	isManagement := len(subsys.Management.SubsysName) > 0
	isAnalytics := len(subsys.Analytics.SubsysName) > 0
	isPortal := len(subsys.Portal.SubsysName) > 0
	isCloudInit := len(subsys.CloudInit.CloudInitFile) > 0

	// parse templates
	var mgtt, analyt, ptl, cloudinitt *template.Template

	if isManagement {
		mgtt = parseTemplates(tbox, tpdir(tbox) + "management-vm.tmpl", tpdir(tbox) + "helpers.tmpl")
		//mgtt = parseTemplates("templates/management-vm.tmpl", "templates/helpers.tmpl")
	}

	if isAnalytics {
		analyt = parseTemplates(tbox, tpdir(tbox) + "analytics-vm.tmpl", tpdir(tbox) + "helpers.tmpl")
	}

	if isPortal {
		ptl = parseTemplates(tbox, tpdir(tbox) + "portal-vm.tmpl", tpdir(tbox) + "helpers.tmpl")
	}

	if isCloudInit {
		cloudinitt = parseTemplates(tbox, tpdir(tbox) + "cloud-init-vm.tmpl", tpdir(tbox) + "helpers.tmpl")
	}

	var osenv OsEnv
	osenv.init2(subsys.Version, subsys.UseVersion)

	// execute templates
	//shellext := subsys.OsEnv.ShellExt
	shellext := osenv.ShellExt

	var outpath string

	oneof := func(a, b bool) bool { if a {return a} else if b {return b}; return false }

	if isManagement && !oneof(certsOnly, datapowerOnly) {
		outpath = fileName(outfiles[outdir], tagOutputFileName(outfiles[managementOut], subsys.Tag)) + shellext
		writeTemplate(mgtt, outpath, subsys.Management)
	}

	if isAnalytics && !oneof(certsOnly, datapowerOnly) {
		outpath = fileName(outfiles[outdir], tagOutputFileName(outfiles[analyticsOut], subsys.Tag)) + shellext
		writeTemplate(analyt, outpath, subsys.Analytics)
	}

	if isPortal && !oneof(certsOnly, datapowerOnly) {
		outpath = fileName(outfiles[outdir], tagOutputFileName(outfiles[portalOut], subsys.Tag)) + shellext
		writeTemplate(ptl, outpath, subsys.Portal)
	}

	// this outputs default cloud-init file... each subsystem can have it's own
	if isCloudInit && !oneof(certsOnly, datapowerOnly) {
		outpath = fileName(outfiles[outdir], subsys.CloudInit.CloudInitFile)
		writeTemplate(cloudinitt, outpath, subsys.CloudInit.InitValues)
	}

	// certs
	if  !oneof(subsysOnly, datapowerOnly) {
		updateCertSpecs(&subsys.Certs, &subsys.Management, &subsys.Analytics, &subsys.Portal, &subsys.Gateway, outfiles[CommonCsrOutDir], outfiles[CustomCsrOutDir])
		outputCerts(&subsys.Certs, outfiles, subsys.Tag, subsys.Version, subsys.UseVersion,  tbox)
	}

	// datapower
	if !oneof(subsysOnly, certsOnly) {
		datapowerCluster(subsys, outfiles, tbox)
	}
}

//func CopyCertVm(certfile string, isdir bool, subsys *SubsysVm, commonCsrDir string, customCsrDir string) error {
//
//	if isdir {
//		return copyCerts(certfile, &subsys.Certs, &subsys.Management, &subsys.Analytics,
//			&subsys.Portal, &subsys.Gateway, commonCsrDir, customCsrDir, true)
//
//	} else {
//		return copyCert(certfile, &subsys.Certs, &subsys.Management, &subsys.Analytics,
//			&subsys.Portal, &subsys.Gateway, commonCsrDir, customCsrDir, true)
//	}
//}

func SomaUpload(uploadfile, dpdir, dpfile, dpdomain, dpenv, url string, tbox *rice.Box) (status string, statusCode int, reply string, err error) {

	return SomaUploadFile(uploadfile, dpdir, dpfile, dpdomain, dpenv, url, tbox)
}