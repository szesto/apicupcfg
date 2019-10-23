package apicupcfg

import (
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"text/template"
)

type HostVm struct {
	Name string
	HardDiskPassword string // luks storage encryption password
	Device string // eth0
	IpAddress string
	SubnetMask string // dot notation
	Gateway string
}

func (vm *HostVm) validateIp() {
	DecodeAddress(vm.IpAddress, vm.Gateway, vm.SubnetMask)
}

type IpRanges struct {
	PodNetwork string
	ServiceNetwork string
}

func (ipr *IpRanges) copyDefaults(from IpRanges) {
	if len(ipr.PodNetwork) == 0 {
		ipr.PodNetwork = from.PodNetwork
	}

	if len(ipr.ServiceNetwork) == 0 {
		ipr.ServiceNetwork = from.ServiceNetwork
	}
}

type CloudInit struct {
	CloudInitFile string
	InitValues map[string]interface {}
}

func (c *CloudInit) copyDefaults(from CloudInit) {

	if len(c.CloudInitFile) == 0 {
		c.CloudInitFile = from.CloudInitFile
	}

	// todo: do not overwrite values...
	c.InitValues = from.InitValues

	if c.InitValues == nil {
		c.InitValues = make(map[string]interface {})
	}
}

type VmFirstBoot struct {
	DnsServers []string
	VmwareConsolePasswordHash string
	IpRanges IpRanges
	Hosts []HostVm
}

func (fb *VmFirstBoot) copyDefaults(from VmFirstBoot) {

	// copy dns servers
	fb.DnsServers = copySlices(fb.DnsServers, from.DnsServers)

	// copy hash password
	if len(fb.VmwareConsolePasswordHash) == 0 {
		fb.VmwareConsolePasswordHash = from.VmwareConsolePasswordHash
	}

	// copy ip ranges
	fb.IpRanges.copyDefaults(from.IpRanges)

	// do not copy hosts
	if fb.Hosts == nil {
		fb.Hosts = []HostVm{}
	}
}

type SubsysVmBase struct {
	SubsysName string
	Mode string

	CloudInit CloudInit
	SearchDomains []string
	VmFirstBoot VmFirstBoot
	SshPublicKeyFile string

	OsEnv
}

func (b *SubsysVmBase) copyDefaults(from SubsysVm) {
	// copy osenv
	b.OsEnv.copyDefaults(from.OsEnv)

	// copy mode
	if len(b.Mode) == 0 {
		b.Mode = from.Mode
	}

	// copy cloud-init
	b.CloudInit.copyDefaults(from.CloudInit)

	// copy search domains
	b.SearchDomains = copySlices(b.SearchDomains, from.SearchDomains)

	// copy first-boot
	b.VmFirstBoot.copyDefaults(from.VmFirstBoot)

	// copy ssh-public-key
	if len(b.SshPublicKeyFile) == 0 {
		b.SshPublicKeyFile = from.SshPublicKeyFile
	}
}

type MgtSubsysVm struct {
	SubsysVmBase

	CassandraBackup Backup

	PlatformApi string
	ApiManagerUi string
	CloudAdminUi string
	ConsumerApi string
}

type AltSubsysVm struct {
	SubsysVmBase

	AnalyticsIngestion string
	AnalyticsClient string

	EnableMessageQueue bool
}

type PtlSubsysVm struct {
	SubsysVmBase

	SiteBackup Backup

	PortalAdmin string
	PortalWww string
}

type SubsysVm struct {
	InstallTypeHeader

	// defaults
	Mode string // dev|standard
	SshPublicKeyFile string
	SearchDomains[] string
	VmFirstBoot VmFirstBoot
	CloudInit CloudInit

	Certs Certs

	Management MgtSubsysVm
	Analytics  AltSubsysVm
	Portal PtlSubsysVm

	OsEnv
}

func ValidateHostIpVm(subsys *SubsysVm) {

	isManagement := len(subsys.Management.SubsysName) > 0
	isAnalytics := len(subsys.Analytics.SubsysName) > 0
	isPortal := len(subsys.Portal.SubsysName) > 0

	if isManagement {
		fmt.Printf("\n--- ip check for the management subsystem \"%s\"\n", subsys.Management.SubsysName)
		for _, hostvm := range subsys.Management.VmFirstBoot.Hosts {
			fmt.Printf("-host: %s\n", hostvm.Name)
			hostvm.validateIp()
		}
	}

	if isAnalytics {
		fmt.Printf("\n--- ip check for the analytics subsystem \"%s\"\n", subsys.Analytics.SubsysName)
		for _, hostvm := range subsys.Analytics.VmFirstBoot.Hosts {
			fmt.Printf("-host: %s\n", hostvm.Name)
			hostvm.validateIp()
		}
	}

	if isPortal {
		fmt.Printf("\n--- ip check for the portal subsystem \"%s\"\n", subsys.Portal.SubsysName)
		for _, hostvm := range subsys.Portal.VmFirstBoot.Hosts {
			fmt.Printf("-host: %s\n", hostvm.Name)
			hostvm.validateIp()
		}
	}
}

func LoadSubsysVm(jsonConfigFile string) *SubsysVm {

	// unmarshal values file
	subsys := &SubsysVm{}
	subsys.OsEnv.init()

	unmarshallJsonFile(jsonConfigFile, &subsys)

	isManagement := len(subsys.Management.SubsysName) > 0
	isAnalytics := len(subsys.Analytics.SubsysName) > 0
	isPortal := len(subsys.Portal.SubsysName) > 0

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

	return subsys
}

func ApplyTemplateVm(subsys *SubsysVm, outfiles map[string]string, tbox *rice.Box) {

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

	// execute templates
	shellext := subsys.OsEnv.ShellExt

	var outpath string

	if isManagement {
		outpath = fileName(outfiles[outdir], outfiles[managementOut]) + shellext
		writeTemplate(mgtt, outpath, subsys.Management)
	}

	if isAnalytics {
		outpath = fileName(outfiles[outdir], outfiles[analyticsOut]) + shellext
		writeTemplate(analyt, outpath, subsys.Analytics)
	}

	if isPortal {
		outpath = fileName(outfiles[outdir], outfiles[portalOut]) + shellext
		writeTemplate(ptl, outpath, subsys.Portal)
	}

	// this outputs default cloud-init file... each subsystem can have it's own
	if isCloudInit {
		outpath = fileName(outfiles[outdir], subsys.CloudInit.CloudInitFile)
		writeTemplate(cloudinitt, outpath, subsys.CloudInit.InitValues)
	}

	// update cert specs
	gwySubsysName := "gwy" // no gwy subystem in vm deployment

	certs := updateCertSpecs(subsys.Certs, subsys.Management.SubsysName,
		subsys.Analytics.SubsysName, subsys.Portal.SubsysName, gwySubsysName,
		outfiles[commonCsrOutDir], outfiles[customCsrOutDir])

	outputCerts(&certs, outfiles, tbox)
}
