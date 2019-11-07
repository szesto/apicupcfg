package apicupcfg

import (
	"fmt"
	"testing"
)

func TestCertCopy1(t *testing.T) {

	certfile := "testcert.pem"
	cert, err := ParseCertFile(certfile)

	if err == nil {

		hostname := "cm.apim.cloud"
		err = VerifyHostName(hostname, cert)

		if err == nil {
			fmt.Printf("hostname %s verifies...\n", hostname)

		} else {
			fmt.Printf("verify hostname %s..., %v\n", hostname, err)
		}

		hostname = "apim.apim.com"
		err = VerifyHostName(hostname, cert)

		if err == nil {
			fmt.Printf("hostname %s verifies...\n", hostname)

		} else {
			fmt.Printf("failed verify..., %v\n", err)
		}
	}
}

func TestCertVerify1(t *testing.T) {

	certfile := "/Users/simon/local/aws/certbot/letsencrypt/live/apim.cloud/cert.pem"
	chain := "/Users/simon/local/aws/certbot/letsencrypt/live/apim.cloud/chain.pem"
	//rootca := "/Users/simon/local/aws/certbot/letsencrypt/live/apim.cloud/fullchain.pem"
	rootca := "/Users/simon/local/aws/certbot/letsencrypt/live/apim.cloud/rootca.pem"

	err := CertVerify(certfile, chain, rootca)

	fmt.Printf("%v\n", err)
}