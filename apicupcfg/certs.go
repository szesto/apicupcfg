package apicupcfg

import (
	"fmt"
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
	CertSubdir string // cert subdirectory relative to the base output directory
	KeySubdir string // key subdirectory relative tp the base output directory
	CaSubdir string // ca subdirectory relative to the base output directory
	AltCns [] string // a list of alt cn's
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

	SharedEndpointTrust bool

	etCertSpecMgt CertSpec
	etCertSpecPtl CertSpec
	etCertSpecAlyt CertSpec
	etCertSpecGwy CertSpec

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

// shared trust endpoint key files
const KeyFileSharedEndpointTrustManagement = "shared-endpoint-trust-management.key"
const KeyFileSharedEndpointTrustPortal = "shared-endpoint-trust-portal.key"
const KeyFileSharedEndpointTrustAnalytics = "shared-endpoint-trust-analytics.key"
const KeyFileSharedEndpointTrustGateway = "shared-endpoint-trust-gateway.key"

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

	if len(certSpec.CertSubdir) == 0 {
		certSpec.CertSubdir = csrSubdir
	}

	if len(certSpec.KeySubdir) == 0 {
		certSpec.KeySubdir = csrSubdir
	}

	if len(certSpec.CaSubdir) == 0 {
		certSpec.CaSubdir = csrSubdir
	}

	// todo: reallocate existing alt-cns
	certSpec.AltCns = make([]string, 0, 50)
}

func createCertMaps(certs *Certs) {
	if certs.PublicUserFacingEkuServerAuth == nil {
		certs.PublicUserFacingEkuServerAuth = make(map[string]CertSpec)
	}

	if certs.PublicEkuServerAuth == nil {
		certs.PublicEkuServerAuth = make(map[string]CertSpec)
	}

	if certs.MutualAuthEkuServerAuth == nil {
		certs.MutualAuthEkuServerAuth = make(map[string]CertSpec)
	}

	if certs.CommonEkuClientAuth == nil {
		certs.CommonEkuClientAuth = make(map[string]CertSpec)
	}
}

