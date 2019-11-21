package apicupcfg

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

	var verifyErrors []error = make([]error, 100)

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
		fmt.Printf("hostname %s verifies...\n", certSpec.Cn)

		// copy cert...
		dstfile := certSpec.CsrSubdir + string(os.PathSeparator) + certSpec.CertFile

		// check if destination file exists
		exists, err := isFileExist(dstfile)
		if err != nil {
			return err
		}

		if exists {
			fmt.Printf("cert file '%s' destination '%s' already exists... skip...\n", certfile, dstfile)

		} else {
			copyFile(certfile, dstfile)
			fmt.Printf("cert file '%s' copied to destination %s\n", certfile, dstfile)

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
	err := cert.VerifyHostname(hostname);
	return err
}
