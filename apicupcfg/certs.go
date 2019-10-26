package apicupcfg

import (
	rice "github.com/GeertJohan/go.rice"
	"strings"
)

type CertSpec struct {
	Cn           string
	CertName     string
	DnFields     []string
	K8sNamespace string
	SubsysName string
	KeyFile string
	CertFile string
	CaFile string
	CsrConf string
	CsrSubdir string // csr subdirectory relative to the base output directory
}

type OsEnvCerts struct {
	OsEnv     OsEnv
	CertSpecs map[string]CertSpec
}

type OsEnvCert struct {
	OsEnv OsEnv
	CertSpec CertSpec
}

type Certbot struct {
	CertDir string
	Cert string
	Key string
	CaChain string
}

type Certs struct {
	DnFields []string
	K8sNamespace string
	CaFile string

	PublicUserFacingCerts bool
	PublicCerts bool
	CommonCerts bool

	PublicUserFacingEkuServerAuth map[string]CertSpec
	PublicEkuServerAuth map[string]CertSpec

	MutualAuthEkuServerAuth map[string]CertSpec
	CommonEkuClientAuth map[string]CertSpec

	Certbot Certbot

	OsEnv
}

// subsystem certs
const CertKeyPlatformApi = "platform-api"
const CertKeyConsumerApi = "consumer-api"
const CertKeyApiManagerUi = "api-manager-ui"
const CertKeyCloudAdminUi = "cloud-admin-ui"
const CertKeyPortalAdminIngress = "portal-admin-ingress"
const CertKeyPortalWwwIngress = "portal-www-ingress"
const CertKeyAnalyticsClientIngress = "analytics-client-ingress"
const CertKeyAnalyticsIngestionIngress = "analytics-ingestion-ingress"
const CertKeyApicGwServiceIngress = "apic-gw-service-ingress"

// common certs
const CertKeyPortalClient = "portal-client"
const CertKeyAnalyticsClientClient = "analytics-client-client"
const CertKeyAnalyticsIngestionClient = "analytics-ingestion-client"

func updateCertSpec(certs *Certs, subsysName string, certName string, certSpec *CertSpec, csrSubdir string) {

	// check if no cn...
	if len(certSpec.Cn) == 0 {
		return
	}

	if len(certSpec.CertName) == 0 {
		certSpec.CertName = certName
	}

	if len(certSpec.DnFields) == 0 {
		certSpec.DnFields = certs.DnFields
	}

	if len(certSpec.K8sNamespace) == 0 {
		certSpec.K8sNamespace = certs.K8sNamespace
	}

	if len(certSpec.CaFile) == 0 {
		certSpec.CaFile = certs.CaFile
	}

	certSpec.SubsysName = subsysName

	// default file names
	if len(certSpec.CsrConf) == 0 {
		certSpec.CsrConf = strings.ReplaceAll(certSpec.Cn, ".", "-") + ".conf"
	}

	if len(certSpec.KeyFile) == 0 {
		certSpec.KeyFile = strings.ReplaceAll(certSpec.Cn, ".", "-") + ".key"
	}

	if len(certSpec.CertFile) == 0 {
		certSpec.CertFile = strings.ReplaceAll(certSpec.Cn, ".", "-") + ".crt"
	}

	if len(certSpec.CsrSubdir) == 0 {
		certSpec.CsrSubdir = csrSubdir
	}
}

func createCertMaps(certs *Certs) {
	certs.PublicUserFacingEkuServerAuth = make(map[string]CertSpec)
	certs.PublicEkuServerAuth = make(map[string]CertSpec)
	certs.MutualAuthEkuServerAuth = make(map[string]CertSpec)
	certs.CommonEkuClientAuth = make(map[string]CertSpec)
}