func updateCertSpecs(certs *Certs, mgmt ManagementSubsysDescriptor, alyt AnalyticsSubsysDescriptor,
	ptl PortalSubsysDescriptor, gwy GatewaySubsysDescriptor, commonCsrOutDir string, customCsrOutDir string) {

	certs.OsEnv.init()
	createCertMaps(certs)

	isManagement := len(mgmt.GetManagementSubsysName()) > 0
	isAnalytics := len(alyt.GetAnalyticsSubsysName()) > 0
	isPortal := len(ptl.GetPortalSubsysName()) > 0
	isGateway := len(gwy.GetGatewaySubsysName()) > 0

	getCertSpec := func (certmap map[string]CertSpec, key string) CertSpec {
		if certspec, ok := certmap[key]; ok {return certspec } else {return CertSpec{}}
	}

	if isManagement && certs.PublicUserFacingCerts {
		// management subsystem contributes public-user-facing certs

		// build cert specs
		certmap := certs.PublicUserFacingEkuServerAuth

		certSpec := getCertSpec(certmap, CertKeyPlatformApi)
		certSpec.Cn = mgmt.GetPlatformApiEndpoint()
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyPlatformApi, &certSpec, customCsrOutDir)
		certmap[CertKeyPlatformApi] = certSpec

		certSpec = getCertSpec(certmap, CertKeyConsumerApi)
		certSpec.Cn = mgmt.GetConsumerApiEndpoint()
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyConsumerApi, &certSpec, customCsrOutDir)
		certmap[CertKeyConsumerApi] = certSpec

		certSpec = getCertSpec(certmap, CertKeyApiManagerUi)
		certSpec.Cn = mgmt.GetApiManagerUIEndpoint()
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyApiManagerUi, &certSpec, customCsrOutDir)
		certmap[CertKeyApiManagerUi] = certSpec

		certSpec = getCertSpec(certmap,CertKeyCloudAdminUi)
		certSpec.Cn = mgmt.GetCloudAdminUIEndpoint()
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyCloudAdminUi, &certSpec, customCsrOutDir)
		certmap[CertKeyCloudAdminUi] = certSpec
	}

	if isPortal && certs.PublicUserFacingCerts {
		// portal subsystem contributes public-user-facing certs
		certmap := certs.PublicUserFacingEkuServerAuth

		certSpec := getCertSpec(certmap, CertKeyPortalWwwIngress)
		certSpec.Cn = ptl.GetPortalWWWEndpoint()
		updateCertSpec(certs, ptl.GetPortalSubsysName(), CertKeyPortalWwwIngress, &certSpec, customCsrOutDir)
		certmap[CertKeyPortalWwwIngress] = certSpec
	}

	if isGateway && certs.PublicCerts {
		// gateway contributes to public certs
		certmap := certs.PublicEkuServerAuth

		certSpec := getCertSpec(certmap, CertKeyApicGwServiceIngress)
		certSpec.Cn = CertKeyApicGwServiceIngress
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyApicGwServiceIngress, &certSpec, commonCsrOutDir)
		certmap[CertKeyApicGwServiceIngress] = certSpec
	}

	if isManagement && certs.CommonCerts {
		// common certs are set on the management subystem
		certmap := certs.CommonEkuClientAuth

		certSpec := getCertSpec(certmap, CertKeyPortalClient)
		certSpec.Cn = CertKeyPortalClient
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyPortalClient, &certSpec, commonCsrOutDir)
		certmap[CertKeyPortalClient] = certSpec

		certSpec = getCertSpec(certmap, CertKeyAnalyticsClientClient)
		certSpec.Cn = CertKeyAnalyticsClientClient
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyAnalyticsClientClient, &certSpec, commonCsrOutDir)
		certmap[CertKeyAnalyticsClientClient] = certSpec

		certSpec = getCertSpec(certmap, CertKeyAnalyticsIngestionClient)
		certSpec.Cn = CertKeyAnalyticsIngestionClient
		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyAnalyticsIngestionClient, &certSpec, commonCsrOutDir)
		certmap[CertKeyAnalyticsIngestionClient] = certSpec
	}

	if isPortal && certs.CommonCerts {
		// portal subsystem contributes mutual auth server cert
		certmap := certs.MutualAuthEkuServerAuth

		certSpec := getCertSpec(certmap, CertKeyPortalAdminIngress)
		certSpec.Cn = ptl.GetPortalAdminEndpoint()
		updateCertSpec(certs, ptl.GetPortalSubsysName(), CertKeyPortalAdminIngress, &certSpec, commonCsrOutDir)
		certmap[CertKeyPortalAdminIngress] = certSpec
	}

	if isAnalytics && certs.CommonCerts {
		// analytics subsystem contributes mutual auth server certs
		certmap := certs.MutualAuthEkuServerAuth

		certSpec := getCertSpec(certmap, CertKeyAnalyticsIngestionIngress)
		certSpec.Cn = alyt.GetAnalyticsIngestionEndpoint()
		updateCertSpec(certs, alyt.GetAnalyticsSubsysName(), CertKeyAnalyticsIngestionIngress, &certSpec, commonCsrOutDir)

		certSpec.AltCns = append(certSpec.AltCns, fmt.Sprintf("*.%s", certSpec.K8sNamespace))
		certSpec.AltCns = append(certSpec.AltCns, fmt.Sprintf("*.%s.svc", certSpec.K8sNamespace))
		certSpec.AltCns = append(certSpec.AltCns, fmt.Sprintf("*.%s.svc.cluster.local", certSpec.K8sNamespace))

		certmap[CertKeyAnalyticsIngestionIngress] = certSpec

		certSpec = getCertSpec(certmap, CertKeyAnalyticsClientIngress)
		certSpec.Cn = alyt.GetAnalyticsClientEndpoint()
		updateCertSpec(certs, alyt.GetAnalyticsSubsysName(), CertKeyAnalyticsClientIngress, &certSpec, commonCsrOutDir)

		certSpec.AltCns = append(certSpec.AltCns, fmt.Sprintf("*.%s", certSpec.K8sNamespace))
		certSpec.AltCns = append(certSpec.AltCns, fmt.Sprintf("*.%s.svc", certSpec.K8sNamespace))
		certSpec.AltCns = append(certSpec.AltCns, fmt.Sprintf("*.%s.svc.cluster.local", certSpec.K8sNamespace))

		certmap[CertKeyAnalyticsClientIngress] = certSpec
	}

	// shared endpoint trust
	if certs.SharedEndpointTrust {
		setupSharedEndpointTrust(certs, mgmt, alyt, ptl, gwy, commonCsrOutDir, customCsrOutDir)
	}
}

