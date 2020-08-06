package apicupcfg

import "encoding/json"

type GatewayRoute struct {
	Destination string
	NextHopRouter string
	Metric string
}

type GatewayInterface struct {
	NetworkInterface
	Routes []GatewayRoute
}

type HostGateway struct {
	HostVm

	Interfaces []GatewayInterface // advanced datapower configuration

	GwdPeeringPriority int
	GwdPeeringInterface string

	RateLimitPeeringPriority int
	RateLimitPeeringInterface string

	SubsPeeringPriority int
	SubsPeeringInterface string

	ApiProbePeeringPriority int
	ApiProbePeeringInterface string

	v10 bool
	subsysgen bool
}

func (h *HostGateway) MarshalJSON() ([]byte, error) {
	if h.v10 {
		if h.subsysgen {
			return json.Marshal(&struct{
				Name string
				Device string
				IpAddress string
				SubnetMask string
				Gateway string

				GwdPeeringPriority int
				RateLimitPeeringPriority int
				SubsPeeringPriority int
				ApiProbePeeringPriority int
			}{
				Name: h.Name, Device: h.Device, IpAddress: h.IpAddress,
				SubnetMask: h.SubnetMask, Gateway: h.Gateway,
				GwdPeeringPriority: h.GwdPeeringPriority, RateLimitPeeringPriority: h.RateLimitPeeringPriority,
				SubsPeeringPriority: h.SubsPeeringPriority, ApiProbePeeringPriority: h.ApiProbePeeringPriority,
			})

		} else {
			return json.Marshal(&struct{
				Name string
				Device string
				IpAddress string
				SubnetMask string
				Gateway string
				HostAlias string

				GwdPeeringPriority int
				RateLimitPeeringPriority int
				SubsPeeringPriority int
				ApiProbePeeringPriority int
			}{
				Name: h.Name, Device: h.Device, IpAddress: h.IpAddress,
				SubnetMask: h.SubnetMask, Gateway: h.Gateway, HostAlias: h.HostAlias,
				GwdPeeringPriority: h.GwdPeeringPriority, RateLimitPeeringPriority: h.RateLimitPeeringPriority,
				SubsPeeringPriority: h.SubsPeeringPriority, ApiProbePeeringPriority: h.ApiProbePeeringPriority,
			})
		}

	} else if h.subsysgen {
		return json.Marshal(&struct{
			Name string
			Device string
			IpAddress string
			SubnetMask string
			Gateway string

			GwdPeeringPriority int
			RateLimitPeeringPriority int
			SubsPeeringPriority int
			ApiProbePeeringPriority int
		}{
			Name: h.Name, Device: h.Device, IpAddress: h.IpAddress,
			SubnetMask: h.SubnetMask, Gateway: h.Gateway,
			GwdPeeringPriority: h.GwdPeeringPriority, RateLimitPeeringPriority: h.RateLimitPeeringPriority,
			SubsPeeringPriority: h.SubsPeeringPriority, ApiProbePeeringPriority: h.ApiProbePeeringPriority,
		})

	} else {
		return json.Marshal(&struct{
			Name string
			Device string
			IpAddress string
			SubnetMask string
			Gateway string
			HostAlias string

			GwdPeeringPriority int
			RateLimitPeeringPriority int
			SubsPeeringPriority int
			ApiProbePeeringPriority int
		}{
			Name: h.Name, Device: h.Device, IpAddress: h.IpAddress,
			SubnetMask: h.SubnetMask, Gateway: h.Gateway, HostAlias: h.HostAlias,
			GwdPeeringPriority: h.GwdPeeringPriority, RateLimitPeeringPriority: h.RateLimitPeeringPriority,
			SubsPeeringPriority: h.SubsPeeringPriority, ApiProbePeeringPriority: h.ApiProbePeeringPriority,
		})
	}
}

func (h *HostGateway) gen(v10, verbose bool) {
	h.subsysgen = true

	// gateway = true
	h.HostVm.gen(v10, verbose, true)
}

type GwySubsysVm struct {
	SubsysName string
	Mode string // copy from defaults.mode
	Platform string // vmware, linux, docker
	v10 bool

	subsysgen bool // created by the subsystem generator with the gen() function

	SearchDomains []string // copy from defaults.search-domains
	DnsServers []string // copy from defaults.vm-first-boot.dns-servers

	// NTP servers (todo: list)
	NTPServers []string

	Hosts []HostGateway

	//PassiveDatapowerCluster []string

	ApiGateway string
	ApicGwService string

	// web-gui timeout
	WebGuiIdleTimeout int

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

	// low level peering configuration

	GwdPeering string
	GwdPeeringLocalPort int
	GwdPeeringMonitorPort int
	GwdPeeringPersistence string // memory | raid

	RateLimitPeering string
	RateLimitPeeringLocalPort int
	RateLimitPeeringMonitorPort int

	SubsPeering string
	SubsPeeringLocalPort int
	SubsPeeringMonitorPort int

	ApiProbePeering string
	ApiProbePeeringLocalPort int
	ApiProbePeeringMonitorPort int

	// todo: api-probe settings

	// datapower cert configuration
	DatapowerCryptoDir string

	// gateway director key: gwd_key
	DatapowerGwdKey string

	// gateway director cert: gwd_cert
	DatapowerGwdCert string

	// peering key and cert
	DatapowerPeeringKey string
	DatapowerPeeringCert string

	// web-gui key and cert
	DatapowerWebGuiKey string
	DatapowerWebGuiCert string

	// datapower api endpoint key and cert
	DatapowerApiEndpointKey string
	DatapowerApiEndpointCert string

	// ca file, root-ca file, ca-chain file (ca+root-ca)
	//CaFile string
	//RootCaFile string
	//CaChainFile string

	// datapower crypto ca certs
	DatapowerCaCert string
	DatapowerRootCert string
}

