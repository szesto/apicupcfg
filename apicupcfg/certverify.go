package apicupcfg

import (
	"crypto/x509"
	"fmt"
	"time"
)

func CertVerify(certfile string, chainfile string, rootcafle string) error {

	cert, err := ParseCertFile(certfile)
	if err != nil {
		return err
	}

	fmt.Printf("not-before: %v, not after %v\n", cert.NotBefore, cert.NotAfter)

	chaincert, err := ParseCertFile(chainfile)
	if err != nil {
		return err
	}

	rootcert, err := ParseCertFile(rootcafle)
	if err != nil {
		return err
	}

	isroot := rootcert.Subject.CommonName == rootcert.Issuer.CommonName

	fmt.Printf("rootcert... subj-cn... %v, is-ca... %t, issuer-cn... %v, is-root... %t\n",
		rootcert.Subject.CommonName, rootcert.IsCA, rootcert.Issuer.CommonName, isroot)

	chainpool := x509.NewCertPool()
	chainpool.AddCert(chaincert)

	rootpool := x509.NewCertPool()
	rootpool.AddCert(rootcert)

	opts := x509.VerifyOptions{Roots:rootpool, DNSName:"cm.apim.cloud",
		CurrentTime: cert.NotBefore.Add(time.Duration(time.Hour * 24)),
		Intermediates:chainpool}

	var chains [][]*x509.Certificate
	chains, err = cert.Verify(opts)
	fmt.Printf("%v, %v\n", chains, err)

	if err == nil {
		for chi, chain := range chains {
			for ci, cp := range chain {
				isroot := cp.Subject.CommonName == cp.Issuer.CommonName

				fmt.Printf("cert[%d,%d]... subj-cn... %v, is-ca... %t, issuer-cn... %v, is-root... %t\n",
					chi, ci, cp.Subject.CommonName, cp.IsCA, cp.Issuer.CommonName, isroot)
			}
		}
	}

	return err
}
