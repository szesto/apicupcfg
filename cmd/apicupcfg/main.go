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
	soma, req, auth, url, setfile, dpdir, dpfile, datapowerOnly, dpdomain, dpcacopy, certchaincopy,
	trustdir, version := apicupcfg.Input()

	func (certdir string, certconcat bool, dpcacopy bool, certchaincopy bool) {} (certdir, certconcat, dpcacopy, certchaincopy)

	// input actions
	isValidateIpActionf := func() bool {return validateIp}
	isCertCopyActionf := func() bool {return certcopy && len(certdir) == 0 }
	isCertDirActionf := func() bool {return certcopy && len(certdir) > 0}
	isCertVerifyActionf := func() bool {return certverify}
	isGenActionf := func() bool {return gen}

	isSomaf := func() bool {return soma}
	isReqf := func() bool {return len(req) > 0}
	isUrlf := func() bool {return len(url) > 0}
	isSetfilef := func() bool {return len(setfile) > 0}
	isDpdirf := func() bool {return len(dpdir) > 0}
	isDpfilef := func() bool {return len(dpfile) > 0}

	isInitConfigActionf := func() bool {return initConfig}

	isVersionf := func() bool { return version }

	if isVersionf() {
		fmt.Printf("%s\n", showVersion())
		return
	}

	// check input actions
	if !isValidateIpActionf() && !isCertCopyActionf() &&
		!isCertVerifyActionf() && !isGenActionf() &&
		!isSomaf() && !isInitConfigActionf() && !isCertDirActionf() && !isVersionf() {

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
				status, statusCode, reply, err := apicupcfg.SomaUpload(setfile, dpdir, dpfile, dpdomain, auth, url, tbox)
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

			if isCertVerifyActionf() {
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

			} else if isCertCopyActionf() {

				err = apicupcfg.CopyCertChain(certfile, cafile, rootcafile, &subsysvm.Certs,
					&subsysvm.Management, &subsysvm.Analytics, &subsysvm.Portal, &subsysvm.Gateway,
					apicupcfg.CommonCsrOutDir, apicupcfg.CustomCsrOutDir, true)

				if err != nil {
					log.Fatal(err)

				}

			} else if isCertDirActionf() {

				err = apicupcfg.CopyCertDir(certdir, trustdir, &subsysvm.Certs,
					&subsysvm.Management, &subsysvm.Analytics, &subsysvm.Portal, &subsysvm.Gateway,
					apicupcfg.CommonCsrOutDir, apicupcfg.CustomCsrOutDir, true)

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
			datapowerdir := apicupcfg.DatapowerOutDir
			err := apicupcfg.CreateOutputDirectories(outdir, apicupcfg.CommonCsrOutDir,
				apicupcfg.CustomCsrOutDir, apicupcfg.SharedCsrOutDir, apicupcfg.ProjectOutDir, datapowerdir)

			if err != nil {
				log.Fatal(err)
			}

			if isCertVerifyActionf() {
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

			} else if isCertCopyActionf() {

				err = apicupcfg.CopyCertChain(certfile, cafile, rootcafile, &subsysk8s.Certs,
					&subsysk8s.Management, &subsysk8s.Analytics, &subsysk8s.Portal, &subsysk8s.Gateway,
					apicupcfg.CommonCsrOutDir, apicupcfg.CustomCsrOutDir, false)

				if err != nil {
					log.Fatal(err)

				}

			} else if isCertDirActionf() {

				err = apicupcfg.CopyCertDir(certdir, trustdir, &subsysk8s.Certs,
					&subsysk8s.Management, &subsysk8s.Analytics, &subsysk8s.Portal, &subsysk8s.Gateway,
					apicupcfg.CommonCsrOutDir, apicupcfg.CustomCsrOutDir, false)

				if err != nil {
					log.Fatal(err)
				}

			} else if isGenActionf() {
				// apply templates
				apicupcfg.ApplyTemplatesK8s(subsysk8s, output, subsysOnly, certsOnly, datapowerOnly, tbox)
			}
		}

	case apicupcfg.InstallTypeInit:
		apicupcfg.InitConfig(input, initConfigType, tbox)

	default:
		log.Fatalf("unsupported install type %s\n", installType)
	}
}