func (gwy *GwySubsysVm) MarshalJSON() ([]byte, error) {
	if gwy.v10 {
		if gwy.subsysgen {
			return json.Marshal(&struct{
				SubsysName string
				Platform string

				NTPServers []string
				Hosts []HostGateway

				ApiGateway string
				ApicGwService string
			}{
				SubsysName: gwy.SubsysName, Platform: gwy.Platform,
				NTPServers: gwy.NTPServers, Hosts: gwy.Hosts,
				ApiGateway: gwy.ApiGateway,ApicGwService: gwy.ApicGwService,
			})

		} else {
			return json.Marshal(&struct{

			}{

			})
		}

	} else if gwy.subsysgen {
		return json.Marshal(&struct{
			SubsysName string
			Platform string

			NTPServers []string
			Hosts []HostGateway

			ApiGateway string
			ApicGwService string
		}{
			SubsysName: gwy.SubsysName, Platform: gwy.Platform,
			NTPServers: gwy.NTPServers, Hosts: gwy.Hosts,
			ApiGateway: gwy.ApiGateway,ApicGwService: gwy.ApicGwService,
		})

	} else {
		return json.Marshal(&struct{

		}{

		})
	}
}

func (gwy *GwySubsysVm) gen(v10, verbose bool) {
	gwy.v10 = v10
	gwy.subsysgen = true

	gwy.SubsysName = "gwy"
	gwy.Mode = "standard"
	gwy.Platform = "vmware"

	gwy.NTPServers = []string{"time-a.ntp.org", "time-b.ntp.org"}

	gwy.Hosts = make([]HostGateway, 3)
	for j := range gwy.Hosts {
		gwy.Hosts[j].gen(v10, verbose)
	}

	gwy.ApicGwService = "gwd.my.domain.com"
	gwy.ApiGateway = "gwy.my.domain.com"
}

//func (gwy *GwySubsysVm) GetCaFileOrDefault() string {
//	caFileDefault := "dp-ca.pem"
//	if len(gwy.CaFile) == 0 { return caFileDefault }
//	return gwy.CaFile
//}
//
//func (gwy *GwySubsysVm) GetRootCaFileOrDefault() string {
//	rootCaFileDefault := "dp-root-ca.pem"
//	if len(gwy.RootCaFile) == 0 { return rootCaFileDefault }
//	return gwy.RootCaFile
//}
//
//func (gwy *GwySubsysVm) GetCaChainFileOrDefault() string {
//	caChainDefault := "dp-ca-chain.pem"
//	if len(gwy.CaChainFile) == 0 { return caChainDefault }
//	return gwy.CaChainFile
//}

func (gwy *GwySubsysVm) GetWebGuiIdleTimeoutOrDefault() int {
	webGuiIdleTimeoutDefault := 600
	if gwy.WebGuiIdleTimeout == 0 { return webGuiIdleTimeoutDefault }
	return gwy.WebGuiIdleTimeout
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

func (gwy *GwySubsysVm) GetNTPServersOrDefault() []string {
	if len(gwy.NTPServers) == 0 {
		return []string{"pool.ntp.org"}
	}
	return gwy.NTPServers
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

func (gwy *GwySubsysVm) GetApiProbePeeringOrDefault() string {
	if len(gwy.ApiProbePeering) == 0 {
		return apiProbePeering
	}
	return gwy.ApiProbePeering
}

func (gwy *GwySubsysVm) GetApiProbePeeringLocalPortOrDefault() int {
	if gwy.ApiProbePeeringLocalPort == 0 {
		return apiProbePeeringLocalPort
	}
	return gwy.ApiProbePeeringLocalPort
}

func (gwy *GwySubsysVm) GetApiProbePeeringMonitorPortOrDefault() int {
	if gwy.ApiProbePeeringLocalPort == 0 {
		return apiProbePeeringMonitorPort
	}
	return gwy.ApiProbePeeringMonitorPort
}

func (gwy *GwySubsysVm) GetApiGatewayPortOrDefault() int {
	if gwy.DatapowerApiGatewayPort == 0 { return defaultApiGatewayPort }
	return gwy.DatapowerApiGatewayPort
}

func (gwy *GwySubsysVm) GetApiGatewayAddressOrDefault() string {
	if len(gwy.DatapowerApiGatewayAddress) == 0 {return defaultApiGatewayAddress}
	return gwy.DatapowerApiGatewayAddress
}

func (gwy *GwySubsysVm) GetApicGwServicePortOrDefault() int {
	if gwy.DatapowerApicGwServicePort == 0 {return defaultApicGwServicePort}
	return gwy.DatapowerApicGwServicePort
}

func (gwy *GwySubsysVm) GetApicGwServiceAddressOrDefault() string {
	if len(gwy.DatapowerApicGwServiceAddress) == 0 {return defaultApicGwServiceAddress}
	return gwy.DatapowerApicGwServiceAddress
}

