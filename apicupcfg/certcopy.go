package apicupcfg

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func copyCerts(certdir string, certs *Certs, mgmt ManagementSubsysDescriptor, alyt AnalyticsSubsysDescriptor,
	ptl PortalSubsysDescriptor, gwy GatewaySubsysDescriptor, commonCsrOutDir string, customCsrOutDir string, isOva bool) error {

	// read certs dir
	infos, err := ioutil.ReadDir(certdir)
	if err != nil {
		return err
	}

	for _, info := range infos {

		if info.IsDir() {
			continue
		}

		certfile := certdir + string(os.PathSeparator) + info.Name()

		err = copyCert(certfile, certs, mgmt, alyt, ptl, gwy, commonCsrOutDir, customCsrOutDir, isOva)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	return nil
}

func CopyCertDir(certdir, trustdir string, certs *Certs, mgmt ManagementSubsysDescriptor,
	alyt AnalyticsSubsysDescriptor, ptl PortalSubsysDescriptor, gwy GatewaySubsysDescriptor,
	commonCsrOutDir string, customCsrOutDir string, isOva bool) error {

	msubjfile := make(map[string]string)
	mcerts := make(map[string]*x509.Certificate)
	mcacerts := make(map[string]*x509.Certificate)
	mrootcerts := make(map[string]*x509.Certificate)

	parseCertFilef := func(certfile string) (*x509.Certificate, string, error) {

		certlist, blockTypes, err := ParseCertFile2(certfile)
		if err != nil {
			return nil, "", err
		}

		// it is possible that certlist is empty eg for private key file
		if len(certlist) == 0 {
			fmt.Printf("file '%s' is not a cert, skip...\n", certfile)
			return nil, "", nil
		}

		if len(certlist) > 1 {
			fmt.Printf("file '%s' is a cert chain... %v, skip", certfile, blockTypes)
			return nil, "", nil
		}

		return certlist[0], blockTypes[0], nil
	}

	saveCertf := func(certfile string, cert *x509.Certificate, blockType string) {

		if _, ok := msubjfile[cert.Subject.CommonName]; ok {
			fmt.Printf("cert file '%s' is a duplicate, skip...\n", certfile)
			return
		}

		if cert.IsCA {
			if cert.Subject.CommonName == cert.Issuer.CommonName {
				fmt.Printf("found Root CA cert '%s', block-type '%s'\n", certfile, blockType)

				msubjfile[cert.Subject.CommonName] = certfile
				mrootcerts[cert.Subject.CommonName] = cert

			} else {
				fmt.Printf("found CA cert '%s', block-type '%s'\n", certfile, blockType)

				msubjfile[cert.Subject.CommonName] = certfile
				mcacerts[cert.Subject.CommonName] = cert
			}

		} else {
			fmt.Printf("found cert '%s', block-type '%s'\n", certfile, blockType)

			msubjfile[cert.Subject.CommonName] = certfile
			mcerts[cert.Subject.CommonName] = cert
		}
	}

	processDirf := func(certdir string) error {

		infos, err := ioutil.ReadDir(certdir)
		if err != nil {
			return err
		}

		for _, info := range infos {

			if info.IsDir() {
				continue
			}

			// look for certs, issuing-ca certs, and root-ca certs

			// parse cert file
			certfile := certdir + string(os.PathSeparator) + info.Name()
			cert, blockType, err := parseCertFilef(certfile)

			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}

			if cert == nil {
				continue
			}

			// save cert
			saveCertf(certfile, cert, blockType)
		}

		return nil
	}

	if err := processDirf(certdir); err != nil {
		return err
	}

	if len(trustdir) > 0 {
		if err := processDirf(trustdir); err != nil {
			return err
		}
	}

	//fmt.Printf("msubjfile... %v\n", msubjfile)
	//fmt.Printf("mcerts... %v\n", mcerts)
	//fmt.Printf("mcacerts... %v\n", mcacerts)
	//fmt.Printf("mrootcerts... %v\n", mrootcerts)

	// iterate over certs and try to find a path to a leaf node

	fmt.Printf("\nsearching for trust paths...\n\n")

	isleaf := false

	for subj, cert := range mcerts {
		fmt.Printf("cert subject: %s\n", subj)

		// start at the current cert
		isleaf = false

		for casubj, cacert := range mcacerts {
			if isleaf {
				// got to the leaf for the current cert
				break
			}

			if cert.Issuer.CommonName == casubj {
				for rootsubj, rootcert := range mrootcerts {
					if isleaf {
						// got to the leaf for the current cert
						break
					}

					if cacert.Issuer.CommonName == rootsubj {
						// got to the leaf...
						fmt.Printf("leaf root subject: %s\n", rootcert.Subject.CommonName)

						// copy cert chain
						certfile := msubjfile[subj]
						cafile := msubjfile[casubj]
						rootcafile := msubjfile[rootsubj]

						fmt.Printf("copy-cert-chain args: '%s', '%s', '%s'", certfile, cafile, rootcafile)

						err := CopyCertChain(certfile, cafile, rootcafile,
							certs, mgmt, alyt, ptl, gwy,
							commonCsrOutDir, customCsrOutDir, isOva)

						if err != nil {
							fmt.Printf("%v\n", err)
						}

						// at the leaf for the current cert
						isleaf = true
					}
				}
			}
		}
	}

	return nil
}

