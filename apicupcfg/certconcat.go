package apicupcfg

import (
	"fmt"
	"os"
)

func CertConcat(cacertfile string, rootcafile string, outfile string, outdir, commoncsrdir, customcsrdir string) error {

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
