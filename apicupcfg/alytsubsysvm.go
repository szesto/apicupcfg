package apicupcfg

import "encoding/json"

type AltSubsysVm struct {
	SubsysVmBase

	AnalyticsIngestion string
	AnalyticsClient string

	EnableMessageQueue bool

	subsysGen bool
}

func (alt *AltSubsysVm) MarshalJSON() ([]byte, error) {

	if alt.v10 {
		if alt.subsysGen {
			return json.Marshal(&struct {
				SubsysName string
				VmFirstBoot VmFirstBoot

				AnalyticsIngestion string
				AnalyticsClient string
			}{
				SubsysName: alt.SubsysName, VmFirstBoot: alt.VmFirstBoot,
				AnalyticsIngestion: alt.AnalyticsIngestion, AnalyticsClient: alt.AnalyticsClient,
			})

		} else {
			return json.Marshal(&struct {
				SubsysName string
				Mode string
				CloudInit CloudInit
				SearchDomains []string
				VmFirstBoot VmFirstBoot
				SshPublicKeyFile string

				AnalyticsIngestion string
				AnalyticsClient string
			}{
				SubsysName: alt.SubsysName, Mode: alt.Mode, CloudInit: alt.CloudInit,
				SearchDomains: alt.SearchDomains, VmFirstBoot: alt.VmFirstBoot,
				SshPublicKeyFile: alt.SshPublicKeyFile,
				AnalyticsIngestion: alt.AnalyticsIngestion,
				AnalyticsClient: alt.AnalyticsClient,
			})
		}

	} else if alt.subsysGen {
		return json.Marshal(&struct {
			SubsysName string
			VmFirstBoot VmFirstBoot

			EnableMessageQueue bool

			AnalyticsIngestion string
			AnalyticsClient string
		}{
			SubsysName: alt.SubsysName, VmFirstBoot: alt.VmFirstBoot, EnableMessageQueue: alt.EnableMessageQueue,
			AnalyticsIngestion: alt.AnalyticsIngestion, AnalyticsClient: alt.AnalyticsClient,
		})

	} else {
		return json.Marshal(&struct {
			SubsysName string
			Mode string
			CloudInit CloudInit
			SearchDomains []string
			VmFirstBoot VmFirstBoot
			SshPublicKeyFile string

			EnableMessageQueue bool

			AnalyticsIngestion string
			AnalyticsClient string
		}{
			SubsysName: alt.SubsysName, Mode: alt.Mode, CloudInit: alt.CloudInit,
			SearchDomains: alt.SearchDomains, VmFirstBoot: alt.VmFirstBoot,
			SshPublicKeyFile: alt.SshPublicKeyFile, EnableMessageQueue: alt.EnableMessageQueue,
			AnalyticsIngestion: alt.AnalyticsIngestion,
			AnalyticsClient: alt.AnalyticsClient,
		})
	}
}

func (alt *AltSubsysVm) gen(v10, verbose bool) {
	// subsys = true; gateway = false
	alt.SubsysVmBase.gen(v10, verbose, true, false)

	alt.AnalyticsClient = "ac.my.domain.com"
	alt.AnalyticsIngestion = "ai.my.domain.com"

	alt.EnableMessageQueue = false

	alt.subsysGen = true
}