func logCertCopy(dir, logm string) error {

	logf, err := openFileAppend(dir + string(os.PathSeparator) + "cert-copy.log")
	if err != nil {
		return err
	}

	defer func() {_ = logf.Close()}()

	_, _ = fmt.Fprintf(logf, "%s: %s\n", time.Now().Format(time.RFC822), logm)
	return nil
}

func concatChainFile(cafile, rootcafile, chainfile string, cacert, rootcert *x509.Certificate) error {

	// check if destination chain file exists
	exists, err := isFileExist(chainfile)
	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("skip cert chain merge, destination file '%s' already exists...\n", chainfile)

	} else {
		concatFiles(cafile, rootcafile, chainfile)

		logm := fmt.Sprintf("merged issuing-ca cert file '%s' (%s<-%s) and root-ca cert file '%s' (%s<-%s) into ca chain file '%s'",
			cafile, cacert.Subject.CommonName, cacert.Issuer.CommonName,
			rootcafile, rootcert.Subject.CommonName, rootcert.Issuer.CommonName, chainfile)

		fmt.Printf("%s\n", logm)

		if err := logCertCopy(filepath.Dir(chainfile), logm); err != nil {
			return err
		}
	}

	return nil
}

func copyCaFile(cacertfile, dstfile string, cacert *x509.Certificate, carole string) error {

	// check if destination file exists
	exists, err := isFileExist(dstfile)
	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("skip copying %s ca cert file '%s' (%s<-%s)... destination file '%s' already exists...\n",
			carole, cacertfile, cacert.Subject.CommonName, cacert.Issuer.CommonName, dstfile)

	} else {
		copyFile(cacertfile, dstfile)

		logm := fmt.Sprintf("copied %s ca cert file '%s' (%s<-%s) to '%s'",
			carole, cacertfile, cacert.Subject.CommonName, cacert.Issuer.CommonName, dstfile)

		fmt.Printf("%s\n", logm)

		if err := logCertCopy(filepath.Dir(dstfile), logm); err != nil {
			return err
		}
	}

	return nil
}

