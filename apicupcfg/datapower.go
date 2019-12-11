package apicupcfg

import (
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"os"
	"strings"
	"text/template"
)

const cryptoDir = "sharedcert"
const gwdKey = "gwd_key"
const gwdCert = "gwd_cert"
const gwdIdCred = "gwd_id_cred"
const sslGwdServer = "gwd_server"
const sslGwdClient = "gwd_client"

const defaultApiGatewayPort = 9443
const defaultApiGatewayAddress = "if_eth0"

const defaultApicGwServicePort = 3000
const defaultApicGwServiceAddress = "if_eth0"

const gwdPeering = "gwd"
const rateLimitPeering = "rate-limit"
const subsPeering = "subs"

const gwdPeeringLocalPort = 16380
const gwdPeeringMonitorPort = 26380

const rateLimitPeeringLocalPort = 16383
const rateLimitPeeringMonitorPort = 26383

const subsPeeringLocalPort = 16384
const subsPeeringMonitorPort = 26384

type DpConfigSequence struct {
	Domain string
	ConfigSequenceName string
	ConfigurationExecutionInterval int
}

func nonl(buf string) string { return buf /* strings.ReplaceAll(buf, "\n", "")*/}
func dot2dash(buf string) string { return strings.ReplaceAll(buf, ".", "-") }

type DpFile struct {
	Domain string
	Directory string
	FileName string
	FileContent string
}

func dpSetFile(outdir, outfile, dpdomain, dpdir, dpfile string, tbox *rice.Box) {

	dp := DpFile{
		Domain:      dpdomain,
		Directory:   dpdir,
		FileName:    dpfile,
		FileContent: "hello",
	}

	dpWriteTemplate(outdir, outfile, dp, "dp-set-file.tmpl", tbox)
}

type DpCryptoKey struct {
	Domain string
	CryptoKeyName string
	CryptoKeyFile string
}

func dpCryptoKey(outdir, outfile string, dp DpCryptoKey, tbox *rice.Box) {

	dpWriteTemplate(outdir, outfile, dp, "dp-crypto-key.tmpl", tbox)
}

type DpCryptoCertificate struct {
	Domain string
	CryptoCertName string
	CryptoCertFile string
}

func dpCryptoCertificate(outdir, outfile string, dp DpCryptoCertificate, tbox *rice.Box) {

	dpWriteTemplate(outdir, outfile, dp, "dp-crypto-certificate.tmpl", tbox)
}

type DpCryptoIdentCredentials struct {
	Domain string
	Name string
	CryptoKeyName string
	CryptoCertName string
	CaName string
}

func dpCryptoIndentCredentials(outdir, outfile string, dp DpCryptoIdentCredentials, tbox *rice.Box) {

	dpWriteTemplate(outdir, outfile, dp, "dp-crypto-ident-cred.tmpl", tbox)
}

type DpSSLServerProfile struct {
	Domain string
	Name string
	CryptoIdentCreds string
	CryptoValCreds string
}

func dpSSLServerProfile(outdir, outfile string, dp DpSSLServerProfile, tbox *rice.Box) {

	dpWriteTemplate(outdir, outfile, dp, "dp-ssl-server-profile.tmpl", tbox)
}

type DpSSLClientProfile struct {
	Domain string
	Name string
	CryptoIdentCreds string
	CryptoValCreds string
}

func dpSSLClientProfile(outdir, outfile string, dp DpSSLClientProfile, tbox *rice.Box) {

	dpWriteTemplate(outdir, outfile, dp, "dp-ssl-client-profile.tmpl", tbox)
}

func dpConfigSequence(outdir, outfile string, dp DpConfigSequence, tbox *rice.Box) {

	dpWriteTemplate(outdir, outfile, dp, "dp-config-sequence.tmpl", tbox)
}

type DpGatewayPeering struct {
	Domain string
	Name string
	Summary string
	LocalAddress string
	LocalPort int
	MonitorPort int
	PeerGroupSwitch string
	Peer1 string
	Peer2 string
	Priority int
	SSLSwitch string
	CryptoIdentCreds string
	CryptoValCreds string
	PersistenceLocation string
	LocalDirectory string
}

