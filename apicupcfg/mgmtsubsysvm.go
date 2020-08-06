package apicupcfg

import "encoding/json"

type MgtSubsysVm struct {
	SubsysVmBase

	CassandraEncryptionKeyFile string

	Backup Backup

	PlatformApi string
	ApiManagerUi string
	CloudAdminUi string
	ConsumerApi string

	subsysGen bool
}

func (mgt *MgtSubsysVm) MarshalJSON() ([]byte, error) {

	if mgt.v10 {
		if mgt.subsysGen {
			return json.Marshal(&struct{
				SubsysName string
				VmFirstBoot VmFirstBoot

				Backup Backup
				PlatformApi string
				ApiManagerUi string
				CloudAdminUi string
				ConsumerApi string
			}{
				SubsysName: mgt.SubsysName,
				VmFirstBoot: mgt.VmFirstBoot,
				Backup: mgt.Backup, PlatformApi: mgt.PlatformApi, ApiManagerUi: mgt.ApiManagerUi,
				CloudAdminUi: mgt.CloudAdminUi, ConsumerApi: mgt.ConsumerApi,
			})

		} else {
			return json.Marshal(&struct{
				SubsysName string
				Mode string
				CloudInit CloudInit
				SearchDomains []string
				VmFirstBoot VmFirstBoot
				SshPublicKeyFile string

				Backup Backup
				PlatformApi string
				ApiManagerUi string
				CloudAdminUi string
				ConsumerApi string
			}{
				SubsysName: mgt.SubsysName, Mode: mgt.Mode, CloudInit: mgt.CloudInit, SearchDomains: mgt.SearchDomains,
				VmFirstBoot: mgt.VmFirstBoot, SshPublicKeyFile: mgt.SshPublicKeyFile,
				Backup: mgt.Backup, PlatformApi: mgt.PlatformApi, ApiManagerUi: mgt.ApiManagerUi,
				CloudAdminUi: mgt.CloudAdminUi, ConsumerApi: mgt.ConsumerApi,
			})
		}

	} else if mgt.subsysGen {
		return json.Marshal(&struct{
			SubsysName string
			VmFirstBoot VmFirstBoot

			Backup Backup

			PlatformApi string
			ApiManagerUi string
			CloudAdminUi string
			ConsumerApi string
		}{
			SubsysName: mgt.SubsysName,
			VmFirstBoot: mgt.VmFirstBoot,
			Backup: mgt.Backup, PlatformApi: mgt.PlatformApi, ApiManagerUi: mgt.ApiManagerUi,
			CloudAdminUi: mgt.CloudAdminUi, ConsumerApi: mgt.ConsumerApi,
		})

	} else {
		return json.Marshal(&struct{
			SubsysName string
			Mode string
			CloudInit CloudInit
			SearchDomains []string
			VmFirstBoot VmFirstBoot
			SshPublicKeyFile string

			CassandraEncryptionKeyFile string

			Backup Backup

			PlatformApi string
			ApiManagerUi string
			CloudAdminUi string
			ConsumerApi string
		}{
			SubsysName: mgt.SubsysName, Mode: mgt.Mode, CloudInit: mgt.CloudInit, SearchDomains: mgt.SearchDomains,
			VmFirstBoot: mgt.VmFirstBoot, SshPublicKeyFile: mgt.SshPublicKeyFile,
			CassandraEncryptionKeyFile: mgt.CassandraEncryptionKeyFile,
			Backup: mgt.Backup, PlatformApi: mgt.PlatformApi, ApiManagerUi: mgt.ApiManagerUi,
			CloudAdminUi: mgt.CloudAdminUi, ConsumerApi: mgt.ConsumerApi,
		})
	}
}

func (mgt *MgtSubsysVm) gen(v10, verbose bool) {
	// subsys = true; gateway = false;
	mgt.SubsysVmBase.gen(v10, verbose, true, false)

	mgt.Backup.gen(v10, verbose)

	mgt.PlatformApi = "api.my.domain.com"
	mgt.ApiManagerUi = "apim.my.domain.com"
	mgt.CloudAdminUi = "cm.my.domain.com"
	mgt.ConsumerApi = "consumer.my.domain.com"

	mgt.subsysGen = true
}

