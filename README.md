** Running apicupcfg tool **

apcupcfg is a tool to generate configuration scripts for the IBM API Connect 2018.4.1.x installation.

Both Kubernetes and OVA installation types are supported.

apicupcfg tool generates all subsystem and certificate scripts, all supportring certificate configuration
scripts, and provides a number of validation options.

To see all available command line options run:

`apicupcfg -h` on linux or macos
`apicupcfg.exe -h` on windows

`apicupcfg` collects all required configuration data in one input json file.

The directory with the configuration file is working directory.

To create input configuration file, change to the working directory and generate json 
configuration file template:

`apicupcfg -gen -initconfig -configtype ova|k8s [-config subsys-config.json]`

`subsys-config.json` is the default configuration file name and is assumed when the `-config` option 
is ommitted.

Ova and kuberenetes configration files are different.

There are a number of sections in the input configuration file. There is a section for each subsystem,
certificate configuration section, and defaults section that applies to all subystems.

To generate configuration scripts and directories:
`apicupcfg -gen -out output-directory [-config subsys-config.json]`

Default output subdirectory is `output` and is assumed if `-out` command line option is ommitted.
Default input configuration file is subsys-config.json

To generate configuration scripts and directories in the working directory:
`apicupcfg -gen -out . [-config config-subsys.json]`

Generated directories:

bin - directory for apicup executable. Copy your version of apicup executable to the bin directory and then copy it to the apicup executable.
project - apicup install configuration.
custom-csr - openssl scripts for custom certificates.
common-csr - openssl scripts for common certificates.
shared-dir - openssl scripts for shared endpoint trust.

Generated scripts and configuration data depend on the input configuration file.
General rule is that all required scripts and configuration files are created.

Each generated configuration file is tagged with the value defined in the input json configuration file.

Generated subsystem configuration scripts:

apicup-subsys-set-management.tag.(sh|bat)
apicup-subsys-set-analytics.tag.(sh|bat)
apicup-subsys-set-portal.tag.(sh|bat)
apicup-subsys-set-gateway.tag.(sh|bat)

Generated certificate setting scripts.

apicup-certs-set-user-facing-public.tag.(sh|bat)
apicup-certs-set-user-facing-public.tag.(sh|bat)
apicup-certs-set-mutual-auth.tag.(sh|bat)
apicup-certs-set-common.tag.(sh|bat)

apicup-certs-set-certbot-user-facing-public.tag.(sh|bat)
apicup-certs-set-certbot-public.tag.(sh|bat)

apicup-certs-set-shared-trust-user-facing-public.tag.(sh|bat)
apicup-certs-set-shared-trust-public.tag.(sh|bat)
apicup-certs-set-shared-trust-mutual-auth.tag.(sh|bat)

**Running generated scripts.**

Generated scripts (other then openssl scripts in csr subdirectories) must be run from the project subdirectory:

`cd project`
`../apicup-subsys-set-management.tag.sh` (linux, macos)
`..\apicup-subsys-set-management.tag.bat` (windows)

**Working with certificates**

Certificate generation is driven by the settings in the input configuration file.

`Certs: {
    PublicUserFacingCerts: true|false,
    PublicCerts: true|false,
    CommonCerts: true|false
}
`
It is recommended to generate public user facing certs only. Other types of certificates are advanced use cases.

**Openssl scripts**.

custom-csr subdirectory contains openssl scripts for public user facing certs and public certs.
common-csr subdirectory contains openssl scripts for common certificates, this includes mutual auth and client certs.
shared-csr subdirectory contains openssl scripts for shared endpoint trust.

Each endpoint the json configuration file is transformed into the cn and openssl csr configuration
file and script are generated in the csr subdirectory.

The Certs.DnFields value defines components of the dn that are required by the certificate authority
and included in the csr configuration file.

To simplify running of these scripts, the custom-csr subdirectory contains 
`all-user-facing-public-csr.tag.(sh|bat)` and `all-public-csr.tag.(sh|bat)` scripts.

The common-csr subdirectory contains `all-mutual-auth-csr.tag.(sh|bat)` and `all-common-csr.tag.sh` scripts.

Run these scripts to generate key pairs and csr's. Submit csr's to the certifiacte authority to get signed certificates.
Copy received certificates to the correct destination.

** Copying certificates to the correct files **.

Certificate settings scripts expect to find certificates, private keys and root certificates at specific locations.
Certificates recieved from the ca can be manually copied to the correct destination but this is error prone.

To copy a certificate received from the ca to the correct destination:
`apicupcfg -certcopy path-to-certificate-file.pem [-config subsys-config.json]`

This command will introspect the certificate match it with endpoints defined in the configuration file
and copy certificate to the correct destination. Note that if certificate matches mulitple endpoints
(wildcard or shared trust) then a separate copy will be made for each endpoint.

To process all certificates received from the ca together, place them in a directory and run:
`apicupcfg -certdir path-to-a-directory-with-certificates [-config subsys-config.json]`

This command will copy all certificates in the directory to the correct destination.

** Copying ca trust file **

You must create a file that concatenates intermidiate ca certificate and root ca certificate and copy it to the correct destination.
This could be done manually but it is error prone.

To concatenate intermediate ca cert and root ca cert and copy this file to correct destination:
`apicupcfg -certconcat -ca path-to-ca.pem -rootca path-to-root-ca.pem`

This command will verify certificates and if valid copy combined file to a destination specified
in the Certs.CaFile value. Combined file will be copied to the custom-csr and common-csr subdirectories.

** Verifying certificates **/

To verify certificate:
`apicupcfg -certverify [-noexpire] -cert path-to-cert.pem -ca path-to-intermediate-ca.pem -rootca path-to-root-ca.pem`

To verify intermidiate ca certificate:
`apicupcfg -certverify [-noexpire] -ca path-to-intermediate-ca.pem -rootca path-to-root-ca.pem`

This command will compute and display trust chain. Pass `-noexpire` to ignore certificate expiration.

**Buid**

This code was developed with go version 1.13.

Install rice dependencies (see *Template Embedding* section)

`go get github.com/Masterminds/sprig`
`go get github.com/GeertJohan/go.rice`
`go get github.com/GeertJohan/go.rice/rice`

Change to the cmd/apicupcfg: then `go install`.

*apicupcfg* executable will be in the $GOPATH/bin (or %GOPATH%\bin) directory:

`$GOPATH/bin/apicupcfg -help`
`%GOPATH%\bin\apicupcfg.exe -help`

Resulting executable is operating-system specific. 
File path syntax and command file syntax are native to the target operating system.

**Tempate embedding.**

rice is a tool for embedding go templates into the executable.

install rice package:
`go get github.com/GeertJohan/go.rice`
`go get github.com/GeertJohan/go.rice/rice`

if you *udpate* templates, generate new rice-box.go:
in the cmd/apicupcfg directory:
`rice clean`
`rice embed-go`

**General Info.**

*No IBM API Connect software is required to build or use this tool.*

**Command line reference.**

`apicupcfg -help`
`apicupcfg -gen -initconfig -configtype ova|k8s [-config subsys-config.json]`
`apicupcfg -gen -out . [-config subsys-config]`
`apicupcfg -certcopy cerftile.pem [-config subsys-config.json]`
`apicupcfg -certdir certdir [-config subsys-config.json]`
`apicupcfg -certvalidate [-noexpire] [-cert cert.pem] -ca ca.pem -rootca rootca.pem`
`apicupcfg -certconcat -ca ca.pem -rootca rootca.pem`
`apicupcfg -validateip [-config subsys-config.json]`

@todo

**Configuraton reference.**

@todo