func setupSharedEndpointTrust(certs *Certs, mgmt ManagementSubsysDescriptor, alyt AnalyticsSubsysDescriptor,
	ptl PortalSubsysDescriptor, gwy GatewaySubsysDescriptor, commonCsrOutDir string, customCsrOutDir string) {

	// mgmt
	cs := &certs.etCertSpecMgt
	cs.Cn = mgmt.GetPlatformApiEndpoint()
	cs.CsrConf = "shared-endpoint-trust-management.conf"
	cs.CsrSubdir = SharedCsrOutDir
	cs.KeySubdir = SharedCsrOutDir
	cs.KeyFile = KeyFileSharedEndpointTrustManagement

	updateCertSpec(certs, mgmt.GetManagementSubsysName(), mgmt.GetManagementSubsysName(), cs, customCsrOutDir)
	//cs.AltCns = append(cs.AltCns, mgmt.GetPlatformApiEndpoint())
	cs.AltCns = append(cs.AltCns, mgmt.GetConsumerApiEndpoint())
	cs.AltCns = append(cs.AltCns, mgmt.GetCloudAdminUIEndpoint())
	cs.AltCns = append(cs.AltCns, mgmt.GetApiManagerUIEndpoint())

	// ptl
	cs = &certs.etCertSpecPtl
	cs.Cn = ptl.GetPortalWWWEndpoint()
	cs.CsrConf = "shared-endpoint-trust-portal.conf"
	cs.CsrSubdir = SharedCsrOutDir
	cs.KeySubdir = SharedCsrOutDir
	cs.KeyFile = KeyFileSharedEndpointTrustPortal

	updateCertSpec(certs, ptl.GetPortalSubsysName(), ptl.GetPortalSubsysName(), cs, customCsrOutDir)
	//cs.AltCns = append(cs.AltCns, ptl.GetPortalWWWEndpoint())
	cs.AltCns = append(cs.AltCns, ptl.GetPortalAdminEndpoint())

	// alyt
	cs = &certs.etCertSpecAlyt
	cs.Cn = alyt.GetAnalyticsClientEndpoint()
	cs.CsrConf = "shared-endpoint-trust-analytics.conf"
	cs.CsrSubdir = SharedCsrOutDir
	cs.KeySubdir = SharedCsrOutDir
	cs.KeyFile = KeyFileSharedEndpointTrustAnalytics

	updateCertSpec(certs, alyt.GetAnalyticsSubsysName(), alyt.GetAnalyticsSubsysName(), cs, customCsrOutDir)
	//cs.AltCns = append(cs.AltCns, alyt.GetAnalyticsClientEndpoint())
	cs.AltCns = append(cs.AltCns, alyt.GetAnalyticsIngestionEndpoint())

	// gwy
	cs = &certs.etCertSpecGwy
	cs.Cn = gwy.GetApicGatewayServiceEndpoint()
	cs.CsrConf = "shared-endpoint-trust-gateway.conf"
	cs.CsrSubdir = SharedCsrOutDir
	cs.KeySubdir = SharedCsrOutDir
	cs.KeyFile = KeyFileSharedEndpointTrustGateway

	updateCertSpec(certs, gwy.GetGatewaySubsysName(), gwy.GetGatewaySubsysName(), cs, customCsrOutDir)
	//cs.AltCns = append(cs.AltCns, gwy.GetApicGatewayServiceEndpoint())
	cs.AltCns = append(cs.AltCns, gwy.GetApiGatewayEndpoint())
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
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[CustomCsrOutDir]), certSpec.CsrConf)
				writeTemplate(ekuServerAuth, outpath, certSpec)

				// key-pair
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[CustomCsrOutDir]), certSpec.CsrConf + osenv.ShellExt)
				writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec})
			}
		}
	}

	if certs.PublicCerts {
		for _, certSpec := range certs.PublicEkuServerAuth {

			if len(certSpec.Cn) > 0 {
				// csr
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[CustomCsrOutDir]), certSpec.CsrConf)
				writeTemplate(ekuServerAuth, outpath, certSpec)

				// key-pair
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[CustomCsrOutDir]), certSpec.CsrConf + osenv.ShellExt)
				writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec})
			}
		}
	}

	if certs.CommonCerts {
		for _, certSpec := range certs.MutualAuthEkuServerAuth {

			if len(certSpec.Cn) > 0 {
				// csr
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[CommonCsrOutDir]), certSpec.CsrConf)
				writeTemplate(ekuServerAuth, outpath, certSpec)

				// key-pair
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[CommonCsrOutDir]), certSpec.CsrConf + osenv.ShellExt)
				writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec})
			}
		}

		for _, certSpec := range certs.CommonEkuClientAuth {

			if len(certSpec.Cn) > 0 {
				// csr
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[CommonCsrOutDir]), certSpec.CsrConf)
				writeTemplate(ekuClientAuth, outpath, certSpec)

				// key-pair
				outpath = fileName(concatSubdir(outfiles[outdir], outfiles[CommonCsrOutDir]), certSpec.CsrConf + osenv.ShellExt)
				writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec})
			}
		}
	}

	if certs.PublicUserFacingCerts {
		// combine public-user-facing key and csr scripts
		outpath = fileName(concatSubdir(outfiles[outdir], outfiles[CustomCsrOutDir]), tagOutputFileName("all-user-facing-public-csr", tag) + osenv.ShellExt)
		writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.PublicUserFacingEkuServerAuth})
	}

	if certs.PublicCerts {
		// combine public key and csr scripts
		outpath = fileName(concatSubdir(outfiles[outdir], outfiles[CustomCsrOutDir]), tagOutputFileName("all-public-csr", tag) + osenv.ShellExt)
		writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.PublicEkuServerAuth})
	}

	if certs.CommonCerts {
		// combine mutual-auth key and csr scripts
		outpath = fileName(concatSubdir(outfiles[outdir], outfiles[CommonCsrOutDir]), tagOutputFileName("all-mutual-auth-csr", tag) + osenv.ShellExt)
		writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.MutualAuthEkuServerAuth})

		// combine common key and csr scripts
		outpath = fileName(concatSubdir(outfiles[outdir], outfiles[CommonCsrOutDir]), tagOutputFileName("all-common-csr", tag) + osenv.ShellExt)
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

	if certs.SharedEndpointTrust {

		updateFromEndpointTrust := func(certs map[string]CertSpec, et map[string]CertSpec) map[string]CertSpec {
			cbmap := make(map[string]CertSpec)
			for cname, cs := range certs {
				// shared endpoint trust - one key shared by subsystem endpoints
				cs.KeySubdir = et[cname].KeySubdir
				cs.KeyFile = et[cname].KeyFile
				cbmap[cname] = cs
			}
			return cbmap
		}

		if certs.PublicUserFacingCerts {
			// mgmt
			certSpec := certs.etCertSpecMgt

			// csr
			outpath = fileName(concatSubdir(outfiles[outdir], outfiles[SharedCsrOutDir]), certSpec.CsrConf)
			writeTemplate(ekuServerAuth, outpath, certSpec)

			// key-pair
			outpath = fileName(concatSubdir(outfiles[outdir], outfiles[SharedCsrOutDir]), certSpec.CsrConf + osenv.ShellExt)
			writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec})

			// ptl
			certSpec = certs.etCertSpecPtl

			// csr
			outpath = fileName(concatSubdir(outfiles[outdir], outfiles[SharedCsrOutDir]), certSpec.CsrConf)
			writeTemplate(ekuServerAuth, outpath, certSpec)

			// key-pair
			outpath = fileName(concatSubdir(outfiles[outdir], outfiles[SharedCsrOutDir]), certSpec.CsrConf + osenv.ShellExt)
			writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec})

			// combine public-user-facing key and csr scripts
			certmap := make(map[string]CertSpec)
			certmap["shared-endpoint-trust-management"] = certs.etCertSpecMgt
			certmap["shared-endpoint-trust-portal"] = certs.etCertSpecPtl

			outpath = fileName(concatSubdir(outfiles[outdir], outfiles[SharedCsrOutDir]), tagOutputFileName("all-user-facing-public-csr", tag) + osenv.ShellExt)
			writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certmap})

			// apicup certs user-facing-public
			m := make(map[string]CertSpec)
			m[CertKeyPlatformApi] = certs.etCertSpecMgt
			m[CertKeyConsumerApi] = certs.etCertSpecMgt
			m[CertKeyApiManagerUi] = certs.etCertSpecMgt
			m[CertKeyCloudAdminUi] = certs.etCertSpecMgt
			m[CertKeyPortalWwwIngress] = certs.etCertSpecPtl

			outpath = fileName(outfiles["outdir"], tagOutputFileName(outfiles[etUserFacingPublicCertsOut], tag)) + osenv.ShellExt
			writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: updateFromEndpointTrust(certs.PublicUserFacingEkuServerAuth, m)})
		}

		if certs.PublicCerts {
			// gateway
			certSpec := certs.etCertSpecGwy

			// csr
			outpath = fileName(concatSubdir(outfiles[outdir], outfiles[SharedCsrOutDir]), certSpec.CsrConf)
			writeTemplate(ekuServerAuth, outpath, certSpec)

			// key-pair
			outpath = fileName(concatSubdir(outfiles[outdir], outfiles[SharedCsrOutDir]), certSpec.CsrConf + osenv.ShellExt)
			writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec})

			// combine public key and csr scripts
			certmap := make(map[string]CertSpec)
			certmap["shared-endpoint-trust-gateway"] = certs.etCertSpecGwy

			outpath = fileName(concatSubdir(outfiles[outdir], outfiles[SharedCsrOutDir]), tagOutputFileName("all-public-csr", tag) + osenv.ShellExt)
			writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certmap})

			// apicup certs public
			m := make(map[string]CertSpec)
			m[CertKeyApicGwServiceIngress] = certs.etCertSpecGwy

			outpath = fileName(outfiles["outdir"], tagOutputFileName(outfiles[etPublicCertsOut], tag)) + osenv.ShellExt
			writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: updateFromEndpointTrust(certs.PublicEkuServerAuth, m)})
		}

		if certs.CommonCerts {
			// alyt
			certSpec := certs.etCertSpecAlyt

			// csr
			outpath = fileName(concatSubdir(outfiles[outdir], outfiles[SharedCsrOutDir]), certSpec.CsrConf)
			writeTemplate(ekuServerAuth, outpath, certSpec)

			// key-pair
			outpath = fileName(concatSubdir(outfiles[outdir], outfiles[SharedCsrOutDir]), certSpec.CsrConf + osenv.ShellExt)
			writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec})

			// combine mutual-auth key and csr scripts
			certmap := make(map[string]CertSpec)
			certmap["shared-endpoint-trust-analytics"] = certs.etCertSpecAlyt

			outpath = fileName(concatSubdir(outfiles[outdir], outfiles[SharedCsrOutDir]), tagOutputFileName("all-mutual-auth-csr", tag) + osenv.ShellExt)
			writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certmap})

			// combine common key and csr scripts
			//outpath = fileName(concatSubdir(outfiles[outdir], outfiles[CommonCsrOutDir]), tagOutputFileName("all-common-csr", tag) + osenv.ShellExt)
			//writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.CommonEkuClientAuth})

			// apicup certs mutual auth
			m := make(map[string]CertSpec)
			m[CertKeyAnalyticsClientIngress] = certs.etCertSpecAlyt
			m[CertKeyAnalyticsIngestionIngress] = certs.etCertSpecAlyt
			m[CertKeyPortalAdminIngress] = certs.etCertSpecPtl

			outpath = fileName(outfiles["outdir"], tagOutputFileName(outfiles[etMutualAuthCertsOut], tag)) + osenv.ShellExt
			writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: updateFromEndpointTrust(certs.MutualAuthEkuServerAuth,m)})
		}
	}
}
