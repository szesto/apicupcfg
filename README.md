**Running apicupcfg tool**

`apcupcfg` is a tool to generate configuration scripts for the IBM API Connect 2018.4.1.x installation.

Both Kubernetes and OVA installation types are supported.

**Typical Steps for the OVA install.**
- Create working directory.
- Generate json configuration file *subsys-config.json*:
    - `apicupcfg -initconfig [-configtype ova]`
- Edit *subsys-config.json* configuration file. In what follows, *tag* is the value of the *Tag* property.
- Validate subsystem (and datapower) ip addresses:
    - `apicupcfg -validateip`
- Generate subsystem and certificate setting scripts:
    - `apicupcfg -gen`
- Generate CSR's. 
    - From the *custom-csr* directory:
        - `all-user-facing-public-csr.tag.bat|sh`
- Submit generated csr files to the ca.
- Place certificates recieved from the ca in a directory. You can optionally split certificates in one directory (received-certs)
and ca certificates in another directory (trust-root). 
- Copy certificates to correct destinations:
    - `apicupcfg -certcopy -certdir dir [-trustdir dir] `
- Run subsystem and certificate setting scripts. 
    - From the *project* directory: 
        - `..\apicup-subsys-set-management.tag.bat|sh`
        - `..\apicup-subsys-set-analytics.tag.bat|sh`
        - `..\apicup-subsys-set-portal.tag.bat|sh`
        - `..\apicup-certs-set-user-facing-public.tag.bat|sh`
- Install subsystems with the `apicup subsys install` command as usual. 
    - From the *project* directory:
        - `..\bin\apicup subsys install mgmt --out mgmt-plan-out`
        - `..\bin\apicup subsys install alyt --out alyt-plan-out`
        - `..\bin\apicup subsys install ptl --out ptl-plan-out`
- Configure datapower cluster.

**Steps for datapower configuration**
- change to the *datapower* directory.
    - run `all-datapower-csr.tag.bat|sh`
    - submit 1 csr to the ca.
    - for each datapower, run initial configuration, set timezone, and enable xml management interface. Apply fixpack.
    - create *dp.env* file with the datapower admin creds: 1st line is username, 2nd line is password
    - run `*zoma-crypto-self-(datapower-name).bat|sh*` file for each datpower instance.
    - run `*zoma-(datapower-name.bat|sh)*` file for each datapower instance.
    - complete datapower crypto update 
    - run `*zoma-web-gui-(datapower-name).bat|sh*` file for each datpower instance.

**Datapower crypto update**
- Place signed certificates in the `received-certs` directory 
- Place issuing ca and root ca certificates in `trust-root` directory. You can also place trust certs together with other certificates.
- Copy certificates: `apicupcfg -certcopy -certdir received-certs [-trustdir trust-root]`
- Change to the *datapower* directory
    - run *zoma-crypto-update-datapower-name.bat|sh* script for each datapower instance

**Typical Steps for the Kubernetes install.** 
- Create working directory.
- Generate json configuration file *subsys-config.json*:
    - `apicupcfg -initconfig -configtype k8s`
- Edit *subsys-config.json* configuration file. In what follows, *tag* is the value of the *Tag* property.
- Validate subsystem ip addresses:
    - `apicupcfg -validateip`
- Generate subsystem and certificate setting scripts:
    - `apicupcfg -gen`
- Generate CSR's. 
    - From the *custom-csr* directory:
        - `all-user-facing-public-csr.tag.bat|sh`
- Submit generated csr files to the ca.
- Place certificates recieved from the ca in a directory. You can optionally split certificates in one directory (received-certs)
and ca certificates in another directory (trust-root). 
- Copy certificates to correct destinations:
    - `apicupcfg -certcopy -certdir dir [-trustdir dir] `