func CopyCertChain(certfile, cafile, rootcafile string, certs *Certs, mgmt ManagementSubsysDescriptor,
	alyt AnalyticsSubsysDescriptor, ptl PortalSubsysDescriptor, gwy GatewaySubsysDescriptor,
	commonCsrOutDir string, customCsrOutDir string, isOva bool) error {

	// check input
	if len(certfile) == 0 {
		return fmt.Errorf("certfile name is empty")
	}

	if len(cafile) == 0 {
		return fmt.Errorf("issuing ca certfile name is empty")
	}

	if len(rootcafile) == 0 {
		return fmt.Errorf("root ca certfile name is empty")
	}

	// build trust chain
	const noexpire = false
	if ischain, err := CertVerify(certfile, cafile, rootcafile, noexpire); err != nil {
		return err

	} else if ischain == false {
		// no trust chain, return...
		return fmt.Errorf("could not build a trust chain from certificates %s, %s, %s\n", certfile, cafile, rootcafile)
	}

	var cert *x509.Certificate
	var cacert *x509.Certificate
	var rootcert *x509.Certificate
	var err error

	// parse certfile
	if cert, err = ParseCertFile(certfile); err != nil {
		return err
	}

	// parse ca file
	if cacert, err = ParseCertFile(cafile); err != nil {
		return err
	}

	// parse root ca file
	if rootcert, err = ParseCertFile(rootcafile); err != nil {
		return err
	}

	// update cert specs
	updateCertSpecs(certs, mgmt, alyt, ptl, gwy, commonCsrOutDir, customCsrOutDir)

	var verifyErrors = make([]error, 100)

	// wildcard cert can match one or more hostnames
	matchcount := 0

	// public user facing certs
	for _, certSpec := range certs.PublicUserFacingEkuServerAuth {
		// copy cert file
		err = verifyCopyCertfile(certfile, cert, &certSpec)

		if err != nil {
			verifyErrors = append(verifyErrors, err)
		} else {
			matchcount++

			// concat issuing ca and root ca files into ca chain file
			dstfile := certSpec.CsrSubdir + string(os.PathSeparator) + certSpec.CaFile
			if err = concatChainFile(cafile, rootcafile, dstfile, cacert, rootcert); err != nil {
				return err
			}
		}
	}

	// mutual auth server auth
	for _, certSpec := range certs.MutualAuthEkuServerAuth {
		err := verifyCopyCertfile(certfile, cert, &certSpec)
		if err != nil {
			verifyErrors = append(verifyErrors, err)

		} else {
			matchcount++

			// concat issuing ca and root ca files into ca chain file
			dstfile := certSpec.CsrSubdir + string(os.PathSeparator) + certSpec.CaFile
			if err = concatChainFile(cafile, rootcafile, dstfile, cacert, rootcert); err != nil {
				return err
			}
		}
	}

	// common client certs
	for _, certSpec := range certs.CommonEkuClientAuth {
		err := verifyCopyCertfile(certfile, cert, &certSpec)
		if err != nil {
			verifyErrors = append(verifyErrors, err)

		} else {
			matchcount++

			// concat issuing ca and root ca files into ca chain file
			dstfile := certSpec.CsrSubdir + string(os.PathSeparator) + certSpec.CaFile
			if err = concatChainFile(cafile, rootcafile, dstfile, cacert, rootcert); err != nil {
				return err
			}
		}
	}

	// gateway subsystem for ova
	if isOva {
		// gateway director
		certSpec := CertSpec{}
		certSpec.Cn = gwy.GetApicGatewayServiceEndpoint()
		updateCertSpec(certs, gwy.GetGatewaySubsysName(), "gateway-director", &certSpec, DatapowerOutDir)

		err = verifyCopyCertfile(certfile, cert, &certSpec)
		if err != nil {
			verifyErrors = append(verifyErrors, err)

		} else {
			matchcount++

			// concat issuing-ca and root-ca files into ca chain file
			dstfile := certSpec.CsrSubdir + string(os.PathSeparator) + certSpec.CaFile
			if err = concatChainFile(cafile, rootcafile, dstfile, cacert, rootcert); err != nil {
				return err
			}

			// copy issuing-ca cert file
			dstfile = certSpec.CsrSubdir + string(os.PathSeparator) + dot2dash(certSpec.Cn) + ".issuing-ca.pem"
			if err := copyCaFile(cafile, dstfile, cacert, "issuing"); err != nil {
				return err
			}

			// copy root-ca cert file
			dstfile = certSpec.CsrSubdir + string(os.PathSeparator) + dot2dash(certSpec.Cn) + ".root-ca.pem"
			if err := copyCaFile(rootcafile, dstfile, rootcert, "root"); err != nil {
				return err
			}
		}

		//// api gateway
		//certSpec = CertSpec{}
		//certSpec.Cn = gwy.GetApiGatewayEndpoint()
		//updateCertSpec(certs, gwy.GetGatewaySubsysName(), "api-gateway", &certSpec, DatapowerOutDir)
		//
		//err = verifyCopyCertfile(certfile, cert, &certSpec)
		//if err != nil {
		//	verifyErrors = append(verifyErrors, err)
		//} else {
		//	matchcount++
		//}
	}

	// no hostname match for the cert, show errors
	if matchcount == 0 {
		for _, err := range verifyErrors {
			if err != nil {
				fmt.Printf("failed verify... %v\n", err)
			}
		}
	}

	return nil
}

