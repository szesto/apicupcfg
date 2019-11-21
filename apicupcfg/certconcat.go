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

func CaCopy(cafile, rootcafile, dstcafile, dstrootcafile, outdir, dstdir, configfile string) error {

	// check input parameters
	if len(cafile) == 0 {
		return fmt.Errorf("%s","ca-copy... ca file parameter is emtpy... value required")
	}

	if len(rootcafile) == 0 {
		return fmt.Errorf("%s","ca-copy... root ca file parameter is emtpy... value required")
	}

	if len(dstcafile) == 0 {
		return fmt.Errorf("ca-copy... destination ca file parameter is emtpy... check the setting '%s' in the config file '%s'",
			"Gateway.CaFile", configfile)
	}

	if len(dstrootcafile) == 0 {
		return fmt.Errorf("ca-copy... destination root ca file parameter is emtpy... check the setting '%s' in the config file '%s'",
			"Gateway.RootCaFile", configfile)
	}

	if len(outdir) == 0 {
		return fmt.Errorf("%s", "ca-copy... base output directory parameter is empty... value required")
	}

	if len(dstdir) == 0 {
		return fmt.Errorf("%s", "ca-copy... destination directory parameter is empty... value required")
	}

	// verify ca file
	isvalid, err := CertVerify("", cafile, rootcafile, true)
	if err != nil {
		return err

	} else if ! isvalid {
		return fmt.Errorf("ca cert file %s could not be verified", cafile)
	}

	// copy ca file
	dstfile := outdir + string(os.PathSeparator) + dstdir + string(os.PathSeparator) + dstcafile
	fmt.Printf("ca-copy... copying ca file '%s' to destination '%s'\n", cafile, dstfile)

	exist, err := copyFileErrExist(cafile, dstfile)
	if err != nil {
		return err
	}

	if exist {
		fmt.Printf("ca-file... destination file '%s' already exists... skip copy\n", dstfile)
	}

	// copy root-ca file
	dstfile = outdir + string(os.PathSeparator) + dstdir + string(os.PathSeparator) + dstrootcafile
	fmt.Printf("ca-copy... copying root ca file '%s' to destination '%s'\n", rootcafile, dstfile)

	exist, err = copyFileErrExist(rootcafile, dstfile)
	if err != nil {
		return err
	}

	if exist {
		fmt.Printf("ca-file... destination file '%s' already exists... skip copy\n", dstfile)
	}

	return nil
}

func copyFileErrExist(srcfile, dstfile string) (bool, error) {
	exist, err := isFileExist(dstfile)

	if err != nil {
		return false, err
	}

	copyFile(srcfile, dstfile)
	return exist, nil
}