- Run subsystem and certificate setting scripts. 
    - From the *project* directory: 
        - `..\apicup-subsys-set-management.tag.bat|sh`
        - `..\apicup-subsys-set-analytics.tag.bat|sh`
        - `..\apicup-subsys-set-portal.tag.bat|sh`
        - `..\apicup-subsys-set-gateway.tag.bat|sh`
        - `..\apicup-certs-set-user-facing-public.tag.bat|sh`
- Install subsystems with the `apicup subsys install` command as usual. 
    - From the *project* directory:
        - `..\bin\apicup subsys install mgmt`
        - `..\bin\apicup subsys install gwy`
        - `..\bin\apicup subsys install alyt`
        - `..\bin\apicup subsys install ptl`

**General Info.**

*No IBM API Connect software is required to build or use this tool.*

`apicupcfg` generates all subsystem and certificate scripts, all supportring certificate configuration
scripts, datapower configuration scripts, and provides a number of validation options.

To see all available command line options run:

`apicupcfg -h` on linux or macos
`apicupcfg.exe -h` on windows

`apicupcfg` collects all required configuration data in one input json file.

The directory with the configuration file is working directory.

To create input configuration file, change to the working directory and generate json 
configuration file template:

`apicupcfg -initconfig -configtype ova|k8s [-config subsys-config.json]`

`subsys-config.json` is the default configuration file name and is assumed when the `-config` option 
is ommitted.

Ova and kuberenetes configration files are different.

There are a number of sections in the input configuration file. There is a section for each subsystem,
certificate configuration section, and defaults section that applies to all subystems.

To generate configuration scripts and directories:
`apicupcfg -gen [-out output-directory] [-config subsys-config.json]`

Default output directory is *current* directory and is assumed if `-out` command line option is ommitted.
Default input configuration file is *subsys-config.json*

To generate configuration scripts and directories in the working directory:
`apicupcfg -gen [-out .] [-config config-subsys.json]`

Generated directories:

- *bin* - directory for apicup executable. Copy your version of apicup executable to the bin directory and then copy it to the apicup executable.
- *project* - apicup install configuration.
- *custom-csr* - openssl scripts for custom certificates.
- *common-csr* - openssl scripts for common certificates.
- *datapwer* - datapower configuration scripts.
  
Generated scripts and configuration data depend on the input configuration file. 
General rule is that all required scripts and configuration files are created. 

Each generated configuration file is tagged with the value defined in the input json configuration file. 

Generated subsystem configuration scripts: 
 
- `apicup-subsys-set-management.tag.(sh|bat)`
- `apicup-subsys-set-analytics.tag.(sh|bat)`
- `apicup-subsys-set-portal.tag.(sh|bat)`
- k8s install only: `apicup-subsys-set-gateway.tag.(sh|bat)`
 
Generated certificate setting scripts.
 
- `apicup-certs-set-user-facing-public.tag.(sh|bat)`
- `apicup-certs-set-user-facing-public.tag.(sh|bat)`
- `apicup-certs-set-mutual-auth.tag.(sh|bat)`
- `apicup-certs-set-common.tag.(sh|bat)`
 
Generated datapower scripts. (OVA install only) 

