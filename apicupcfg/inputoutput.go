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
const certbotUserFacingPublicCertOut = "certbotuserfacingpubliccert"
const certbotPublicCertOut = "certbotpubliccert"

const CommonCsrOutDir = "common-csr"
const CustomCsrOutDir = "custom-csr"
const ProjectOutDir = "project"

func OutputFiles(baseout string) map[string]string {

	outfiles := map[string]string{
		outdir:                         baseout,
		managementOut:                  "apicup-subsys-set-management",
		gatewayOut:                     "apicup-subsys-set-gateway",
		analyticsOut:                   "apicup-subsys-set-analytics",
		portalOut:                      "apicup-subsys-set-portal",
		userFacingPublicCertsOut:       "apicup-certs-set-user-facing-public",
		publicCertsOut:                 "apicup-certs-set-public",
		mutualAuthCertsOut:             "apicup-certs-set-mutual-auth",
		commonCertsOut:                 "apicup-certs-set-common",
		CommonCsrOutDir:                CommonCsrOutDir,
		CustomCsrOutDir:                CustomCsrOutDir,
		certbotUserFacingPublicCertOut: "apicup-certs-set-certbot-user-facing-public",
		certbotPublicCertOut:           "apicup-certs-set-certbot-public",
	}

	return outfiles
}

func tagOutputFileName(outName string, tag string) string {
	const dot = "."
	return outName + dot + tag
}

func concatSubdir(dir1 string, dir2 string) string {
	return dir1 + string(os.PathSeparator) + dir2
}

func Input() (input string, outdir string, validateIp bool, initConfig bool, initConfigType string,
	subsysOnly bool, certsOnly bool, certcopy string, certdir string) {

	// define command line flags
	inputArg := flag.String("config", "subsys-config.json", "-config input-file")
	outdirArg := flag.String("out", "output", "-out output-directory")
	//commonCsrSubdirArg := flag.String("commoncsr", "common-csr", "-commoncsr subdir")
	//customCsrSubdirArg := flag.String("customcsr", "custom-csr", "-customcsr subdir")
	//projectSubdirArg := flag.String("project", "project", "-project subdir")

	validateIpArg := flag.Bool("validateip", false, "-validateip [true] validate ip addresses")

	initConfigArg := flag.Bool("initconfig", false, "-initconfig [true] initalize json config, use with configtype option")
	initConfigTypeArg := flag.String("configtype", "ova", "-configtype ova|k8s use with initconfig option")

	subsysOnlyArg := flag.Bool("subsys", false, "-subsys [true] generate subsystem scripts only")
	certsOnlyArg := flag.Bool("certs", false, "-certs [true] generate certs scripts only")

	certCopyArg := flag.String("certcopy", "", "-certcopy certfile copy certificate to destination")
	certDirArg := flag.String("certdir", "","-certdir dir copy all certificate files in dir to destination")

	// scan command line args
	flag.Parse()

	input = *inputArg
	outdir = *outdirArg
	//commonCsrSubdir = commonCsrOutDir
	//customCsrSubdir = customCsrOutDir
	//projectSubdir = projectOutDir
	validateIp = *validateIpArg
	initConfig = *initConfigArg
	initConfigType = *initConfigTypeArg
	subsysOnly = *subsysOnlyArg
	certsOnly = *certsOnlyArg
	certcopy = *certCopyArg
	certdir = *certDirArg

	return input, outdir, validateIp, initConfig, initConfigType, subsysOnly, certsOnly, certcopy, certdir
}
