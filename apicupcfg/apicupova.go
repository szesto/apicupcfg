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

	CassandraEncryptionKeyFile string

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

type GwySubsysVm struct {
	SubsysName string
	Mode string

	SearchDomains []string
	DnsServers []string
	Hosts []HostVm

	ApiGateway string
	ApicGwService string

	// apic datapower domain
	DatapowerDomain string

	// apic configuration sequence (low level)
	ConfigurationSequenceName string
	ConfigurationExecutionInterval int

	// API gateway
	DatapowerApiGatewayPort int
	DatapowerApiGatewayAddress string // host-alias

	// API connect gateway service
	DatapowerApicGwServicePort int // 3000
	DatapowerApicGwServiceAddress string // host-alias

	// NTP server
	NTPServer string

	// low level peering configuration
	GwdPeering string
	GwdPeeringLocalPort int
	GwdPeeringMonitorPort int

	RateLimitPeering string
	RateLimitPeeringLocalPort int
	RateLimitPeeringMonitorPort int

	SubsPeering string
	SubsPeeringLocalPort int
	SubsPeeringMonitorPort int

	// datapower cert configuration
	DatapowerCryptoDir string

	// gwd_key
	DatapowerGwdKey string

	// gwd_cert
	DatapowerGwdCert string

	// ca, root-ca
	CaFile string
	RootCaFile string

	// apic_ca_cert, apic_root_cert
	DatapowerCaCert string
	DatapowerRootCert string
}

func (gwy *GwySubsysVm) GetCryptoDirectoryOrDefault() string {
	if len(gwy.DatapowerCryptoDir) == 0 {
		return cryptoDir
	}
	return gwy.DatapowerCryptoDir
}

func (gwy *GwySubsysVm) GetGwdKeyOrDefault() string {
	if len(gwy.DatapowerGwdKey) == 0 {
		return gwdKey
	}
	return gwy.DatapowerGwdKey
}

func (gwy *GwySubsysVm) GetGwdCertOrDefault() string {
	if len(gwy.DatapowerGwdCert) == 0 {
		return gwdCert
	}
	return gwy.DatapowerGwdCert
}

func (gwy *GwySubsysVm) GetCaCertOrDefault() string {
	if len(gwy.DatapowerCaCert) == 0 {
		return "gwd_ca_cert"
	}
	return gwy.DatapowerCaCert
}

func (gwy *GwySubsysVm) GetRootCertOrDefault() string {
	if len(gwy.DatapowerRootCert) == 0 {
		return "gwd_root_cert"
	}
	return gwy.DatapowerRootCert
}

func (gwy *GwySubsysVm) GetDatapowerDomainOrDefault() string {
	if len(gwy.DatapowerDomain) == 0 {
		return "apiconnect"
	}
	return gwy.DatapowerDomain
}

func (gwy *GwySubsysVm) GetNTPServerOrDefault() string {
	if len(gwy.NTPServer) == 0 {
		return "pool.ntp.org"
	}
	return gwy.NTPServer
}

func (gwy *GwySubsysVm) GetGwdPeeringOrDefault() string {
	if len(gwy.GwdPeering) == 0 {
		return gwdPeering
	}
	return gwy.GwdPeering
}

func (gwy *GwySubsysVm) GetGwdPeeringLocalPortOrDefault() int {
	if gwy.GwdPeeringLocalPort == 0 {
		return gwdPeeringLocalPort
	}
	return gwy.GwdPeeringLocalPort
}

func (gwy *GwySubsysVm) GetGwdPeeringMonitorPortOrDefault() int {
	if gwy.GwdPeeringMonitorPort == 0 {
		return gwdPeeringMonitorPort
	}
	return gwy.GwdPeeringMonitorPort
}

func (gwy *GwySubsysVm) GetRateLimitPeeringOrDefault() string {
	if len(gwy.RateLimitPeering) == 0 {
		return rateLimitPeering
	}
	return gwy.RateLimitPeering
}

func (gwy *GwySubsysVm) GetRateLimitPeeringLocalPortOrDefault() int {
	if gwy.RateLimitPeeringLocalPort == 0 {
		return rateLimitPeeringLocalPort
	}
	return gwy.RateLimitPeeringLocalPort
}

func (gwy *GwySubsysVm) GetRateLimitPeeringMonitorPortOrDefault() int {
	if gwy.RateLimitPeeringMonitorPort == 0 {
		return rateLimitPeeringMonitorPort
	}
	return gwy.RateLimitPeeringMonitorPort
}

func (gwy *GwySubsysVm) GetSubsPeeringOrDefault() string {
	if len(gwy.SubsPeering) == 0 {
		return subsPeering
	}
	return gwy.SubsPeering
}

func (gwy *GwySubsysVm) GetSubsPeeringLocalPortOrDefault() int {
	if gwy.SubsPeeringLocalPort == 0 {
		return subsPeeringLocalPort
	}
	return gwy.SubsPeeringLocalPort
}

func (gwy *GwySubsysVm) GetSubsPeeringMonitorPortOrDefault() int {
	if gwy.SubsPeeringMonitorPort == 0 {
		return subsPeeringMonitorPort
	}
	return gwy.SubsPeeringMonitorPort
}

type SubsysVm struct {
	InstallTypeHeader

	Version string
	Tag string

	UseVersion bool // use version for the apic executable

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
	Gateway GwySubsysVm

	OsEnv

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

func LoadSubsysVm(jsonConfigFile string) *SubsysVm {

	// unmarshal values file
	subsys := &SubsysVm{}
	unmarshallJsonFile(jsonConfigFile, &subsys)

	subsys.OsEnv.init2(subsys.Version, subsys.UseVersion)

	// save input config file
	subsys.configFileName = jsonConfigFile

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

	// execute templates
	shellext := subsys.OsEnv.ShellExt

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

func CopyCertVm(certfile string, isdir bool, subsys *SubsysVm, commonCsrDir string, customCsrDir string) error {

	if isdir {
		return copyCerts(certfile, &subsys.Certs, &subsys.Management, &subsys.Analytics,
			&subsys.Portal, &subsys.Gateway, commonCsrDir, customCsrDir, true)

	} else {
		return copyCert(certfile, &subsys.Certs, &subsys.Management, &subsys.Analytics,
			&subsys.Portal, &subsys.Gateway, commonCsrDir, customCsrDir, true)
	}
}

func SomaUpload(subsys *SubsysVm, uploadfile, dpdir, dpfile, dpdomain, dpenv, url string, tbox *rice.Box) (status string, statusCode int, reply string, err error) {

	return SomaUploadFile(uploadfile, dpdir, dpfile, dpdomain, dpenv, url, tbox)
}