- `all-datapower-csr.tag.bat|sh`
- `*.conf` one ssl configuration file for all endpoints and datapower machines`
- `one *.conf.bat|sh` open ssl script
- `zoma-crypto-self-*.bat|sh` script for each datapower machine for self-signed cryto.
- `zoma-crypto-update-*.bat|sh` crypto script for each datapower machine
- `zoma-web-gui-*.bat|sh` web gui configuration scripts
- `zoma-*.bat|sh` datapower configuration script for each datapower machine
- `soma` subdirectory with datapower xml request files.
 
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

*custom-csr* subdirectory contains openssl scripts for public user facing certs and public certs.
*common-csr* subdirectory contains openssl scripts for common certificates, this includes mutual auth and client certs.

Each endpoint the json configuration file is transformed into the cn and openssl csr configuration
file and script are generated in the csr subdirectory.

The Certs.DnFields value defines components of the dn that are required by the certificate authority
and included in the csr configuration file.

To simplify running of these scripts, the custom-csr subdirectory contains 
`all-user-facing-public-csr.tag.(sh|bat)` and `all-public-csr.tag.(sh|bat)` scripts.

The common-csr subdirectory contains `all-mutual-auth-csr.tag.(sh|bat)` and `all-common-csr.tag.sh` scripts.

Run these scripts to generate key pairs and csr's. Submit csr's to the certifiacte authority to get signed certificates.
Copy received certificates to the correct destination.

**Copying certificates to the correct files**.

Certificate settings scripts expect to find certificates, private keys and root certificates at specific locations. 

To copy certificate chain received from ca:  
`apicupcfg -certcopy -cert cert.pem -ca ca.pem -rootca rootca.pem` 

This command will introspect the certificate, build validation path, match subject common name with endpoints defined 
in the configuration file and copy certificate chain to the correct destinations. As part of a copy process
`ca.pem` and `rootca.pem` will be merged into certificate file. Each copy will be recorded in the cert-copy.log. 

Note that if certificate subject common name matches mulitple endpoints then a separate copy will be made for each endpoint. 

To process all certificates received from ca together, place them in one directory or put certificates in one
directory (certdir) and ca and ca root certificates in another directory (trustdir) and run:  
`apicupcfg -certcopy -certdir dir [-trustdir dir]` 

This command will build trust chains for all certificates in the certificate directory and copy each certificate
chain as described above.

**Verifying certificates**

To verify certificate:
`apicupcfg -certverify [-noexpire] -cert path-to-cert.pem -ca path-to-intermediate-ca.pem -rootca path-to-root-ca.pem`

To verify intermidiate ca certificate:
`apicupcfg -certverify [-noexpire] -ca path-to-intermediate-ca.pem -rootca path-to-root-ca.pem`

This command will compute and display trust chain. Pass `-noexpire` to ignore certificate expiration. 

**Datapower Configuration for the OVA install.**  

To configure datapower cluster, define configuration values in the Gateway:{} object. 

Datapower configuration is generated in the *datapower* directory.  

Datapower configuration defines 2 endpoints: gateway director endpoint and api invocation endpoint.
One csr is generated for all endpoints with subject alternative names listing all datapower instances.

To complete datapower configuration, change to the *datapower* directory and run *all-datapower-csr-tag.bat|shell* script.
This creates private key, csr, and self-signed certificates.

Datapower crypto is first configured with self-signed certificates. Real certificates are installed with the crypto update script.

There are a number of layers in the *Datapower* configuration. Each layer is configured with the *SOMA* request.
Each *SOMA* request is named with the function that it executes, eg `dp-domain.xml` to create application domain.

*SOMA* request is posted to the target datapower by the apicupcfg command:
`apicupcfg -config ../subsys-config.json -soma -req dp-domain.xml -auth dp.env -url https://gw1.my.domain:5550/service/mgmt/3.0`

A number of manual datapower configuration steps is kept to the minimum.
Complete initial datapower setup, set timezone, and enable xml management interface.

Datapower configuration steps are combined into the *zoma* scripts, one for each datapower instance.

Create *dp.env* file with 2 lines, one username and another password for datapower authentication.
dp.env:
admin
dppassword 

Run *zoma* self-signed scirpts for each individual datapower. 
Run *zoma* script for each individual datapower. 
Run *zoma* crypto update script for each individual datapower when signed certs are ready. 
Run *zoma-web-gui-*" scripts for each individual datapower. 

*Copying datapower certificates.*  
Datapower certificates can be copied with either command (see above): 
`apicupcfg -certcopy -cert cert.pem -ca ca.pem -rootca rootca.pem` 
`apicupcfg -certcopy -certdir dir [-trustdir dir]` 

*Updating datapower crypto configuration.*  
After copying datapower certificates and datapower trust certificates run
`zoma-crypto-update...bash|bat` script for each datapower machine.

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

**Command line reference.**

Note that default output directory is current directory: -out .