func dpGatewayPeering(outdir, outfile string, dp DpGatewayPeering, tbox *rice.Box) {
	dpWriteTemplate(outdir, outfile, dp, "dp-gateway-peering.tmpl", tbox)
}

type DpGatewayPeeringManager struct {
	Domain string
	Name string
	GwdPeering string
	RateLimitPeering string
	SubscriptionPeering string
}

func dpGatewayPeeringManager(outdir, outfile string, dp DpGatewayPeeringManager, tbox *rice.Box) {
	dpWriteTemplate(outdir, outfile, dp, "dp-gateway-peering-manager.tmpl", tbox)
}

type DpApicGwService struct {
	Domain string
	Name string
	LocalAddress string
	LocalPort int
	SSLClientProfile string
	SSLServerProfile string
	ApiGateway string
	ApiGatewayPort int
	//GwdPeering string
	GwPeeringManager string
}

func dpApicGwService(outdir, outfile string, dp DpApicGwService, tbox *rice.Box) {
	dpWriteTemplate(outdir, outfile, dp, "dp-apic-gw-service.tmpl", tbox)
}

type DpHostAlias struct {
	Alias string
	IPAddress string
}

func dpHostAlias(outdir, outfile string, dp DpHostAlias, tbox *rice.Box) {
	dpWriteTemplate(outdir, outfile, dp, "dp-host-alias.tmpl", tbox)
}

type DpNTPService struct {
	NTPServer string
}

func dpNTPService(outdir, outfile string, dp DpNTPService, tbox *rice.Box) {
	dpWriteTemplate(outdir, outfile, dp, "dp-ntp-service.tmpl", tbox)
}

type DpSystemSettings struct {
	SystemName string
}

func dpSystemSettings(outdir, outfile string, dp DpSystemSettings, tbox *rice.Box) {
	dpWriteTemplate(outdir, outfile, dp, "dp-system-settings.tmpl", tbox)
}

type DpDomain struct {
	DatapowerDomain string
}

func dpDomain(outdir, outfile string, dp DpDomain, tbox *rice.Box) {
	dpWriteTemplate(outdir, outfile, dp, "dp-domain.tmpl", tbox)
}

type DpSaveConfig struct {
	Domain string
}

func dpSaveConfig(outdir, outfile string, dp DpSaveConfig, tbox *rice.Box) {
	dpWriteTemplate(outdir, outfile, dp, "dp-save-config.tmpl", tbox)
}

type DpDNSService struct {
	SearchDomains []string
	NameServers []string
}

func dpDNSServiceModify(outdir, outfile string, dp DpDNSService, tbox *rice.Box) {
	dpWriteTemplate(outdir, outfile, dp, "dp-dns-modify.tmpl", tbox)
}

func dpWriteTemplate(outdir, outfile string, dp interface{}, templateName string, tbox *rice.Box) {
	t := parseTemplate(tbox, tpdir(tbox) + templateName)
	outpath := fileName(outdir, outfile)
	writeTemplate2(t, outpath, nonl, dp)
}

type SomaSpec struct {
	Req string		// req xml file
	File string		// set-file file path
	Dpdir string	// set-file dp directory (cert, local, etc)
	Dpfile string	// set-file dp file (gwd_key.pem, etc)
	Dpdomain string		// set-file dp domain
	Auth string		// auth env file with username, password
	Url string		// datapower management service url
}

type OsEnvSomaSpecs struct {
	OsEnv
	Config string
	SetFileSpecs []SomaSpec
	ReqSpecs []SomaSpec
}

