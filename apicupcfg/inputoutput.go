package apicupcfg

import (
	"flag"
	"os"
)

const outdir = "outdir"
const managementOut = "management"
const gatewayOut = "gateway"
const analyticsOut = "analytics"
const portalOut = "portal"
const userFacingPublicCertsOut = "userfacingpubliccerts"
const publicCertsOut = "publiccerts"
const mutualAuthCertsOut = "mutualauthcerts"
const commonCertsOut = "commoncerts"
const commonCsrOutDir = "common-csr"
const customCsrOutDir = "custom-csr"
const certbotUserFacingPublicCertOut = "certbotuserfacingpubliccert"
const certbotPublicCertOut = "certbotpubliccert"

func OutputFiles(baseout string, commonCsrSubdir string, customCsrSubdir string) map[string]string {

	outfiles := map[string]string{
		outdir:                   baseout,
		managementOut:            "apicup-subsys-set-management",
		gatewayOut:               "apicup-subsys-set-gateway",
		analyticsOut:             "apicup-subsys-set-analytics",
		portalOut:                "apicup-subsys-set-portal",
		userFacingPublicCertsOut: "apicup-certs-set-user-facing-public",
		publicCertsOut:           "apicup-certs-set-public",
		mutualAuthCertsOut:       "apicup-certs-set-mutual-auth",
		commonCertsOut:           "apicup-certs-set-common",
		commonCsrOutDir:          commonCsrSubdir,
		customCsrOutDir:          customCsrSubdir,
		certbotUserFacingPublicCertOut: "apicup-certs-set-certbot-user-facing-public",
		certbotPublicCertOut: "apicup-certs-set-certbot-public",
	}

	return outfiles
}

func updateOutputFileName(outName string, subsysName string) string {
	const dot = "."
	return outName + dot + subsysName
}

func concatSubdir(dir1 string, dir2 string) string {
	return dir1 + string(os.PathSeparator) + dir2
}

func Input() (input string, outdir string, commonCsrSubdir string, customCsrSubdir string, projectSubdir string, validateIp bool,
	initConfig bool, initConfigType string) {

	// define command line flags
	inputArg := flag.String("config", "subsys-config.json", "-config input-file")
	outdirArg := flag.String("out", "output", "-out output-directory")
	commonCsrSubdirArg := flag.String("commoncsr", "common-csr", "-commoncsr subdir")
	customCsrSubdirArg := flag.String("customcsr", "custom-csr", "-customcsr subdir")
	projectSubdirArg := flag.String("project", "project", "-project subdir")

	validateIpArg := flag.Bool("validateip", false, "-validateip true|false")

	initConfigArg := flag.Bool("initconfig", false, "-initconfig true")
	initConfigTypeArg := flag.String("configtype", "ova", "-configtype ova|k8s")

	// scan command line args
	flag.Parse()

	input = *inputArg
	outdir = *outdirArg
	commonCsrSubdir = *commonCsrSubdirArg
	customCsrSubdir = *customCsrSubdirArg
	projectSubdir = *projectSubdirArg
	validateIp = *validateIpArg
	initConfig = *initConfigArg
	initConfigType = *initConfigTypeArg

	return input, outdir, commonCsrSubdir, customCsrSubdir, projectSubdir, validateIp, initConfig, initConfigType
}