- help:  
`apicupcfg -help`  
- init config:  
`apicupcfg -initconfig -configtype ova|k8s [-config subsys-config.json]`  
- generate subsys, cert and datapower scripts:  
`apicupcfg -gen [-out .] [-config subsys-config.json]`  
- generate subsys or certs or datapower only:  
`apicupcfg -gen -subsys [-out .] [-config subsys-config.json]`   
`apicupcfg -gen -certs [-out .] [-config subsys-config.json]` 
`apicupcfg -gen -datapower [-out .] [-config subsys-config.json]` 
- copy certificate(s) to correct destination:  
`apicupcfg -certcopy -cert cert.pem -ca ca.pem -rootca rootca.pem [-out .] [-config subsys-config.json]` 
`apicupcfg -certcopy -certdir dir [-trustdir dir] [-out .] [-config subsys-config.json]` 
- verify certificate:  
`apicupcfg -certverify [-noexpire] [-cert cert.pem] -ca ca.pem -rootca rootca.pem`  
- validate subsystem ip addresses (ova install only):  
`apicupcfg -validateip [-config subsys-config.json]`

**Configuraton reference.** 

Create *subsys-config.json* configuration file for the **OVA** install: `apicupcfg -initconfig`. 
Comments are not part of *JSON* syntax.
```
{

    //
    // install type
    //
    "InstallType": "ova",

    //
    // example: windows_lts_v2018.4.1.9, linux_lts_v2018.4.1.9, mac_lts_v2018.4.1.9
    //
    "Version": "windows_lts_v2018.4.1.9", 

    //
    // tag generated scripts
    //
    "Tag": "tag",

    //
    // if set to 'true', then generated scripts use full version for the apicup executable.
    //
    "UseVersion": false, 

    //
    // apicup install mode
    //
    "Mode": "dev|standard",

    //
    // path to the public key for the ssh login to the management, portal, and analytics vm's.
    // put this file into working directory
    //
    "SshPublicKeyFile": "/path/to/public/key/file",

    //
    // a list of search domains, applies to management, analytics, and portal subsystems.
    //
    "SearchDomains": ["my.domain.com", "domain.com"],

    //
    // these parameters take effect at first boot only.
    // changing these parameters after the first boot does not take effect.
    // applies to managment, portal, and analytics subystems.
    //
    "VmFirstBoot": {
        "DnsServers": ["dns-ip-1","dns-ip-2"],
        "VmwareConsolePasswordHash": "hash-output-b64",

        "IpRanges": {
            "PodNetwork": "172.16.0.0/16",
            "ServiceNetwork": "172.17.0.0/16"
        }
    },

    //
    // cloud-init file. If CloudInitFile is empty no file will be generated.
    //
    "CloudInit": {
        "CloudInitFile": "cloud-init-file.yml",
        "InitValues": {
            "rsyslog": {
                "remotes": {
                    "syslog_server1": "syslog-collector-ip-1:514|601",
                    "syslog_server2": "syslog-collector-ip-2:514|601"
                }
            }
        }
    },

    //
    // certificates
    //
    "Certs": {
        //
        // dn fields that must be present in the csr.
        // note that state is ST
        //
        "DnFields": ["O=APIC","C=US"],
        
        //
        // kubernetes namespace. Do not change for the ova install
        //
        "K8sNamespace": "default",
        
        //
        // Types of certifiactes to generate.
        //
        // for the public user facing certs all open ssl scripts will be in the custom-csr directory
        // to generate key pairs and csr's change to the custom-csr directory and run:
        // all-user-facing-public-csr.tag.(bat|sh)
        //
        // submit csr's to ca
        //
        // To copy signed certificates to files expected by scripts: 
        // apicupcfg -certcopy -certdir dir [-trustdir dir]
        // 
        "PublicUserFacingCerts": true,

        //
        // public-certs and common-certs are advanced options.
        //
        "PublicCerts": false,
        "CommonCerts": false
    },

    //
    // management subsystem
    //
    "Management": {
        //
        // management subsystem name
        //
        "SubsysName": "mgmt",

        //
        // management subsystem parameters for the first boot
        //
        "VmFirstBoot": {
        
            //
            // management subsystem hosts. 1 host for the dev mode, 3 hosts for the standard mode.
            //
            "Hosts": [
                {"Name": "m1.my.domain.com", "HardDiskPassword": "password", "Device": "eth0",
                    "IpAddress": "ip-address", "SubnetMask": "dot.subnet.mask","Gateway": "gw-ip-address"},
                {"Name": "m2.my.domain.com", "HardDiskPassword": "password", "Device": "eth0",
                    "IpAddress": "ip-address", "SubnetMask": "dot.subnet.mask", "Gateway": "gw-ip-address"},
                {"Name": "m3.my.domain.com", "HardDiskPassword": "password", "Device": "eth0",
                    "IpAddress": "ip-address", "SubnetMask": "dot.subnet.mask", "Gateway": "gw-ip-address"}
            ]
        },

        //
        // Cassandra backup configuration.
        // If BackupProtocol is empty, backup configuration is skipped.
        //
        "CassandraBackup": {
            //
            // Backup protocol. 
            // For the sftp protocol, specify Backup* parameters. For the objstore protocol specify Objstore* parameters.
            //
            "BackupProtocol": "sftp|objstore",
            "BackupAuthUser": "admin",
            "BackupAuthPass": "secret",
            "BackupHost": "backup.my.domain.com",
            "BackupPort": 1022,
            "BackupPath": "/backup",
            "ObjstoreS3SecretKeyId": "",
            "ObjstoreS3SecretAccessKey": "",
            "ObjstoreEndpointRegion": "",
            "ObjstoreBucketSubfolder": "",
            "BackupEncoding": "min(0-59) hour(0-23) dayofmonth(1-31) month(1-12) dayofweek(0-6)",
            "BackupSchedule": "0 0 * * 0"
        },

        //
        // Cassandra backup encryption key.
        //
        "CassandraEncryptionKeyFile": "encryption-secret.bin",

        //
        // Management subsystem endpoints
        //
        "PlatformApi": "platform.my.domain.com",
        "ApiManagerUi": "ui.my.domain.com",
        "CloudAdminUi": "cm.my.domain.com",
        "ConsumerApi": "consumer.my.domain.com"
    },

    //
    // Analytics subsystem
    //
    "Analytics": {
        //
        // Analytics subsystem name
        //
        "SubsysName": "alt",

        //
        // analytics subsystem parameters for the first boot
        //
        "VmFirstBoot": {
        
            //
            // analytics subsystem hosts. 1 host for the dev mode, 3 hosts for the standard mode.
            //
            "Hosts": [
                {"Name": "a1.my.domain.com", "HardDiskPassword": "password", "Device": "eth0",
                    "IpAddress": "ip-address", "SubnetMask": "dot.subnet.mask","Gateway": "gw-ip-address"},
                {"Name": "a2.my.domain.com", "HardDiskPassword": "password", "Device": "eth0",
                    "IpAddress": "ip-address", "SubnetMask": "dot.subnet.mask", "Gateway": "gw-ip-address"},
                {"Name": "a3.my.domain.com", "HardDiskPassword": "password", "Device": "eth0",
                    "IpAddress": "ip-address", "SubnetMask": "dot.subnet.mask", "Gateway": "gw-ip-address"}
            ]
        },

        "EnableMessageQueue": false,

        //
        // analytics subsystem endpoints
        //
        "AnalyticsIngestion": "ai.my.domain.com",
        "AnalyticsClient": "ac.my.domain.com"
    },

    //
    // Portal subsystem
    //
    "Portal": {
        //
        // portal subsystem name
        //
        "SubsysName": "ptl",

        //
        // portal subsystem parameters for the first boot
        //
        "VmFirstBoot": {
            //
            // portal subsystem hosts. 1 host for the dev mode, 3 hosts for the standard mode.
            //
            "Hosts": [
                {"Name": "p1.my.domain.com", "HardDiskPassword": "password", "Device": "eth0",
                    "IpAddress": "ip-address", "SubnetMask": "dot.subnet.mask","Gateway": "gw-ip-address"},
                {"Name": "p2.my.domain.com", "HardDiskPassword": "password", "Device": "eth0",
                    "IpAddress": "ip-address", "SubnetMask": "dot.subnet.mask", "Gateway": "gw-ip-address"},
                {"Name": "p3.my.domain.com", "HardDiskPassword": "password", "Device": "eth0",
                    "IpAddress": "ip-address", "SubnetMask": "dot.subnet.mask", "Gateway": "gw-ip-address"}
            ]
        },

        //
        // Portal site backup configuration.
        // If BackupProtocol is empty, backup configuration is skipped.
        //
        "SiteBackup": {
            //
            // Backup protocol. 
            // For the sftp protocol, specify Backup* parameters. For the objstore protocol specify Objstore* parameters.
            //
            "BackupProtocol": "sftp|objstore",
            "BackupAuthUser": "admin",
            "BackupAuthPass": "secret",
            "BackupHost": "backup.my.domain.com",
            "BackupPort": 1022,
            "BackupPath": "/backup",
            "ObjstoreS3SecretKeyId": "",
            "ObjstoreS3SecretAccessKey": "",
            "ObjstoreEndpointRegion": "",
            "ObjstoreBucketSubfolder": "",
            "BackupEncoding": "min(0-59) hour(0-23) dayofmonth(1-31) month(1-12) dayofweek(0-6)",
            "BackupSchedule": "0 2 * * *"
        },

        //
        // portal subsystem endpoints
        //
        "PortalAdmin": "padmin.my.domain.com",
        "PortalWww": "portal.my.domain.com"
    },

    //
    // Datapower gateway
    //
    "Gateway": {
        //
        // Gateway subsystem name. Not used for the OVA install.
        //
        "SubsysName": "gwy",
        
        //
        // Datapower install mode.
        //
        "Mode": "standard",

        //
        // datapower search domains
        //
        "SearchDomains": ["my.domain.com","domain.com"],
        
        //
        // datapower dns servers
        //
        "DnsServers": ["dns-ip-1","dns-ip-2"],

        //
        // datapower ntp server
        //
        "NTPServer": "pool.ntp.org",

        //
        // Datapower hosts. 1 host in dev mode, 3 hosts in standard mode
        // [...] denotes optional fields with default value.
        //
        "Hosts": [
            {"Name": "dp1.my.domain.com", "Device": "eth0", ["HostAlias": "if_eth0",]
                "IpAddress": "ip-address", "SubnetMask": "dot.subnet.mask","Gateway": "gw-ip-address",
                "GwdPeeringPriority": 100, "RateLimitPeeringPriority": 100, "SubsPeeringPriority": 100,
                "ApiProbePeeringPriority: 100},

            {"Name": "dp2.my.domain.com", "Device": "eth0", ["HostAlias": "if_eth0",]
                "IpAddress": "ip-address", "SubnetMask": "dot.subnet.mask", "Gateway": "gw-ip-address",
                "GwdPeeringPriority": 90, "RateLimitPeeringPriority": 90, "SubsPeeringPriority": 90,
                "ApiProbePeeringPriority": 90},

            {"Name": "dp3.my.domain.com", "Device": "eth0", ["HostAlias": "if_eth0",]
                "IpAddress": "ip-address", "SubnetMask": "dot.subnet.mask", "Gateway": "gw-ip-address",
                "GwdPeeringPriority": 80, "RateLimitPeeringPriority": 80, "SubsPeeringPriority": 80,
                "ApiProbePeeringPriority": 80}
        ],

        //
        // datapower endpoints
        //
        "ApiGateway": "gw.my.domain.com",
        "ApicGwService": "gwd.my.domain.com",

        //
        // datapower domain name for api connect
        //
        "DatapowerDomain": "apiconnect",
        
        //
        // datapower api gateway port
        //
        "DatapowerApiGatewayPort": 9443,
    }
}
```