func updateCertSpecs(certs *Certs, mgmt ManagementSubsysDescriptor, alyt AnalyticsSubsysDescriptor,
	ptl PortalSubsysDescriptor, gwy GatewaySubsysDescriptor, commonCsrOutDir string, customCsrOutDir string) {

	certs.OsEnv.init()
	createCertMaps(certs)

	isManagement := len(mgmt.GetManagementSubsysName()) > 0
	isAnalytics := len(alyt.GetAnalyticsSubsysName()) > 0
	isPortal := len(ptl.GetPortalSubsysName()) > 0
	isGateway := len(gwy.GetGatewaySubsysName()) > 0

	if isManagement && certs.PublicUserFacingCerts {
		// management subsystem contributes public-user-facing certs

		// build cert specs
		certmap := certs.PublicUserFacingEkuServerAuth

		certSpec := CertSpec{}
		certSpec.Cn = mgmt.GetPlatformApiEndpoint()
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyPlatformApi, &certSpec, customCsrOutDir)
		certmap[CertKeyPlatformApi] = certSpec

		certSpec = CertSpec{}
		certSpec.Cn = mgmt.GetConsumerApiEndpoint()
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyConsumerApi, &certSpec, customCsrOutDir)
		certmap[CertKeyConsumerApi] = certSpec

		certSpec = CertSpec{}
		certSpec.Cn = mgmt.GetApiManagerUIEndpoint()
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyApiManagerUi, &certSpec, customCsrOutDir)
		certmap[CertKeyApiManagerUi] = certSpec

		certSpec = CertSpec{}
		certSpec.Cn = mgmt.GetCloudAdminUIEndpoint()
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyCloudAdminUi, &certSpec, customCsrOutDir)
		certmap[CertKeyCloudAdminUi] = certSpec
	}

	if isPortal && certs.PublicUserFacingCerts {
		// portal subsystem contributes public-user-facing certs
		certmap := certs.PublicUserFacingEkuServerAuth

		certSpec := CertSpec{}
		certSpec.Cn = ptl.GetPortalWWWEndpoint()
		updateCertSpec(certs, ptl.GetPortalSubsysName(), CertKeyPortalWwwIngress, &certSpec, customCsrOutDir)
		certmap[CertKeyPortalWwwIngress] = certSpec
	}

	if isGateway && certs.PublicCerts {
		// gateway contributes to public certs
		certmap := certs.PublicEkuServerAuth

		certSpec := CertSpec{}
		certSpec.Cn = CertKeyApicGwServiceIngress
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyApicGwServiceIngress, &certSpec, commonCsrOutDir)
		certmap[CertKeyApicGwServiceIngress] = certSpec
	}

	if isManagement && certs.CommonCerts {
		// common certs are set on the management subystem
		certmap := certs.CommonEkuClientAuth

		certSpec := CertSpec{}
		certSpec.Cn = CertKeyPortalClient
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyPortalClient, &certSpec, commonCsrOutDir)
		certmap[CertKeyPortalClient] = certSpec

		certSpec = CertSpec{}
		certSpec.Cn = CertKeyAnalyticsClientClient
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyAnalyticsClientClient, &certSpec, commonCsrOutDir)
		certmap[CertKeyAnalyticsClientClient] = certSpec

		certSpec = CertSpec{}
		certSpec.Cn = CertKeyAnalyticsIngestionClient
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyAnalyticsIngestionClient, &certSpec, commonCsrOutDir)
		certmap[CertKeyAnalyticsIngestionClient] = certSpec
	}

	if isPortal && certs.CommonCerts {
		// portal subsystem contributes mutual auth server cert
		certmap := certs.MutualAuthEkuServerAuth

		certSpec := CertSpec{}
		certSpec.Cn = ptl.GetPortalAdminEndpoint()
		updateCertSpec(certs, ptl.GetPortalSubsysName(), CertKeyPortalAdminIngress, &certSpec, commonCsrOutDir)
		certmap[CertKeyPortalAdminIngress] = certSpec
	}

	if isAnalytics && certs.CommonCerts {
		// analytics subsystem contributes mutual auth server certs
		certmap := certs.MutualAuthEkuServerAuth

		certSpec := CertSpec{}
		certSpec.Cn = alyt.GetAnalyticsIngestionEndpoint()
		updateCertSpec(certs, alyt.GetAnalyticsSubsysName(), CertKeyAnalyticsIngestionIngress, &certSpec, commonCsrOutDir)
		certmap[CertKeyAnalyticsIngestionIngress] = certSpec

		certSpec = CertSpec{}
		certSpec.Cn = alyt.GetAnalyticsClientEndpoint()
		updateCertSpec(certs, alyt.GetAnalyticsSubsysName(), CertKeyAnalyticsClientIngress, &certSpec, commonCsrOutDir)
		certmap[CertKeyAnalyticsClientIngress] = certSpec
	}
}