func datapowerOpensslConfig(subsys *SubsysVm, outputdir string, osenv OsEnv, tbox *rice.Box) {

	// parse templates
	ekuServerClient := parseTemplate(tbox, tpdir(tbox) + "csr-server-client-eku.tmpl")
	keypairTemplate := parseTemplates(tbox, tpdir(tbox) + "keypair.tmpl", tpdir(tbox) + "helpers.tmpl")
	combinedCsrTemplate := parseTemplates(tbox, tpdir(tbox) + "combined-csr.tmpl", tpdir(tbox) + "helpers.tmpl")

	certmap := make(map[string]CertSpec)

	// apic gateway service
	cs := CertSpec{}
	cs.Cn = subsys.Gateway.ApicGwService
	updateCertSpec(&subsys.Certs, subsys.Gateway.SubsysName, "datapower", &cs, DatapowerOutDir)

	// add api gateway to alt cns
	cs.AltCns = append(cs.AltCns, subsys.Gateway.ApiGateway)

	for _, host := range subsys.Gateway.Hosts {
		if len(host.Name) == 0 {
			continue
		}
		cs.AltCns = append(cs.AltCns, host.Name)
	}

	// save cert-spec in combined cert-map
	certmap[cs.Cn] = cs

	// csr
	outpath := fileName(outputdir, cs.CsrConf)
	writeTemplate(ekuServerClient, outpath, cs)

	// key-pair
	outpath = fileName(outputdir, cs.CsrConf + osenv.ShellExt)
	writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: cs})

/*
	// api gateway
	cs = CertSpec{}
	cs.Cn = subsys.Gateway.ApiGateway
	updateCertSpec(&subsys.Certs, subsys.Gateway.SubsysName, "datapower", &cs, DatapowerOutDir)

	for _, host := range subsys.Gateway.Hosts {
		if len(host.Name) == 0 {
			continue
		}
		cs.AltCns = append(cs.AltCns, host.Name)
	}

	// save cert-spec in combined cert-map
	certmap[cs.Cn] = cs

	// csr
	outpath = fileName(outputdir, cs.CsrConf)
	writeTemplate(ekuServerClient, outpath, cs)

	// key-pair
	outpath = fileName(outputdir, cs.CsrConf + osenv.ShellExt)
	writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: cs})
*/
	// combine mutual-auth key and csr scripts
	outpath = fileName(outputdir, tagOutputFileName("all-datapower-csr", subsys.Tag) + osenv.ShellExt)
	writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certmap})
}

func datapowerCryptoConfig(gwy *GwySubsysVm, outputdir string, tbox *rice.Box) {

	dpdomain := gwy.GetDatapowerDomainOrDefault()

	// crypto-key
	dpcryptokey := DpCryptoKey{
		Domain:        dpdomain,
		CryptoKeyName: gwy.GetGwdKeyOrDefault(),
		CryptoKeyFile: fmt.Sprintf("%s:///%s.pem",
			gwy.GetCryptoDirectoryOrDefault(), gwy.GetGwdKeyOrDefault()),
	}
	outfile := fmt.Sprintf("%s", "dp-crypto-key.xml")
	dpCryptoKey(outputdir, outfile, dpcryptokey, tbox)

	// cryto-certificate self-signed
	dpcryptocert := DpCryptoCertificate{
		Domain:         dpdomain,
		CryptoCertName: fmt.Sprintf("%s_self", gwy.GetGwdCertOrDefault()),
		CryptoCertFile: fmt.Sprintf("%s:///%s_self.pem",
			gwy.GetCryptoDirectoryOrDefault(), gwy.GetGwdCertOrDefault()),
	}
	outfile = fmt.Sprintf("%s","dp-crypto-cert-self.xml")
	dpCryptoCertificate(outputdir, outfile, dpcryptocert, tbox)

	// crypto certificate
	dpcryptocert = DpCryptoCertificate{
		Domain:         dpdomain,
		CryptoCertName: fmt.Sprintf("%s", gwy.GetGwdCertOrDefault()),
		CryptoCertFile: fmt.Sprintf("%s:///%s.pem",
			gwy.GetCryptoDirectoryOrDefault(), gwy.GetGwdCertOrDefault()),
	}
	outfile = fmt.Sprintf("%s", "dp-crypto-cert.xml")
	dpCryptoCertificate(outputdir, outfile, dpcryptocert, tbox)

	// crypto cert - ca
	dpcryptocert = DpCryptoCertificate{
		Domain:         dpdomain,
		CryptoCertName: fmt.Sprintf("%s", gwy.GetCaCertOrDefault()),
		CryptoCertFile: fmt.Sprintf("%s:///%s.pem",
			gwy.GetCryptoDirectoryOrDefault(), gwy.GetCaCertOrDefault()),
	}
	outfile = fmt.Sprintf("%s", "dp-crypto-cert-ca.xml")
	dpCryptoCertificate(outputdir, outfile, dpcryptocert, tbox)

	// crypto cert - root ca
	dpcryptocert = DpCryptoCertificate{
		Domain:         dpdomain,
		CryptoCertName: fmt.Sprintf("%s", gwy.GetRootCertOrDefault()),
		CryptoCertFile: fmt.Sprintf("%s:///%s.pem",
			gwy.GetCryptoDirectoryOrDefault(), gwy.GetRootCertOrDefault()),
	}
	outfile = fmt.Sprintf("%s", "dp-crypto-cert-root-ca.xml")
	dpCryptoCertificate(outputdir, outfile, dpcryptocert, tbox)

	// crypto identification credentials
	dpidcreds := DpCryptoIdentCredentials{
		Domain:         dpdomain,
		Name:           gwdIdCred,
		CryptoKeyName:  gwy.GetGwdKeyOrDefault(),
		CryptoCertName: fmt.Sprintf("%s_self", gwy.GetGwdCertOrDefault()),
		CaName:         "",
	}
	outfile = fmt.Sprintf("%s", "dp-crypto-id-creds.xml")
	dpCryptoIndentCredentials(outputdir, outfile, dpidcreds, tbox)

	// valcred: gwd_val_cred -- not used...
}

