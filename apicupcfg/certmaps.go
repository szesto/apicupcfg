package apicupcfg

type CertInput struct {
	DnFields []string
	K8sNamespace string

	PublicUserFacingCerts bool
	PublicCerts bool // datapower api invoke
	CommonCerts bool

	BringYourOwnKey bool
}

type CertMaps struct {
	CertInput

	PublicUserFacingServerAuth map[string]CertSpec
	PublicServerAuth map[string]CertSpec
	InternalServerAuth map[string]CertSpec
	CommonClientAuth map[string]CertSpec
}

/*
func (certmaps *CertMaps) createMaps() {
	if certmaps.PublicServerAuth == nil {
		certmaps.PublicServerAuth = make(map[string]CertSpec)
	}

	if certmaps.PublicServerAuth == nil {
		certmaps.PublicServerAuth = make(map[string]CertSpec)
	}

	if certmaps.InternalServerAuth == nil {
		certmaps.InternalServerAuth = make(map[string]CertSpec)
	}

	if certmaps.CommonClientAuth == nil {
		certmaps.CommonClientAuth = make(map[string]CertSpec)
	}
}

func (certmaps *CertMaps) updateCertSpecs(mgmt ManagementSubsysDescriptor, alyt AnalyticsSubsysDescriptor, ptl PortalSubsysDescriptor, gwy GatewaySubsysDescriptor) {
	certmaps.createMaps()

	isManagement := mgmt != nil && len(mgmt.GetManagementSubsysName()) > 0

	getCertSpec := func (certmap map[string]CertSpec, key string) CertSpec {
		if certspec, ok := certmap[key]; ok {return certspec } else {return CertSpec{}}
	}

	if mgmt != nil && isManagement && certmaps.PublicUserFacingCerts {
		certmap := certmaps.PublicUserFacingServerAuth
		certspec := getCertSpec(certmap, CertKeyPlatformApi)
		certspec.Cn = mgmt.GetPlatformApiEndpoint()

		updateCertSpec(certs, mgmt.GetManagementSubsysName(), CertKeyPlatformApi, &certSpec, customCsrOutDir)
		certmap[CertKeyPlatformApi] = certSpec
	}
}

func (certSpec *CertSpec) updateCertSpec(subsysName string, certName string, certMaps *CertMaps) {

	// check if no cn...
	if len(certSpec.Cn) == 0 {
		return
	}

	if len(certSpec.CertName) == 0 {
		certSpec.CertName = certName
	}

	if len(certSpec.DnFields) == 0 {
		certSpec.DnFields = certMaps.DnFields
	}

	if len(certSpec.K8sNamespace) == 0 {
		certSpec.K8sNamespace = certMaps.K8sNamespace
	}

	// ca chain file: issuing-ca + root-ca
	if len(certSpec.CaFile) == 0 {
		//certSpec.CaFile = certs.CaFile
		certSpec.CaFile = dot2dash(certSpec.Cn) + ".chain.pem"
	}

	certSpec.SubsysName = subsysName

	// default file names
	if len(certSpec.CsrConf) == 0 {
		certSpec.CsrConf = dot2dash(certSpec.Cn) + ".conf"
	}

	if len(certSpec.KeyFile) == 0 {
		certSpec.KeyFile = dot2dash(certSpec.Cn) + ".key"
	}

	if len(certSpec.CertFile) == 0 {
		certSpec.CertFile = dot2dash(certSpec.Cn) + ".crt"
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

	// alt-cns
	certSpec.AltCns = make([]string, 0, 500)
}

func (certmaps *CertMaps) outputCerts(outfiles Outfiles, tag string, version string, useVersion bool, tbox *rice.Box) {

	// parse templates
	serverAuthTemplate := parseTemplate(tbox, tpdir(tbox) + "csr-server-auth.tmpl")
	clientAuthTemplate := parseTemplate(tbox, tpdir(tbox) + "csr-client-auth.tmpl")
	keypairTemplate := parseTemplates(tbox, tpdir(tbox) + "keypair.tmpl", tpdir(tbox) + "helpers.tmpl")
	combinedCsrTemplate := parseTemplates(tbox, tpdir(tbox) + "combined-csr.tmpl", tpdir(tbox) + "helpers.tmpl")
	subsysCertsTemplate := parseTemplates(tbox, tpdir(tbox) + "subsys-certs.tmpl", tpdir(tbox) + "helpers.tmpl")

	// setup os env for template output
	var osenv OsEnv
	osenv.init2(version, useVersion)

	if certmaps.PublicUserFacingCerts {
		for _, certSpec := range certmaps.PublicUserFacingServerAuth {
			if len(certSpec.Cn) > 0 {
				// csr
				outpath := fileName(outfiles.CustomCsrOutDir(), certSpec.CsrConf)
				writeTemplate(serverAuthTemplate, outpath, certSpec)

				// key-pair
				outpath = fileName(outfiles.CustomCsrOutDir(), certSpec.CsrConf + osenv.ShellExt)
				writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec, Passive: certmaps.BringYourOwnKey})
			}
		}
	}

	if certmaps.PublicCerts {
		for _, certSpec := range certmaps.PublicServerAuth {

			if len(certSpec.Cn) > 0 {
				// csr
				outpath := fileName(outfiles.CustomCsrOutDir(), certSpec.CsrConf)
				writeTemplate(serverAuthTemplate, outpath, certSpec)

				// key-pair
				outpath = fileName(outfiles.CustomCsrOutDir(), certSpec.CsrConf + osenv.ShellExt)
				writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec, Passive: certmaps.BringYourOwnKey})
			}
		}
	}

	if certmaps.CommonCerts {
		for _, certSpec := range certmaps.InternalServerAuth {

			if len(certSpec.Cn) > 0 {
				// csr
				outpath := fileName(outfiles.CommonCsrOutDir(), certSpec.CsrConf)
				writeTemplate(serverAuthTemplate, outpath, certSpec)

				// key-pair
				outpath = fileName(outfiles.CommonCsrOutDir(), certSpec.CsrConf + osenv.ShellExt)
				writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec, Passive: certmaps.BringYourOwnKey})
			}
		}

		for _, certSpec := range certmaps.CommonClientAuth {

			if len(certSpec.Cn) > 0 {
				// csr
				outpath := fileName(outfiles.CommonCsrOutDir(), certSpec.CsrConf)
				writeTemplate(clientAuthTemplate, outpath, certSpec)

				// key-pair
				outpath = fileName(outfiles.CommonCsrOutDir(), certSpec.CsrConf + osenv.ShellExt)
				writeTemplate(keypairTemplate, outpath, OsEnvCert{OsEnv: osenv, CertSpec: certSpec, Passive: certmaps.BringYourOwnKey})
			}
		}
	}

	if certmaps.PublicUserFacingCerts {
		// combine public-user-facing key and csr scripts
		outpath := tagOutputFileName(outfiles.AllPublicUserFacingCsrFileName(), tag) + osenv.ShellExt
		writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certmaps.PublicUserFacingServerAuth})
	}

	if certmaps.PublicCerts {
		// combine public key and csr scripts
		outpath := tagOutputFileName(outfiles.AllPublicCsrFileName(), tag) + osenv.ShellExt
		writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certmaps.PublicServerAuth})
	}

	if certmaps.CommonCerts {
		// combine internal server key and csr scripts
		outpath :=  tagOutputFileName(outfiles.AllInternalCsrFileName(), tag) + osenv.ShellExt
		writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certmaps.InternalServerAuth})

		// combine common key and csr scripts
		outpath = tagOutputFileName(outfiles.AllCommonCsrFileName(), tag) + osenv.ShellExt
		writeTemplate(combinedCsrTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certmaps.CommonClientAuth})
	}

	if certmaps.PublicUserFacingCerts {
		// apicup certs user-facing-public
		outpath := tagOutputFileName(outfiles.UserFacingPublicCertsOutFileName(), tag) + osenv.ShellExt
		writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certmaps.PublicUserFacingServerAuth})
	}

	if certmaps.PublicCerts {
		// apicup certs public
		outpath := tagOutputFileName(outfiles.PublicCertsOutFileName(), tag) + osenv.ShellExt
		writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certmaps.PublicServerAuth})
	}

	if certmaps.CommonCerts {
		// apicup certs internal
		outpath := tagOutputFileName(outfiles.AllInternalCsrFileName(), tag) + osenv.ShellExt
		writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certmaps.InternalServerAuth})

		// apicup certs common
		outpath = tagOutputFileName(outfiles.AllCommonCsrFileName(), tag) + osenv.ShellExt
		writeTemplate(subsysCertsTemplate, outpath, OsEnvCerts{OsEnv: osenv, CertSpecs: certmaps.CommonClientAuth})
	}
}
 */
