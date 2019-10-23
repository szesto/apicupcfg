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

func subsysNameForCertName(certName string, managementSubsysName string, analyticsSubsysName string,
	portalSubsysName string, gatewaySusbsysName string) (string, bool) {

	if
		certName == CertKeyPlatformApi ||
		certName == CertKeyConsumerApi ||
		certName == CertKeyApiManagerUi ||
		certName == CertKeyCloudAdminUi {

		return managementSubsysName, true
	} else if
		certName == 	CertKeyPortalAdminIngress ||
		certName == CertKeyPortalWwwIngress {

		return portalSubsysName, true
	} else if
		certName == CertKeyAnalyticsClientIngress ||
		certName == CertKeyAnalyticsIngestionIngress {

		return analyticsSubsysName, true
	} else if
		certName == CertKeyApicGwServiceIngress {

		return gatewaySusbsysName, true
	} else if
		certName == CertKeyPortalClient ||
		certName == CertKeyAnalyticsClientClient ||
		certName ==  CertKeyAnalyticsIngestionClient {

		return managementSubsysName, true
	}

	return "", false
}

func updateCertSpecs(certs Certs, managementSubsysName string, analyticsSubsysName string,
	portalSubsysName string, gatewaySusbsysName string,
	commonCsrSubdir string, publicCsrSubdir string) Certs {

	// output certs
	outCerts := certs

	// public-user-facing-eku-server-auth
	for certName, certSpec := range certs.PublicUserFacingEkuServerAuth {

		if subsysName, found := subsysNameForCertName(certName, managementSubsysName,
			analyticsSubsysName, portalSubsysName, gatewaySusbsysName); found == true {

			updateCertSpec(&certs, subsysName, certName, &certSpec, publicCsrSubdir)

			// copy updated spec
			outCerts.PublicUserFacingEkuServerAuth[certName] = certSpec
		}
	}

	// public certs
	for certName, certSpec := range certs.PublicEkuServerAuth {

		if subsysName, found := subsysNameForCertName(certName, managementSubsysName,
			analyticsSubsysName, portalSubsysName, gatewaySusbsysName); found == true {

			updateCertSpec(&certs, subsysName, certName, &certSpec, publicCsrSubdir)

			// copy updated spec
			outCerts.PublicEkuServerAuth[certName] = certSpec
		}
	}

	// mutual auth certs
	for certName, certSpec := range certs.MutualAuthEkuServerAuth {

		if subsysName, found := subsysNameForCertName(certName, managementSubsysName,
			analyticsSubsysName, portalSubsysName, gatewaySusbsysName); found == true {

			updateCertSpec(&certs, subsysName, certName, &certSpec, commonCsrSubdir)

			// copy updated spec
			outCerts.MutualAuthEkuServerAuth[certName] = certSpec
		}
	}

	// common client certs
	for certName, certSpec := range certs.CommonEkuClientAuth {
		// common certs are set on the management subsystem
		updateCertSpec(&certs, managementSubsysName, certName, &certSpec, commonCsrSubdir)

		// copy updated spec
		outCerts.CommonEkuClientAuth[certName] = certSpec
	}

	outCerts.OsEnv.init()
	return outCerts
}

func outputCerts(certs *Certs, outfiles map[string]string, tbox *rice.Box) {

	ekuServerAuth := parseTemplate(tbox, tpdir(tbox) + "csr-server-auth.tmpl")
	ekuClientAuth := parseTemplate(tbox, tpdir(tbox) + "csr-client-auth.tmpl")
	keypairTemplate := parseTemplates(tbox, tpdir(tbox) + "keypair.tmpl", tpdir(tbox) + "helpers.tmpl")
	combinedCsrTemplate := parseTemplates(tbox, tpdir(tbox) + "combined-csr.tmpl", tpdir(tbox) + "helpers.tmpl")

	subsysCertsTemplate := parseTemplates(tbox, tpdir(tbox) + "subsys-certs.tmpl", tpdir(tbox) + "helpers.tmpl")

	var outpath string

	var osenv OsEnv
	osenv.init()

	// public-user-facing-eku-server-auth csr conf
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

	// combine public-user-facing key and csr scripts
	outpath = fileName(concatSubdir(outfiles[outdir], outfiles[customCsrOutDir]), "all-user-facing-public-csr" + osenv.ShellExt)
	writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.PublicUserFacingEkuServerAuth})

	// combine public key and csr scripts
	outpath = fileName(concatSubdir(outfiles[outdir], outfiles[customCsrOutDir]), "all-public-csr" + osenv.ShellExt)
	writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.PublicEkuServerAuth})

	// combine mutual-auth key and csr scripts
	outpath = fileName(concatSubdir(outfiles[outdir], outfiles[commonCsrOutDir]), "all-mutual-auth-csr" + osenv.ShellExt)
	writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.MutualAuthEkuServerAuth})

	// combine common key and csr scripts
	outpath = fileName(concatSubdir(outfiles[outdir], outfiles[commonCsrOutDir]), "all-common-csr" + osenv.ShellExt)
	writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.CommonEkuClientAuth})

	// apicup certs user-facing-public
	outpath = fileName(outfiles["outdir"], outfiles[userFacingPublicCertsOut]) + osenv.ShellExt
	writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.PublicUserFacingEkuServerAuth})

	// apicup certs public
	outpath = fileName(outfiles["outdir"], outfiles[publicCertsOut]) + osenv.ShellExt
	writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.PublicEkuServerAuth})

	// apicup certs mutual auth
	outpath = fileName(outfiles["outdir"], outfiles[mutualAuthCertsOut]) + osenv.ShellExt
	writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.MutualAuthEkuServerAuth})

	// apicup certs common
	outpath = fileName(outfiles["outdir"], outfiles[commonCertsOut]) + osenv.ShellExt
	writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certs.CommonEkuClientAuth})

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
		cbmap := updateFromCertbot(certs.PublicUserFacingEkuServerAuth, certs.Certbot)
		outpath = fileName(outfiles["outdir"], outfiles[certbotUserFacingPublicCertOut]) + osenv.ShellExt
		writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: cbmap})

		// apicup certs public certbot
		cbmap = updateFromCertbot(certs.PublicEkuServerAuth, certs.Certbot)
		outpath = fileName(outfiles["outdir"], outfiles[certbotPublicCertOut]) + osenv.ShellExt
		writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: cbmap})
	}
}
