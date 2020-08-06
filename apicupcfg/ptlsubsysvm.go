package apicupcfg

import "encoding/json"

type PtlSubsysVm struct {
	SubsysVmBase

	SiteBackup Backup

	PortalAdmin string
	PortalWww string

	// v10 fields
	DeploymentProfile string // move to subsys-vm-base, do not overload Mode

	MultiSiteHAEnabled bool // tech preview
	MultiSiteHAMode string // tech preview
	ReplicationPeer string // tech preview
	SiteName string // tech preview
	PortalReplication string // tech preview

	subsysGen bool
}

func (ptl *PtlSubsysVm) MarshalJSON() ([]byte, error) {
	if ptl.v10 {
		if ptl.subsysGen {
			return json.Marshal(&struct{
				SubsysName string
				VmFirstBoot VmFirstBoot

				SiteBackup Backup
				PortalAdmin string
				PortalWww string
			}{
				SubsysName: ptl.SubsysName, VmFirstBoot: ptl.VmFirstBoot,
				SiteBackup: ptl.SiteBackup, PortalAdmin: ptl.PortalAdmin, PortalWww: ptl.PortalWww,
			})

		} else {
			return json.Marshal(&struct{
				SubsysName string
				Mode string
				CloudInit CloudInit
				SearchDomains []string
				VmFirstBoot VmFirstBoot
				SshPublicKeyFile string

				SiteBackup Backup
				PortalAdmin string
				PortalWww string
			}{
				SubsysName: ptl.SubsysName, Mode: ptl.Mode, CloudInit: ptl.CloudInit,
				SearchDomains: ptl.SearchDomains, VmFirstBoot: ptl.VmFirstBoot, SshPublicKeyFile: ptl.SshPublicKeyFile,
				SiteBackup: ptl.SiteBackup, PortalAdmin: ptl.PortalAdmin, PortalWww: ptl.PortalWww,
			})
		}

	} else if ptl.subsysGen {
		return json.Marshal(&struct{
			SubsysName string
			VmFirstBoot VmFirstBoot

			SiteBackup Backup
			PortalAdmin string
			PortalWww string
		}{
			SubsysName: ptl.SubsysName, VmFirstBoot: ptl.VmFirstBoot,
			SiteBackup: ptl.SiteBackup, PortalAdmin: ptl.PortalAdmin, PortalWww: ptl.PortalWww,
		})

	} else {
		return json.Marshal(&struct{
			SubsysName string
			Mode string
			CloudInit CloudInit
			SearchDomains []string
			VmFirstBoot VmFirstBoot
			SshPublicKeyFile string

			SiteBackup Backup
			PortalAdmin string
			PortalWww string
		}{
			SubsysName: ptl.SubsysName, Mode: ptl.Mode, CloudInit: ptl.CloudInit,
			SearchDomains: ptl.SearchDomains, VmFirstBoot: ptl.VmFirstBoot, SshPublicKeyFile: ptl.SshPublicKeyFile,
			SiteBackup: ptl.SiteBackup, PortalAdmin: ptl.PortalAdmin, PortalWww: ptl.PortalWww,
		})
	}
}

func (ptl *PtlSubsysVm) gen(v10, verbose bool) {
	// subsys = true; gateway = false
	ptl.SubsysVmBase.gen(v10, verbose, true, false)

	// somehow backup config refers to the k8s install
	ptl.SiteBackup.gen(false, verbose)

	ptl.PortalAdmin = "padmin.my.domain.com"
	ptl.PortalWww = "portal.my.domain.com"

	ptl.subsysGen = true
}