Create *subsys-config.json* configuration file for the **Kubernetes** install: `apicupcfg -initconfig -configtype k8s`. 
Comments are not part of *JSON* syntax. 
```
{

    //
    // install type
    //    
    "InstallType": "k8s",
 
    //
    // example: windows_lts_v2018.4.1.9, linux_lts_v2018.4.1.9, mac_lts_v2018.4.1.9
    //
    "Version": "windows_lts_v2018.4.1.9", 

    //
    // tag generated scripts
    // 
    "Tag": "tag",

    //
    // if set to 'true', then generated scripts use full version for the apicup executable.
    //
    "UseVersion": false,

    //
    // apicup install mode
    //
    "Mode": "dev|standard",

    //
    // kubernetes namespace
    //
    "Namespace": "apic",

    //
    // registry url for container images
    //
    "RegistryUrl": "container-image-registry-url",

    //
    // registry secret to pull images
    //
    "RegistrySecret": "container-image-registry-secret",

    //
    // ingress type: route for open-shift, ingress for all other kubernetes clusters
    //
    "IngressType": "ingress|route",

    //
    // storage class
    //
    "StorageClass": "gp2|etc",

    //
    // certificates
    //
    "Certs": {

        //
        // dn fields that must be present in the csr.
        // note that state is ST
        //
        "DnFields": ["O=APIC","C=US"],

        //
        // kubernetes namespace. match to the value defined in the Namespace parameter
        //
        "K8sNamespace": "apic",

        //
        // Types of certifiactes to generate.
        //
        // for the public user facing certs all open ssl scripts will be in the custom-csr directory
        // to generate key pairs and csr's change to the custom-csr directory and run:
        // all-user-facing-public-csr.tag.(bat|sh)
        //
        // submit csr's to ca
        //
        // To copy signed certificates to files expected by scripts: 
        // apicupcfg -certcopy -certdir dir [-trustdir dir]
        // 
        "PublicUserFacingCerts": true,

        //
        // public-certs and common-certs are advanced options.
        //
        "PublicCerts": false,
        "CommonCerts": false
    },

    //
    // management subsystem
    //
    "Management": {
        //
        // management subsystem name
        //
        "SubsysName": "mgmt",

        //
        // extra values file for the management subsystem chart
        //
        "ExtraValuesFile": "mgmt-values.yaml",

        //
        // values for the management subsystem chart
        //
        "ExtraValues": {},

        //
        // Cassandra backup configuration.
        // If BackupProtocol is empty, backup configuration is skipped.
        //
        "CassandraBackup": {
            //
            // Backup protocol. 
            // For the sftp protocol, specify Backup* parameters. For the objstore protocol specify Objstore* parameters.
            //
            "BackupProtocol": "sftp|objstore",
            "BackupAuthUser": "admin",
            "BackupAuthPass": "secret",
            "BackupHost": "backup.my.domain.com",
            "BackupPort": 1022,
            "BackupPath": "/backup",
            "ObjstoreS3SecretKeyId": "",
            "ObjstoreS3SecretAccessKey": "",
            "ObjstoreEndpointRegion": "",
            "ObjstoreBucketSubfolder": "",
            "BackupEncoding": "min(0-59) hour(0-23) dayofmonth(1-31) month(1-12) dayofweek(0-6)",
            "BackupSchedule": "0 0 * * 0"
        },

        "CassandraMaxMemoryGb": 9,
        "CassandraVolumeSizeGb": 50,
        "CassandraClusterSize": 3,
        "ExternalCassandraHost": "ext.my.domain.com",
        "CreateCrd": true,

        //
        // Cassandra backup encryption key
        //
        "CassandraEncryptionKeyFile": "encryption-secret.bin",

        //
        // Management subsystem endpoints
        //
        "PlatformApi": "api.my.domain.com",
        "ApiManagerUi": "apim.my.domain.com",
        "CloudAdminUi": "cm.my.domain.com",
        "ConsumerApi": "consumer.my.domain.com"
    },

    //
    // Analytics subsystem
    //
    "Analytics": {
        //
        // Analytics subsystem name
        //
        "SubsysName": "analyt",

        //
        // extra values file for the analytics subsystem chart
        //
        "ExtraValuesFile": "analyt-values.yaml",

        //
        // values for the analytics subsystem chart
        //
        "ExtraValues": {},

        "CoordinatingMaxMemoryGb": 12,
        "DataMaxMemoryGb": 12,
        "DataStorageSizeGb": 200,
        "MasterMaxMemoryGb": 12,
        "MasterStorageSizeGb": 5,

        "EnableMessageQueue": false,

        "EsStorageClass": "",
        "MqStorageClass": "",

        //
        // analytics subsystem endpoints
        //
        "AnalyticsIngestionEndpoint": "ai.my.domain.com",
        "AnalyticsClientEndpoint": "ac.my.domain.com"
    },

    //
    // Portal subsystem
    //
    "Portal": {
        //
        // portal subsystem name
        //
        "SubsysName": "ptl",

        //
        // extra values file for the portal subsystem chart
        //
        "ExtraValuesFile": "ptl-values.yaml",

        //
        // extra values for the portal subsystem chart
        //
        "ExtraValues": {},

        //
        // Portal backup configuration.
        // If BackupProtocol is empty, backup configuration is skipped.
        //
        "SiteBackup": {
            //
            // Backup protocol. 
            // For the sftp protocol, specify Backup* parameters. For the objstore protocol specify Objstore* parameters.
            //
            "BackupProtocol": "sftp|objstore",
            "BackupAuthUser": "admin",
            "BackupAuthPass": "secret",
            "BackupHost": "backup.my.domain.com",
            "BackupPort": 1022,
            "BackupPath": "/backup",
            "ObjstoreS3SecretKeyId": "",
            "ObjstoreS3SecretAccessKey": "",
            "ObjstoreEndpointRegion": "",
            "ObjstoreBucketSubfolder": "",
            "BackupEncoding": "min(0-59) hour(0-23) dayofmonth(1-31) month(1-12) dayofweek(0-6)",
            "BackupSchedule": "0 0 * * 0"
        },

        "WwwStorageSizeGb": 5,
        "BackupStorageSizeGb": 5,
        "DbStorageSizeGb": 12,

        //
        // parameters that can not be changed for the portal subsystem
        //
        "Fixed": {
            "DbLogsStorageSizeGb": 2,
            "AdminStorageSizeGb": 1
        },

        //
        // portal subsystem endpoints
        //
        "PortalAdmin": "padmin.my.domain.com",
        "PortalWWW": "portal.my.domain.com"
    },

    //
    // Gateway subsystem
    //
    "Gateway": {
        //
        // Gateway subsystem name
        //
        "SubsysName": "gwy",

        //
        // gateway subsystem deployment mode
        //
        "Mode": "dev",

        //
        // extra values file for the gateway subsystem chart
        //
        "ExtraValuesFile": "gwy-values.yaml",

        //
        // extra values for the gateway subsystem chart
        //
        "ExtraValues": {
            "datapower": {
                "webGuiManagementState": "enabled",
                "apiDebugProbe": "enabled"
            }
        },

        "LicenseVersion": "Production|Development",
        "ImagePullPolicy": "IfNotPresent",

        "ReplicaCount": 3,
        "MaxCpu": 4,
        "MaxMemoryGb": 6,

        "V5ComatabilityMode": false,
        "EnableTms": true,
        "TmsPeeringStorageSizeGb": 10,
        "EnableHighPerformancePeering": "true",

        //
        // Gateway subsystem endpoints
        //
        "ApiGateway": "gw.my.domain.com",
        "ApicGwService": "gwd.my.domain.com"
    }
}
```