func datapowerGatewayPeeringConfig(gwy *GwySubsysVm, outputdir string, tbox *rice.Box) {

	peering := []string {gwy.GetGwdPeeringOrDefault(), gwy.GetRateLimitPeeringOrDefault(), gwy.GetSubsPeeringOrDefault()}
	localPorts := []int {gwy.GetGwdPeeringLocalPortOrDefault(), gwy.GetRateLimitPeeringLocalPortOrDefault(), gwy.GetSubsPeeringLocalPortOrDefault()}
	monitorPorts := []int {gwy.GetGwdPeeringMonitorPortOrDefault(), gwy.GetRateLimitPeeringMonitorPortOrDefault(), gwy.GetSubsPeeringMonitorPortOrDefault()}

	prif := func(hidx int, host, pgroup string) int {switch hidx {case 0: return 100; case 1: return 90; case 2: return 80; default: return 50}}

	for hidx, host := range gwy.Hosts {

		if len(host.Name) == 0 {
			continue
		}

		peer1 := ""
		peer2 := ""

		peergroupswitch := "off"
		if len(gwy.Hosts) == 3 {
			peergroupswitch = "on"
		}

		sslswitch := "off"
		if peergroupswitch == "on" {
			sslswitch = "on"
		}

		for pgidx, pgroup := range peering {

			if len(gwy.Hosts) == 3 {
				if hidx == 0 {
					peer1 = gwy.Hosts[1].Name
					peer2 = gwy.Hosts[2].Name

				} else if hidx == 1 {
					peer1 = gwy.Hosts[0].Name
					peer2 = gwy.Hosts[2].Name

				} else if hidx == 2 {
					peer1 = gwy.Hosts[0].Name
					peer2 = gwy.Hosts[1].Name
				}
			}

			localport := localPorts[pgidx]
			monitorport := monitorPorts[pgidx]

			// gateway peering: gwd, rate-limit, subs, (api-probe)
			dpgwpeering := DpGatewayPeering{
				Domain:              gwy.GetDatapowerDomainOrDefault(),
				Name:                pgroup,
				Summary:             "APIC peering",
				LocalAddress:		fmt.Sprintf("if_%s", host.Device),
				LocalPort:           localport,
				MonitorPort:         monitorport,
				PeerGroupSwitch: peergroupswitch,
				Peer1:               peer1,
				Peer2:               peer2,
				SSLSwitch:           sslswitch,
				Priority:            prif(hidx, host.Name, pgroup),
				CryptoIdentCreds:       gwdIdCred,
				CryptoValCreds:      "", // do not assign validation creds...
				PersistenceLocation: "memory",
				LocalDirectory:      "local:///", // local store or raid
			}

			outfile := fmt.Sprintf("dp-peering-%s-%s.xml", pgroup, dot2dash(host.Name))
			dpGatewayPeering(outputdir, outfile, dpgwpeering, tbox)
		}
	}
}

