package apicupcfg

import (
	"fmt"
	"os"
)

func CertConcat(cacertfile string, rootcafile string, outfile string, outdir,
	commoncsrdir, customcsrdir string, configfile string) error {

	// make sure output file is defined
	if len(outfile) == 0 {
		return fmt.Errorf("cert-concat... output ca file name is empty, value required... check Certs.CaFile parameter in the config file '%s'\n",
			configfile)
	}

	// verify ca cert file
	isvalid, err := CertVerify("", cacertfile, rootcafile, true)
	if err != nil {
		return err
	}

	if ! isvalid {
		return fmt.Errorf("ca cert file %s could not be verified", cacertfile)
	}

	// concatenate
	dstfile := outdir + string(os.PathSeparator) + commoncsrdir + string(os.PathSeparator) + outfile

	fmt.Printf("concat ca cert '%s' and root cert '%s' into '%s'\n", cacertfile, rootcafile, dstfile)
	concatFiles(cacertfile, rootcafile, dstfile)

	dstfile = outdir + string(os.PathSeparator) + customcsrdir + string(os.PathSeparator) + outfile

	fmt.Printf("concat ca cert '%s' and root cert '%s' into '%s'\n", cacertfile, rootcafile, dstfile)
	concatFiles(cacertfile, rootcafile, dstfile)

	return nil
}
