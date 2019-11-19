package main

import (
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/szesto/apicupcfg/apicupcfg"
	"log"
)

func main() {
	// template embedding box
	tbox := rice.MustFindBox("../../templates")

	// input: configuration file, output dir, csr subdirectories
	input, outdir, validateIp, initConfig, initConfigType, subsysOnly, certsOnly,
	certcopy, certdir, certverify, certfile, cafile, rootcafile, noexpire, certconcat, gen,
	soma, req, auth, url, setfile, dpdir, dpfile, datapowerOnly, dpdomain := apicupcfg.Input()

	// input actions
	isValidateIpActionf := func() bool {return validateIp}
	isCertCopyActionf := func() bool {return len(certcopy) > 0}
	isCertDirActionf := func() bool {return len(certdir) > 0}
	isCertVerifyActionf := func() bool {return certverify}
	isCertConcatActionf := func() bool {return certconcat}
	isGenActionf := func() bool {return gen}

	isSomaf := func() bool {return soma}
	isReqf := func() bool {return len(req) > 0}
	//isAuthf := func() bool {return len(auth) > 0}
	isUrlf := func() bool {return len(url) > 0}
	isSetfilef := func() bool {return len(setfile) > 0}
	isDpdirf := func() bool {return len(dpdir) > 0}
	isDpfilef := func() bool {return len(dpfile) > 0}

	// check input actions
	if !isValidateIpActionf() && !isCertCopyActionf() && !isCertDirActionf() &&
		!isCertVerifyActionf() && !isCertConcatActionf() && !isGenActionf() && !isSomaf() {

		log.Fatalf("no action specified... use apicupcfg -h for help...")
	}

	// output files
	output := apicupcfg.OutputFiles(outdir)

	// install type
	installType := apicupcfg.InstallTypeUknown

	if initConfig {
		installType = apicupcfg.InstallTypeInit

	} else {
		installType = apicupcfg.InstallType(input)
	}

	switch installType {
	case apicupcfg.InstallTypeOva:
		// load subsystem
		subsysvm := apicupcfg.LoadSubsysVm(input)

		if isValidateIpActionf() {
			apicupcfg.ValidateHostIpVm(subsysvm)

		} else if isSomaf() {
			// soma request

			if !isUrlf() {
				log.Fatal("-url command line opotion required with -soma flag")
			}

			if isReqf() && isSetfilef() {
				log.Fatal("-req and -setfile options are mutually exclusive with the -soma flag")
			}

			if !isReqf() && !isSetfilef() {
				log.Fatal("-req or -setfile options are required with the -soma flag")
			}

			if isSetfilef() && !isDpdirf() {
				log.Fatal("-dpdir option is required for -soma with -setfile")
			}

			if isSetfilef() && !isDpfilef() {
				log.Fatal("-dpfile option is required for -soma with -setfile")
			}

			if isReqf() {
				// soma request
				status, statusCode, reply, err := apicupcfg.SomaReq(req, auth, url, tbox)
				if err != nil {
					fmt.Printf("%v\n\n", err)

				} else {
					fmt.Printf("%s, %d, %s\n\n", status, statusCode, reply)
				}

			} else if isSetfilef() {
				// soma set-file request
				status, statusCode, reply, err := apicupcfg.SomaUpload(subsysvm, setfile, dpdir, dpfile, dpdomain, auth, url, tbox)
				if err != nil {
					fmt.Printf("%v\n\n", err)

				} else {
					fmt.Printf("%s, %d, %s\n\n", status, statusCode, reply)
				}
			}

		} else {
			// create output directories
			err := apicupcfg.CreateOutputDirectories(outdir, apicupcfg.CommonCsrOutDir,
				apicupcfg.CustomCsrOutDir, apicupcfg.SharedCsrOutDir, apicupcfg.ProjectOutDir, apicupcfg.DatapowerOutDir)

			if err != nil {
				log.Fatal(err)
			}

			if isCertCopyActionf() {
				// certcopy action

				if isCertCopyActionf() && isCertDirActionf() {
					log.Fatalf("%s\n", "-certcopy and -certdir options are mutually exclusive...")
				}

				// copy certs
				err = apicupcfg.CopyCertVm(certcopy, false, subsysvm, apicupcfg.CommonCsrOutDir, apicupcfg.CustomCsrOutDir)

				if err != nil {
					log.Fatal(err)
				}

			} else if isCertDirActionf() {
				// certdir action

				if isCertCopyActionf() && isCertDirActionf() {
					log.Fatalf("%s\n", "-certcopy and -certdir options are mutually exclusive...")
				}

				// copy certs in cert dir
				err = apicupcfg.CopyCertVm(certdir, true, subsysvm, apicupcfg.CommonCsrOutDir, apicupcfg.CustomCsrOutDir)

				if err != nil {
					log.Fatal(err)
				}

			} else if isCertVerifyActionf() {
				// certverify action
				isvalid, err := apicupcfg.CertVerify(certfile, cafile, rootcafile, noexpire)

				if err != nil {
					log.Fatal(err)

				} else if isvalid {
					cfile := certfile; if len(certfile) == 0 {cfile = cafile}
					fmt.Printf("Certificate file '%s' verifies...\n", cfile)

				} else {
					cfile := certfile; if len(certfile) == 0 {cfile = cafile}
					fmt.Printf("Certificate file '%s' does not verify...\n", cfile)
				}

			} else if isCertConcatActionf() {
				// cert concat action

				err = apicupcfg.CertConcat(cafile, rootcafile, subsysvm.Certs.CaFile, outdir,
					apicupcfg.CommonCsrOutDir, apicupcfg.CustomCsrOutDir, input)

				if err != nil {
					log.Fatal(err)

				}

			} else if isGenActionf() {
				// gen action

				// apply templates
				apicupcfg.ApplyTemplateVm(subsysvm, output, subsysOnly, certsOnly, datapowerOnly, tbox)
			}
		}

	case apicupcfg.InstallTypeK8s:
		// load subsystem
		subsysk8s := apicupcfg.LoadSubsysK8s(input)

		if isValidateIpActionf() {
			// not applicable, complain
			fmt.Printf("validateip command line option is not applicable to the %s install type...\n", apicupcfg.InstallTypeK8s)

		} else if isSomaf() {
			// not applicable, complain
			fmt.Printf("soma command line option is not applicable to the %s install type...\n", apicupcfg.InstallTypeK8s)

		} else {
			// create output directories
			datapowerdir := ""
			err := apicupcfg.CreateOutputDirectories(outdir, apicupcfg.CommonCsrOutDir,
				apicupcfg.CustomCsrOutDir, apicupcfg.SharedCsrOutDir, apicupcfg.ProjectOutDir, datapowerdir)

			if err != nil {
				log.Fatal(err)
			}

			if isCertCopyActionf() {

				if isCertCopyActionf() && isCertDirActionf() {
					log.Fatalf("%s\n", "-certcopy and -certdir options are mutually exclusive...")
				}

				// copy certs
				err = apicupcfg.CopyCertK8s(certcopy, false, subsysk8s, apicupcfg.CommonCsrOutDir, apicupcfg.CustomCsrOutDir)

				if err != nil {
					log.Fatal(err)
				}

			} else if isCertDirActionf() {

				if isCertCopyActionf() && isCertDirActionf() {
					log.Fatalf("%s\n", "-certcopy and -certdir options are mutually exclusive...")
				}

				// copy certs in cert dir
				err = apicupcfg.CopyCertK8s(certdir, true, subsysk8s, apicupcfg.CommonCsrOutDir, apicupcfg.CustomCsrOutDir)

				if err != nil {
					log.Fatal(err)
				}

			} else if isCertVerifyActionf() {
				// certverify action
				isvalid, err := apicupcfg.CertVerify(certfile, cafile, rootcafile, noexpire)

				if err != nil {
					log.Fatal(err)

				} else if isvalid {
					cfile := certfile; if len(certfile) == 0 {cfile = cafile}
					fmt.Printf("Certificate file '%s' verifies...\n", cfile)

				} else {
					cfile := certfile; if len(certfile) == 0 {cfile = cafile}
					fmt.Printf("Certificate file '%s' does not verify...\n", cfile)
				}

			} else if isCertConcatActionf() {
				// cert concat action

				err = apicupcfg.CertConcat(cafile, rootcafile, subsysk8s.Certs.CaFile, outdir,
					apicupcfg.CommonCsrOutDir, apicupcfg.CustomCsrOutDir, input)

				if err != nil {
					log.Fatal(err)

				}

			} else if isGenActionf() {
				// apply templates
				apicupcfg.ApplyTemplatesK8s(subsysk8s, output, subsysOnly, certsOnly, tbox)
			}
		}

	case apicupcfg.InstallTypeInit:
		apicupcfg.InitConfig(input, initConfigType, tbox)

	default:
		log.Fatalf("unsupported install type %s\n", installType)
	}
}