type DpCryptoIdentCredModify struct {
	Domain string
	Name string
	Key string
	Cert string
	CaCerts []string
}

func datapowerCryptoIdentCredModify(outputdir, outfile string, dp DpCryptoIdentCredModify, tbox *rice.Box) {

	// set-file cert
	// set-file ca-cert
	// set-file root-cert
	// crypto-key key
	// crypto-cert cert
	// crypto-cert ca-cert
	// crypto-cert root-cert
	// crypto-id-creds-mod crypto-cert, crypto-ca-cert, crypto-root-cert

	dpWriteTemplate(outputdir, outfile, dp, "dp-crypto-ident-cred-modify.tmpl", tbox)
}

func somaServiceUrl(host string) string {
	return fmt.Sprintf("https://%s:5550/service/mgmt/3.0", host)
}

func somaPref() string {
	return "soma" + string(os.PathSeparator)
}

func somaCryptoSelf(gwy *GwySubsysVm, dpoutdir string, osenv OsEnv, configFileName string, somaTemplate *template.Template) {

	// set file private key
	// set file self-signed cert
	// crypto key
	// crypto cert self-signed
	// crypto ident cred

	for _, host := range gwy.Hosts {

		if len(host.Name) == 0 {
			continue
		}

		dpenv := "dp.env"
		url := somaServiceUrl(host.Name)

		setfileSpecs := make([]SomaSpec, 0)
		reqSpecs := make([]SomaSpec, 0)

		// set files - default domain

		setfileSpecs = append(setfileSpecs, SomaSpec{
			Req:    "",
			File:   dot2dash(gwy.ApicGwService) + ".key",
			Dpdir:  gwy.GetCryptoDirectoryOrDefault(),
			Dpfile: gwy.GetGwdKeyOrDefault() + ".pem",
			Dpdomain: "default",
			Auth:   dpenv,
			Url:    url,
		})

		setfileSpecs = append(setfileSpecs, SomaSpec{
			Req:    "",
			File:   dot2dash(gwy.ApicGwService) + ".crt.self",
			Dpdir:  gwy.GetCryptoDirectoryOrDefault(),
			Dpfile: fmt.Sprintf("%s_self.pem", gwy.GetGwdCertOrDefault()),
			Dpdomain: "default",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + "dp-domain.xml",
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		// save default domain
		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + fmt.Sprintf("dp-save-config-%s.xml", "default"),
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		// crypto-objects apiconnect domain
		
		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + "dp-crypto-key.xml",
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + "dp-crypto-cert-self.xml",
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + "dp-crypto-id-creds.xml",
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		// save apiconnect domain
		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + fmt.Sprintf("dp-save-config-%s.xml", gwy.GetDatapowerDomainOrDefault()),
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		// write soma script for a host
		outpath := dpoutdir + osenv.PathSeparator + "zoma-crypto-self-" + dot2dash(host.Name) + osenv.ShellExt
		writeTemplate(somaTemplate, outpath, OsEnvSomaSpecs{
			OsEnv:        osenv,
			Config:	configFileName,
			SetFileSpecs: setfileSpecs,
			ReqSpecs:     reqSpecs,
		})
	}
}

func somaCryptoUpdate(gwy *GwySubsysVm, dpoutdir string, osenv OsEnv, configFileName string, somaTemplate *template.Template) {

	// write crypto ident modify script
	// set file crypto-cert
	// set file ca cert
	// set file root cert
	// crypto cert
	// crypto cert ca
	// crypto cert root ca
	// crypto ident cred mod

	for _, host := range gwy.Hosts {

		if len(host.Name) == 0 {
			continue
		}

		dpenv := "dp.env"
		url := somaServiceUrl(host.Name)

		setfileSpecs := make([]SomaSpec, 0)
		reqSpecs := make([]SomaSpec, 0)

		setfileSpecs = append(setfileSpecs, SomaSpec{
			Req:    "",
			File:   dot2dash(gwy.ApicGwService) + ".crt",
			Dpdir:  gwy.GetCryptoDirectoryOrDefault(),
			Dpfile: fmt.Sprintf("%s.pem", gwy.GetGwdCertOrDefault()),
			Dpdomain: "default",
			Auth:   dpenv,
			Url:    url,
		})

		setfileSpecs = append(setfileSpecs, SomaSpec{
			Req:    "",
			File:   gwy.CaFile,
			Dpdir:  gwy.GetCryptoDirectoryOrDefault(),
			Dpfile: fmt.Sprintf("%s.pem", gwy.GetCaCertOrDefault()),
			Dpdomain: "default",
			Auth:   dpenv,
			Url:    url,
		})

		setfileSpecs = append(setfileSpecs, SomaSpec{
			Req:    "",
			File:   gwy.RootCaFile,
			Dpdir:  gwy.GetCryptoDirectoryOrDefault(),
			Dpfile: fmt.Sprintf("%s.pem", gwy.GetRootCertOrDefault()),
			Dpdomain: "default",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + "dp-crypto-cert.xml",
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + "dp-crypto-cert-ca.xml",
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + "dp-crypto-cert-root-ca.xml",
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + "dp-crypto-ident-cred-modify.xml",
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		// save apiconnect domain
		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + fmt.Sprintf("dp-save-config-%s.xml", gwy.GetDatapowerDomainOrDefault()),
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		// write soma script for a host
		outpath := dpoutdir + osenv.PathSeparator + "zoma-crypto-update-" + dot2dash(host.Name) + osenv.ShellExt
		writeTemplate(somaTemplate, outpath, OsEnvSomaSpecs{
			OsEnv:        osenv,
			Config:	configFileName,
			SetFileSpecs: setfileSpecs,
			ReqSpecs:     reqSpecs,
		})
	}
}

func somaApiConnect(gwy *GwySubsysVm, dpoutdir string, osenv OsEnv, configFileName string, somaTemplate *template.Template) {

	for _, host := range gwy.Hosts {

		if len(host.Name) == 0 {
			continue
		}

		dpenv := "dp.env"
		url := somaServiceUrl(host.Name)

		setfileSpecs := make([]SomaSpec, 0)
		reqSpecs := make([]SomaSpec, 0)

		// default domain

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + fmt.Sprintf("dp-system-settings-%s.xml", dot2dash(host.Name)),
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + fmt.Sprintf("dp-host-alias-%s.xml", dot2dash(host.Name)),
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + "dp-ntp-service.xml",
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:      somaPref() + "dp-dns-service-modify.xml",
			File:     "",
			Dpdir:    "",
			Dpfile:   "",
			Dpdomain: "",
			Auth:     dpenv,
			Url:      url,
		})

		// api connect domain created by the crypto script
		//reqSpecs = append(reqSpecs, SomaSpec{
		//	Req:    somaPref() + "dp-domain.xml",
		//	File:   "",
		//	Dpdir:  "",
		//	Dpfile: "",
		//	Auth:   dpenv,
		//	Url:    url,
		//})

		// save default domain
		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + fmt.Sprintf("dp-save-config-%s.xml", "default"),
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		// api connect domain

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + "dp-ssl-server.xml",
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + "dp-ssl-client.xml",
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		peeringkey := fmt.Sprintf(fmt.Sprintf("dp-peering-%s-%s.xml", gwdPeering, dot2dash(host.Name)))
		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + peeringkey,
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		peeringkey = fmt.Sprintf(fmt.Sprintf("dp-peering-%s-%s.xml", rateLimitPeering, dot2dash(host.Name)))
		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + peeringkey,
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		peeringkey = fmt.Sprintf(fmt.Sprintf("dp-peering-%s-%s.xml", subsPeering, dot2dash(host.Name)))
		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + peeringkey,
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + "dp-peering-manager.xml",
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + "dp-config-sequence.xml",
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + "dp-apic-gw-service.xml",
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		reqSpecs = append(reqSpecs, SomaSpec{
			Req:    somaPref() + fmt.Sprintf("dp-save-config-%s.xml", gwy.GetDatapowerDomainOrDefault()),
			File:   "",
			Dpdir:  "",
			Dpfile: "",
			Auth:   dpenv,
			Url:    url,
		})

		// write soma script for a host
		outpath := dpoutdir + osenv.PathSeparator + "zoma-" + dot2dash(host.Name) + osenv.ShellExt
		writeTemplate(somaTemplate, outpath, OsEnvSomaSpecs{
			OsEnv:        osenv,
			Config:	configFileName,
			SetFileSpecs: setfileSpecs,
			ReqSpecs:     reqSpecs,
		})
	}
}

func datapowerCluster(subsys *SubsysVm, outfiles map[string]string, tbox *rice.Box) {

	// parse templates
	somaTemplate := parseTemplates(tbox, tpdir(tbox) + "soma.tmpl", tpdir(tbox) + "helpers.tmpl")

	var osenv OsEnv
	osenv.init()

	// datapower output directory
	dpoutdir := concatSubdir(outfiles[outdir], outfiles[DatapowerOutDir])

	// soma output directory
	somaoutdir := concatSubdir(dpoutdir, "soma")

	//outdir1 := concatSubdir(outfiles[outdir], outfiles[DatapowerOutDir])

	// openssl configuration
	datapowerOpensslConfig(subsys, dpoutdir, osenv, tbox)

	gwy := subsys.Gateway

	// datapower domain
	dpdomain := gwy.GetDatapowerDomainOrDefault()

	// host alias, system name
	for _, host := range subsys.Gateway.Hosts {
		if len(host.Name) == 0 {
			continue
		}

		dpha := DpHostAlias{
			Alias:     fmt.Sprintf("if_%s", host.Device),
			IPAddress: host.IpAddress,
		}

		outfile := fmt.Sprintf("dp-host-alias-%s.xml", dot2dash(host.Name))
		dpHostAlias(somaoutdir, outfile, dpha, tbox)

		dpsys := DpSystemSettings{SystemName: host.Name}
		outfile = fmt.Sprintf("dp-system-settings-%s.xml", dot2dash(host.Name))
		dpSystemSettings(somaoutdir, outfile, dpsys, tbox)
	}

	// ntp service
	ntpserver := gwy.GetNTPServerOrDefault()
	dpntp := DpNTPService{NTPServer: ntpserver}
	outfile := fmt.Sprintf("%s", "dp-ntp-service.xml")
	dpNTPService(somaoutdir, outfile, dpntp, tbox)

	// update dns service
	dpdns := DpDNSService{
		SearchDomains: subsys.Gateway.SearchDomains,
		NameServers:   subsys.Gateway.DnsServers,
	}

	outfile = fmt.Sprintf("%s", "dp-dns-service-modify.xml")
	dpDNSServiceModify(somaoutdir, outfile, dpdns, tbox)

	// application domain
	dpdatapowerdomain := DpDomain{DatapowerDomain: gwy.GetDatapowerDomainOrDefault()}
	outfile = fmt.Sprintf("%s", "dp-domain.xml")
	dpDomain(somaoutdir, outfile, dpdatapowerdomain, tbox)

	// crypto-key, cryto-certs, identification creds
	datapowerCryptoConfig(&gwy, somaoutdir, tbox)

	// ssl-server
	dpsslsrv := DpSSLServerProfile{
		Domain:               dpdomain,
		Name:                 sslGwdServer, // gwd_server
		CryptoIdentCreds: gwdIdCred,
		CryptoValCreds:   "", // no valcreds...
	}

	outfile = fmt.Sprintf("%s", "dp-ssl-server.xml")
	dpSSLServerProfile(somaoutdir, outfile, dpsslsrv, tbox)

	// ssl-client
	dpsslclient := DpSSLClientProfile{
		Domain:               dpdomain,
		Name:                 sslGwdClient, // gwd_client
		CryptoIdentCreds: gwdIdCred,
		CryptoValCreds:   "",
	}

	outfile = fmt.Sprintf("%s", "dp-ssl-client.xml")
	dpSSLClientProfile(somaoutdir, outfile, dpsslclient, tbox)

	// gateway peering
	datapowerGatewayPeeringConfig(&gwy, somaoutdir, tbox)

	// gateway peering manager: default
	peeringmgr := DpGatewayPeeringManager{
		Domain:              dpdomain,
		Name:                "default", // always
		GwdPeering:          gwdPeering,
		RateLimitPeering:    rateLimitPeering,
		SubscriptionPeering: subsPeering,
	}

	outfile = fmt.Sprintf("%s","dp-peering-manager.xml")
	dpGatewayPeeringManager(somaoutdir, outfile, peeringmgr, tbox)

	// configuration sequence
	dpconfigseq := DpConfigSequence{
		Domain:                         dpdomain,
		ConfigSequenceName:             "apiconnect", // always
		ConfigurationExecutionInterval: gwy.ConfigurationExecutionInterval,
	}

	outfile = fmt.Sprintf("%s","dp-config-sequence.xml")
	dpConfigSequence(somaoutdir, outfile, dpconfigseq, tbox)

	// apic gw service: default
	apicgw := DpApicGwService{
		Domain:           dpdomain,
		Name:             "default", // always
		LocalAddress:     gwy.GetApicGwServiceAddressOrDefault(),
		LocalPort:        gwy.GetApicGwServicePortOrDefault(),
		SSLClientProfile: sslGwdClient,
		SSLServerProfile: sslGwdServer,
		ApiGateway:       gwy.GetApiGatewayAddressOrDefault(),
		ApiGatewayPort:   gwy.GetApiGatewayPortOrDefault(),
		//GwdPeering:       gwdPeering, // v5 only, not used for the api gateway
		GwPeeringManager: "default", // always
	}

	outfile = fmt.Sprintf("%s", "dp-apic-gw-service.xml")
	dpApicGwService(somaoutdir, outfile, apicgw, tbox)

	// save config - default domain
	dpsavedomain := "default"
	dpsaveconfig := DpSaveConfig{Domain: dpsavedomain}
	outfile = fmt.Sprintf("dp-save-config-%s.xml", dpsavedomain)
	dpSaveConfig(somaoutdir, outfile, dpsaveconfig, tbox)

	// save config - api connect domain
	dpsavedomain = gwy.GetDatapowerDomainOrDefault()
	dpsaveconfig = DpSaveConfig{Domain: dpsavedomain}
	outfile = fmt.Sprintf("dp-save-config-%s.xml", dpsavedomain)
	dpSaveConfig(somaoutdir, outfile, dpsaveconfig, tbox)

	// crypto ident cred modify
	dpcryptoidentcredmod := DpCryptoIdentCredModify{
		Domain:  dpdomain,
		Name:    gwdIdCred,
		Key:     fmt.Sprintf("%s", gwy.GetGwdKeyOrDefault()),
		Cert:    fmt.Sprintf("%s", gwy.GetGwdCertOrDefault()),
		CaCerts: []string{gwy.GetCaCertOrDefault(),gwy.GetRootCertOrDefault()},
	}

	outfile = fmt.Sprintf("%s", "dp-crypto-ident-cred-modify.xml")
	datapowerCryptoIdentCredModify(somaoutdir, outfile, dpcryptoidentcredmod, tbox)

	// write out soma scripts...

	// soma crypto self-signed scripts
	somaCryptoSelf(&subsys.Gateway, dpoutdir, osenv, subsys.configFileName, somaTemplate)

	// soma api connect scripts
	somaApiConnect(&subsys.Gateway, dpoutdir, osenv, subsys.configFileName, somaTemplate)

	// soma crypto update scripts
	somaCryptoUpdate(&subsys.Gateway, dpoutdir, osenv, subsys.configFileName, somaTemplate)
}
