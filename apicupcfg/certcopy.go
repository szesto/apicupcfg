package apicupcfg

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func copyCert(certfile string, certs *Certs, mgmt ManagementSubsysDescriptor, alyt AnalyticsSubsysDescriptor,
	ptl PortalSubsysDescriptor, gwy GatewaySubsysDescriptor, commonCsrOutDir string, customCsrOutDir string) {

	// parse cert file
	cert, err := ParseCertFile(certfile)
	if err != nil {
		log.Fatalf("%v\n", err)
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

	// hostname did not verify for any cert
	if matchcount == 0 {
		for _, err := range verifyErrors {
			if err != nil {
				fmt.Printf("failed verify... %v\n", err)
			}
		}
	}

	return
}

func verifyCopyCertfile(certfile string, cert *x509.Certificate, certSpec *CertSpec) error {

	// cn is certified hostname
	err := VerifyHostName(certSpec.Cn, cert)

	if err == nil {
		fmt.Printf("hostname %s verifies...\n", certSpec.Cn)

		// copy cert...
		dstfile := certSpec.CsrSubdir + string(os.PathSeparator) + certSpec.CertFile
		copyFile(certfile, dstfile)

		fmt.Printf("cert file %s copied to destination %s\n", certfile, dstfile)
		return nil
	}

	//fmt.Printf("failed hostname verify %s..., %v\n", certSpec.Cn, err)
	return err

}

func ParseCertFile(certfile string) (*x509.Certificate, error) {
	bytes, err := ReadDecodeCertFile(certfile)
	if err == nil {
		return ParseCertBytes(bytes)
	}
	return nil, err
}

func ReadDecodeCertFile(certfile string) ([]byte, error) {

	certbytes := readFileBytes(certfile)

	block, _ := pem.Decode(certbytes)
	if block == nil {
		var b strings.Builder
		_, _ = fmt.Fprintf(&b, "failed to pem decode cert file %s", certfile)
		fmt.Printf("%s\n", b.String())
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