func outputCerts(certs *Certs, outfiles map[string]string, tag string, tbox *rice.Box) {

	ekuServerAuth := parseTemplate(tbox, tpdir(tbox) + "csr-server-auth.tmpl")
	ekuClientAuth := parseTemplate(tbox, tpdir(tbox) + "csr-client-auth.tmpl")
	keypairTemplate := parseTemplates(tbox, tpdir(tbox) + "keypair.tmpl", tpdir(tbox) + "helpers.tmpl")
	combinedCsrTemplate := parseTemplates(tbox, tpdir(tbox) + "combined-csr.tmpl", tpdir(tbox) + "helpers.tmpl")

	subsysCertsTemplate := parseTemplates(tbox, tpdir(tbox) + "subsys-certs.tmpl", tpdir(tbox) + "helpers.tmpl")

	var outpath string

	var osenv OsEnv
	osenv.init()

	// public-user-facing-eku-server-auth csr conf
	if certs.PublicUserFacingCerts {
		for _, certSpec := range certs.PublicUserFacingEkuServerAuth {

			if len(certSpec.Cn) > 0 {
				// csr
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[customCsrOutDir]), certSpec.CsrConf)
				writeTemplate(ekuServerAuth, outpath, certSpec)

				// key-pair
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[customCsrOutDir]), certSpec.CsrConf + osenv.ShellExt)
				writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec})
			}
		}
	}

	if certs.PublicCerts {
		for _, certSpec := range certs.PublicEkuServerAuth {

			if len(certSpec.Cn) > 0 {
				// csr
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[customCsrOutDir]), certSpec.CsrConf)
				writeTemplate(ekuServerAuth, outpath, certSpec)

				// key-pair
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[customCsrOutDir]), certSpec.CsrConf + osenv.ShellExt)
				writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec})
			}
		}
	}

	if certs.CommonCerts {
		for _, certSpec := range certs.MutualAuthEkuServerAuth {

			if len(certSpec.Cn) > 0 {
				// csr
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[commonCsrOutDir]), certSpec.CsrConf)
				writeTemplate(ekuServerAuth, outpath, certSpec)

				// key-pair
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[commonCsrOutDir]), certSpec.CsrConf + osenv.ShellExt)
				writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec})
			}
		}

		for _, certSpec := range certs.CommonEkuClientAuth {

			if len(certSpec.Cn) > 0 {
				// csr
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[commonCsrOutDir]), certSpec.CsrConf)
				writeTemplate(ekuClientAuth, outpath, certSpec)

				// key-pair
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[commonCsrOutDir]), certSpec.CsrConf + osenv.ShellExt)
				writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec})
			}
		}
	}

	if certs.PublicUserFacingCerts {
		// combine public-user-facing key and csr scripts
		outpath = fileName(concatSubdir(outfiles[outdir], outfiles[customCsrOutDir]), tagOutputFileName("all-user-facing-public-csr", tag) + osenv.ShellExt)
		writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.PublicUserFacingEkuServerAuth})
	}

	if certs.PublicCerts {
		// combine public key and csr scripts
		outpath = fileName(concatSubdir(outfiles[outdir], outfiles[customCsrOutDir]), tagOutputFileName("all-public-csr", tag) + osenv.ShellExt)
		writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.PublicEkuServerAuth})
	}

	if certs.CommonCerts {
		// combine mutual-auth key and csr scripts
		outpath = fileName(concatSubdir(outfiles[outdir], outfiles[commonCsrOutDir]), tagOutputFileName("all-mutual-auth-csr", tag) + osenv.ShellExt)
		writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.MutualAuthEkuServerAuth})

		// combine common key and csr scripts
		outpath = fileName(concatSubdir(outfiles[outdir], outfiles[commonCsrOutDir]), tagOutputFileName("all-common-csr", tag) + osenv.ShellExt)
		writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.CommonEkuClientAuth})
	}

	if certs.PublicUserFacingCerts {
		// apicup certs user-facing-public
		outpath = fileName(outfiles["outdir"], tagOutputFileName(outfiles[userFacingPublicCertsOut], tag)) + osenv.ShellExt
		writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.PublicUserFacingEkuServerAuth})
	}

	if certs.PublicCerts {
		// apicup certs public
		outpath = fileName(outfiles["outdir"], tagOutputFileName(outfiles[publicCertsOut], tag)) + osenv.ShellExt
		writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.PublicEkuServerAuth})
	}

	if certs.CommonCerts {
		// apicup certs mutual auth
		outpath = fileName(outfiles["outdir"], tagOutputFileName(outfiles[mutualAuthCertsOut], tag)) + osenv.ShellExt
		writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.MutualAuthEkuServerAuth})

		// apicup certs common
		outpath = fileName(outfiles["outdir"], tagOutputFileName(outfiles[commonCertsOut], tag)) + osenv.ShellExt
		writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.CommonEkuClientAuth})
	}

	if len(certs.Certbot.CertDir) > 0 {

		updateFromCertbot := func(certs map[string]CertSpec, certbot Certbot) map[string]CertSpec {
			cbmap := make(map[string]CertSpec)
			for cname, cs := range certs {
				cs.CsrSubdir = certbot.CertDir
				cs.KeyFile = certbot.Key
				cs.CertFile = certbot.Cert
				cs.CaFile = certbot.CaChain
				cbmap[cname] = cs
			}
			return cbmap
		}

		// apicup certs user-facing-public certbot
		if certs.PublicUserFacingCerts {
			cbmap := updateFromCertbot(certs.PublicUserFacingEkuServerAuth, certs.Certbot)
			outpath = fileName(outfiles["outdir"], tagOutputFileName(outfiles[certbotUserFacingPublicCertOut], tag)) + osenv.ShellExt
			writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: cbmap})
		}

		// apicup certs public certbot
		if certs.PublicCerts {
			cbmap := updateFromCertbot(certs.PublicEkuServerAuth, certs.Certbot)
			outpath = fileName(outfiles["outdir"], tagOutputFileName(outfiles[certbotPublicCertOut], tag)) + osenv.ShellExt
			writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: cbmap})
		}
	}
}
