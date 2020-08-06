package apicupcfg

import "encoding/json"

type NetworkInterface struct {
	Device     string // eth0
	IpAddress  string
	SubnetMask string // dot notation
	Gateway    string

	// datapower host alias; if host-alias is set then it is used,
	//otherwise device name is formatted as datapower host alias
	HostAlias string
}

func (netif *NetworkInterface) gen(v10, verbose, gateway bool) {
	netif.Device = "eth0"
	netif.IpAddress = "ip4-address"
	netif.SubnetMask = "dot-mask"
	netif.Gateway = "gw-ip4-address"
}

type HostVm struct {
	Name string
	NetworkInterface
}

func (vm *HostVm) gen(v10, verbose, gateway bool) {
	vm.Name = "host-fqdn"
	vm.NetworkInterface.gen(v10, verbose, gateway)
}

type HostVmSubsys struct {
	HostVm
	HardDiskPassword string // luks storage encryption password
}

func (hvm *HostVmSubsys) gen(v10, verbose, gateway bool) {
	hvm.HostVm.gen(v10, verbose, gateway)
	hvm.HardDiskPassword = "luks-storage-encryption-password"
}

func (vm *HostVm) validateIp() {
	DecodeAddress(vm.IpAddress, vm.Gateway, vm.SubnetMask)
}

type IpRanges struct {
	PodNetwork string
	ServiceNetwork string
}

func (ipr *IpRanges) gen(v10, verbose bool) {
	ipr.PodNetwork = "172.16.0.0/16"
	ipr.ServiceNetwork = "172.17.0.0/16"
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

func (c *CloudInit) gen(v10, verbose bool) {
	c.CloudInitFile = "cloud-init.yaml"
	c.InitValues = make(map[string]interface{})
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
	Hosts []HostVmSubsys

	subsys bool
	gateway bool
}

func (fb *VmFirstBoot) MarshalJSON() ([]byte, error) {
	if fb.subsys {
		return json.Marshal(&struct{
			Hosts []HostVmSubsys
		}{fb.Hosts})

	} else {
		return json.Marshal(&struct{
			DnsServers []string
			VmwareConsolePasswordHash string
			IpRanges IpRanges
			Hosts []HostVmSubsys
		}{fb.DnsServers, fb.VmwareConsolePasswordHash, fb.IpRanges, fb.Hosts})
	}
}

func (fb *VmFirstBoot) gen(v10, verbose, subsys, gateway bool) {
	fb.DnsServers = []string {"dns.ip.addr.1", "dns.ip.addr.2"}
	fb.VmwareConsolePasswordHash = "vm-console-hash"
	fb.IpRanges.gen(v10, verbose)

	// 3 hosts
	if subsys {
		fb.Hosts = make([]HostVmSubsys, 3)
		for j := range []int{0,1,2} {
			fb.Hosts[j].gen(v10, verbose, gateway)
		}
	}

	fb.subsys = subsys
	fb.gateway = gateway
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
		fb.Hosts = []HostVmSubsys{}
	}
}

type SubsysVmBase struct {
	SubsysName string
	Mode string
	v10 bool	// private field, not serialized to json

	CloudInit CloudInit
	SearchDomains []string
	VmFirstBoot VmFirstBoot
	SshPublicKeyFile string
}

//func (b *SubsysVmBase) MarshalJSON() ([]byte, error) {
//	return json.Marshal(&struct{VmFirstBoot VmFirstBoot}{b.VmFirstBoot})
//}

func (b *SubsysVmBase) gen(v10, verbose, subsys, gateway bool) {
	b.SubsysName = "mgmt"

	if v10 {
		b.Mode = "Production"
	} else {
		b.Mode = "standard"
	}

	b.v10 = v10
	b.CloudInit.gen(v10, verbose)
	b.SearchDomains = []string{"my.domain.com","domain.com"}

	// subsys = true
	b.VmFirstBoot.gen(v10, verbose, subsys, gateway)

	// subsystem public key not generated, can be set manually
}

func (b *SubsysVmBase) copyDefaults(from SubsysVm) {
	// copy v10 flag
	b.v10 = from.V10

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