func copyCert(certfile string, certs *Certs, mgmt ManagementSubsysDescriptor, alyt AnalyticsSubsysDescriptor,
	ptl PortalSubsysDescriptor, gwy GatewaySubsysDescriptor, commonCsrOutDir string, customCsrOutDir string, isOva bool) error {

	fmt.Printf("\nprocessing file... '%s'\n", certfile)

	// parse cert file
	// cert file could be a chain with the cert and ca
	certlist, blockTypes, err := ParseCertFile2(certfile)
	if err != nil {
		return err
		//log.Fatalf("%v\n", err)
	}

	// it is possible that certlist is empty eg for private key file
	if len(certlist) == 0 {
		fmt.Printf("%s\n", "skip...")
		return nil
	}

	if len(certlist) > 1 {
		return fmt.Errorf("file '%s' is a cert chain... %v, one cert expected", certfile, blockTypes)
	}

	cert := certlist[0]

	if cert.IsCA {
		return fmt.Errorf("file '%s' is ca cert..., Subject CN: %v, Issuer CN: %v",
			certfile, cert.Subject.CommonName, cert.Issuer.CommonName)
		//log.Fatalf("file '%s' is ca cert..., Subject CN: %v, Issuer CN: %v\n",
		//	certfile, cert.Subject.CommonName, cert.Issuer.CommonName)
	}

	var verifyErrors = make([]error, 100)

	// update cert specs
	updateCertSpecs(certs, mgmt, alyt, ptl, gwy, commonCsrOutDir, customCsrOutDir)

	// wildcard cert can match one or more hostnames
	matchcount := 0

	// public user facing certs
	for _, certSpec := range certs.PublicUserFacingEkuServerAuth {
		err := verifyCopyCertfile(certfile, cert, &certSpec)
		if err != nil {
			verifyErrors = append(verifyErrors, err)
		} else {
			matchcount++
		}
	}

	// public certs
	for _, certSpec := range certs.PublicEkuServerAuth {
		err := verifyCopyCertfile(certfile, cert, &certSpec)
		if err != nil {
			verifyErrors = append(verifyErrors, err)
		} else {
			matchcount++
		}
	}

	// mutual auth server auth
	for _, certSpec := range certs.MutualAuthEkuServerAuth {
		err := verifyCopyCertfile(certfile, cert, &certSpec)
		if err != nil {
			verifyErrors = append(verifyErrors, err)
		} else {
			matchcount++
		}
	}

	// common client certs
	for _, certSpec := range certs.CommonEkuClientAuth {
		err := verifyCopyCertfile(certfile, cert, &certSpec)
		if err != nil {
			verifyErrors = append(verifyErrors, err)
		} else {
			matchcount++
		}
	}

	// gateway subsystem for ova
	if isOva {
		// gateway director
		certSpec := CertSpec{}
		certSpec.Cn = gwy.GetApicGatewayServiceEndpoint()
		updateCertSpec(certs, gwy.GetGatewaySubsysName(), "gateway-director", &certSpec, DatapowerOutDir)

		err = verifyCopyCertfile(certfile, cert, &certSpec)
		if err != nil {
			verifyErrors = append(verifyErrors, err)
		} else {
			matchcount++
		}

		// api gateway
		certSpec = CertSpec{}
		certSpec.Cn = gwy.GetApiGatewayEndpoint()
		updateCertSpec(certs, gwy.GetGatewaySubsysName(), "api-gateway", &certSpec, DatapowerOutDir)

		err = verifyCopyCertfile(certfile, cert, &certSpec)
		if err != nil {
			verifyErrors = append(verifyErrors, err)
		} else {
			matchcount++
		}
	}

	// no hostname match for the cert, show errors
	if matchcount == 0 {
		for _, err := range verifyErrors {
			if err != nil {
				fmt.Printf("failed verify... %v\n", err)
			}
		}
	}

	return nil
}

