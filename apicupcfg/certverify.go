package apicupcfg

import (
	"crypto/x509"
	"fmt"
	"time"
)

func CertVerify2(certfile string, chainfiles []string, noexpire bool) (bool, error) {

	var cert *x509.Certificate
	var err error

	if len(certfile) > 0 {
		//fmt.Printf("parsing file '%s'...\n", certfile)

		if cert, err = ParseCertFile(certfile); err != nil {
			return false, err
		}

		// must not be ca cert
		if cert.IsCA {
			err := fmt.Errorf("file '%s' is ca certificate... subject cn %v, issuer cn %v",
				certfile, cert.Subject.CommonName, cert.Issuer.CommonName)
			return false, err
		}
	}

	chaincerts := make([]*x509.Certificate, len(chainfiles))
	for chainidx, chainfile := range chainfiles {

		//fmt.Printf("parsing file '%s'...\n", chainfile)
		if chaincerts[chainidx], err = ParseCertFile(chainfile); err != nil {
			return false, err
		}

		cacert := chaincerts[chainidx]

		// must be ca cert
		if ! cacert.IsCA {
			err := fmt.Errorf("file '%s' is not a ca certificate... subject cn %v, issuer cn %v",
				chainfile, cacert.Subject.CommonName, cacert.Issuer.CommonName)
			return false, err
		}
	}

	chains, err := BuildChain2(cert, chaincerts, noexpire)
	//fmt.Printf("%v, %v\n", chains, err)

	if err != nil {
		return false, err
	}

	//fmt.Printf("not-before: %v, not after %v\n", cert.NotBefore, cert.NotAfter)

	fmt.Printf("\nbuilding trust chain...\n")

	for chi, chain := range chains {

		for ci, cp := range chain {
			isroot := cp.Subject.CommonName == cp.Issuer.CommonName

			fmt.Printf("cert[%d,%d]... subj-cn... %v, is-ca... %t, issuer-cn... %v, is-root... %t\n",
				chi, ci, cp.Subject.CommonName, cp.IsCA, cp.Issuer.CommonName, isroot)

			if ci == len(chain)-1 && isroot {
				// chain is terminated with the root cert
				return true, nil
			}
		}
	}

	// no chains all the way to the root
	return false, nil
}

func CertVerify(certfile string, cafile string, rootcafile string, noexpire bool) (bool, error) {

	var cert *x509.Certificate
	var cacert *x509.Certificate
	var rootcert *x509.Certificate
	var err error

	if len(certfile) > 0 {
		//fmt.Printf("parsing file '%s'...\n", certfile)

		if cert, err = ParseCertFile(certfile); err != nil {
			return false, err
		}

		// must not be ca cert
		if cert.IsCA {
			err := fmt.Errorf("file '%s' is ca certificate... subject cn %v, issuer cn %v",
				certfile, cert.Subject.CommonName, cert.Issuer.CommonName)
			return false, err
		}
	}

	if len(cafile) > 0 {
		//fmt.Printf("parsing file '%s'...\n", chainfile)
		if cacert, err = ParseCertFile(cafile); err != nil {
			return false, err
		}

		// must be ca cert
		if ! cacert.IsCA {
			err := fmt.Errorf("file '%s' is not a ca certificate... subject cn %v, issuer cn %v",
				cafile, cacert.Subject.CommonName, cacert.Issuer.CommonName)
			return false, err
		}
	}

	if len(rootcafile) > 0 {
		//fmt.Printf("parsing file '%s'\n...", rootcafile)
		if rootcert, err = ParseCertFile(rootcafile); err != nil {
			return false, err
		}

		// must be root cert
		isroot := rootcert.Subject.CommonName == rootcert.Issuer.CommonName

		//fmt.Printf("rootcert... subj-cn... %v, is-ca... %t, issuer-cn... %v, is-root... %t\n",
		//	rootcert.Subject.CommonName, rootcert.IsCA, rootcert.Issuer.CommonName, isroot)

		if ! isroot {
			err := fmt.Errorf("file '%s' is not a root cert... subject cn %v, issuer cn %v",
				rootcafile, rootcert.Subject.CommonName, rootcert.Issuer.CommonName)
			return false, err
		}
	}

	chains, err := BuildChain(cert, cacert, rootcert, noexpire)
	//fmt.Printf("%v, %v\n", chains, err)

	if err != nil {
		return false, err
	}

	//fmt.Printf("not-before: %v, not after %v\n", cert.NotBefore, cert.NotAfter)

	fmt.Printf("\nbuilding trust chain...\n")

	for chi, chain := range chains {

		for ci, cp := range chain {
			isroot := cp.Subject.CommonName == cp.Issuer.CommonName

			fmt.Printf("cert[%d,%d]... subj-cn... %v, is-ca... %t, issuer-cn... %v, is-root... %t\n",
				chi, ci, cp.Subject.CommonName, cp.IsCA, cp.Issuer.CommonName, isroot)

			if ci == len(chain)-1 && isroot {
				// chain is terminated with the root cert
				return true, nil
			}
		}
	}

	// no chains all the way to the root
	return false, nil
}

func BuildChain2(cert *x509.Certificate, fullchain []*x509.Certificate, noexpire bool) ([][]*x509.Certificate, error) {

	chainpool := x509.NewCertPool()
	rootpool := x509.NewCertPool()

	if cert != nil {
		for _, c := range fullchain {
			if c.IsCA && c.Issuer.CommonName == c.Subject.CommonName {
				rootpool.AddCert(c)

			} else if c.IsCA {
				chainpool.AddCert(c)
			}
		}

		opts := x509.VerifyOptions{Roots:rootpool, Intermediates:chainpool,}

		if noexpire {
			opts.CurrentTime = cert.NotBefore.Add(time.Duration(time.Hour * 24))
		}

		fmt.Printf("verifying cert... subject-cn %v\n", cert.Subject.CommonName)
		chains, err := cert.Verify(opts)
		return chains, err
	}

	return nil, fmt.Errorf("%s", "build-chain... no certificates to verfiy...")
}

func BuildChain(cert *x509.Certificate, chaincert *x509.Certificate, rootcert *x509.Certificate, noexpire bool) ([][]*x509.Certificate, error) {

	chainpool := x509.NewCertPool()
	rootpool := x509.NewCertPool()

	if cert != nil {
		if chaincert != nil {
			chainpool.AddCert(chaincert)
		}

		if rootcert != nil {
			rootpool.AddCert(rootcert)
		}

		opts := x509.VerifyOptions{Roots:rootpool, Intermediates:chainpool,}

		if noexpire {
			opts.CurrentTime = cert.NotBefore.Add(time.Duration(time.Hour * 24))
		}

		fmt.Printf("verifying cert... subject-cn %v\n", cert.Subject.CommonName)
		chains, err := cert.Verify(opts)
		return chains, err

	} else if chaincert != nil {
		if rootcert != nil {
			rootpool.AddCert(rootcert)
		}

		opts := x509.VerifyOptions{Roots:rootpool,}

		if noexpire {
			opts.CurrentTime = chaincert.NotBefore.Add(time.Duration(time.Hour * 24))
		}

		fmt.Printf("verifying cert... subject-cn %v\n", chaincert.Subject.CommonName)
		chains, err := chaincert.Verify(opts)
		return chains, err
	}

	return nil, fmt.Errorf("%s", "build-chain... no certificates to verfiy...")
}