func verifyCopyCertfile(certfile string, cert *x509.Certificate, certSpec *CertSpec) error {

	// cn is certified hostname
	err := VerifyHostName(certSpec.Cn, cert)

	if err == nil {
		fmt.Printf("\nhostname %s verifies...\n", certSpec.Cn)

		// copy cert...
		dstfile := certSpec.CsrSubdir + string(os.PathSeparator) + certSpec.CertFile

		// check if destination file exists
		exists, err := isFileExist(dstfile)
		if err != nil {
			return err
		}

		if exists {
			fmt.Printf("skip copying cert file '%s' (%s<-%s), destination file '%s' already exists...\n",
				certfile, cert.Subject.CommonName, cert.Issuer.CommonName, dstfile)

		} else {
			copyFile(certfile, dstfile)

			logm := fmt.Sprintf("copied cert file '%s' (%s<-%s) to '%s'", certfile, cert.Subject.CommonName, cert.Issuer.CommonName, dstfile)

			fmt.Printf("%s\n", logm)

			if err := logCertCopy(filepath.Dir(dstfile), logm); err != nil {
				return err
			}

		}

		return nil
	}

	//fmt.Printf("failed hostname verify %s..., %v\n", certSpec.Cn, err)
	return err
}

func ParseCertFile2(certfile string) ([]*x509.Certificate, []string, error) {

	blocks, err := ReadDecodeCertFile2(certfile)

	if err != nil {
		return nil, nil, err
	}

	certs, blockTypes, err := ParseCertBytes2(blocks)
	return certs, blockTypes, err
}

func ReadDecodeCertFile2(certfile string) ([]*pem.Block, error) {

	certbytes := readFileBytes(certfile)

	bytes := certbytes

	blocks := make([]*pem.Block, 0, 10)

	for {
		block, rest := pem.Decode(bytes)

		if block == nil {
			var b strings.Builder
			_, _ = fmt.Fprintf(&b, "failed to pem decode file %s", certfile)
			//fmt.Printf("%s\n", b.String())
			return nil, errors.New(b.String())
		}

		blocks = append(blocks, block)

		if rest != nil && len(rest) > 0 {
			bytes = rest

		} else {
			// done...
			return blocks, nil
		}
	}
}

func ParseCertBytes2(blocks []*pem.Block) ([]*x509.Certificate, []string, error) {

	certs := make([]*x509.Certificate, 0, 100)
	blockTypes := make([]string, 0, 100)

	for _, block := range blocks {
		//fmt.Printf("parsing block... %v\n", block.Type)

		if block.Type == "PRIVATE KEY" {
			// skip
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, nil, err
		}

		certs = append(certs, cert)
		blockTypes = append(blockTypes, block.Type)
	}

	return certs, blockTypes, nil
}

func ParseCertFile(certfile string) (*x509.Certificate, error) {

	bytes, err := ReadDecodeCertFile(certfile)

	if err != nil {
		return nil, err
	}

	certs, err := ParseCertBytes(bytes)
	return certs, err
}

func ReadDecodeCertFile(certfile string) ([]byte, error) {

	certbytes := readFileBytes(certfile)

	block, _ := pem.Decode(certbytes)
	if block == nil {
		var b strings.Builder
		_, _ = fmt.Fprintf(&b, "failed to pem decode file %s", certfile)
		//fmt.Printf("%s\n", b.String())
		return nil, errors.New(b.String())
	}

	return block.Bytes, nil
}

func ParseCertBytes(bytes []byte) (*x509.Certificate, error) {
	cert, err := x509.ParseCertificate(bytes)
	return cert, err
}

func VerifyHostName(hostname string, cert *x509.Certificate) error {
	err := cert.VerifyHostname(hostname)
	return err
}
