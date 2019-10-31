package apicupcfg

import rice "github.com/GeertJohan/go.rice"

type Backup struct {
	// sftp/objstore
	BackupProtocol string

	// sftp backup
	BackupAuthUser string
	BackupAuthPass string
	BackupHost string
	BackupPort int
	BackupPath string

	// objstore backup
	// objstore://s3-secret-key-id@s3-secret-access-key/endpoint/region/bucket/subfolder
	ObjstoreS3SecretKeyId string // -> auth-user
	ObjstoreS3SecretAccessKey string // -> auth-pass
	ObjstoreEndpointRegion string // endpoint/region -> host
	ObjstoreBucketSubfolder string // bucket|bucket/subfolder -> backup-path

	BackupSchedule string
}

type SubsysK8sBase struct {
	OsEnv

	SubsysName string
	Mode string

	ExtraValuesFile string
	ExtraValues map[string]interface {}

	IngressType string
	Namespace string
	RegistryUrl string
	RegistrySecret string
	StorageClass string
}

func (b *SubsysK8sBase) copyDefaults(from SubsysK8s) {
	b.OsEnv.copyDefaults(from.OsEnv)

	if len(b.Mode) == 0 {
		b.Mode = from.Mode
	}

	if len(b.Namespace) == 0 {
		b.Namespace = from.Namespace
	}

	if len(b.RegistryUrl) == 0 {
		b.RegistryUrl = from.RegistryUrl
	}

	if len(b.RegistrySecret) == 0 {
		b.RegistrySecret = from.RegistrySecret
	}

	if len(b.IngressType) == 0 {
		b.IngressType = from.IngressType
	}

	if len(b.StorageClass) == 0 {
		b.StorageClass = from.StorageClass
	}
}

type MgtSubsysK8s struct {
	SubsysK8sBase

	// cassandra
	CassandraBackup Backup

	CassandraMaxMemoryGb int
	CassandraClusterSize int
	CassandraVolumeSizeGb int
	ExternalCassandraHost string

	CassandraEncryptionKeyFile string

	CreateCrd bool

	// endpoints
	PlatformApi string
	ApiManagerUi string
	CloudAdminUi string
	ConsumerApi string
}

type AlytSubsysK8s struct {
	SubsysK8sBase

	// storage classes
	EsStorageClass string
	MqStorageClass string

	// memory/storage
	CoordinatingMaxMemoryGb int
	DataMaxMemoryGb int
	DataStorageSizeGb int
	MasterMaxMemoryGb int
	MasterStorageSizeGb int

	EnableMessageQueue bool

	// endpoints
	AnalyticsIngestionEndpoint string
	AnalyticsClientEndpoint string
}

type PtlSubsysK8s struct {
	SubsysK8sBase

	SiteBackup Backup

	WwwStorageSizeGb int // >= 5gb
	BackupStorageSizeGb int // >= 5gb
	DbStorageSizeGb int // >= 12gb
	DbLogsStorageSizeGb int // = 12gb
	AdminStorageSizeGb int // = 1gb

	// endpoints
	PortalAdmin string
	PortalWWW string
}

type GwSubsysK8s struct {
	SubsysK8sBase

	LicenseVersion string

	ImagePullPolicy string

	ReplicaCount int
	MaxCpu int
	MaxMemoryGb int

	V5CompatibilityMode bool
	EnableTms bool
	TmsPeeringStorageSizeGb int
	EnableHighPerformancePeering string

	ApiGateway string
	ApicGwService string
}

type SubsysK8s struct {
	InstallTypeHeader
	OsEnv

	Version string
	Tag string

	// defaults
	Mode string
	Namespace string
	RegistryUrl string
	RegistrySecret string
	IngressType string
	StorageClass string

	// certs
	Certs Certs

	// subsys
	Management MgtSubsysK8s
	Analytics AlytSubsysK8s
	Portal PtlSubsysK8s
	Gateway GwSubsysK8s
}

func LoadSubsysK8s(jsonConfigFile string) *SubsysK8s {
	subsys := &SubsysK8s{}
	unmarshallJsonFile(jsonConfigFile, &subsys)

	subsys.OsEnv.init()

	// copy defaults
	subsys.Management.copyDefaults(*subsys)
	subsys.Gateway.copyDefaults(*subsys)
	subsys.Analytics.copyDefaults(*subsys)
	subsys.Portal.copyDefaults(*subsys)

	return subsys
}

func ApplyTemplatesK8s(subsys *SubsysK8s, outfiles map[string]string, subsysOnly, certsOnly bool, tbox *rice.Box)  {

	// parse templates
	mgtt := parseTemplates(tbox, tpdir(tbox) + "management-k8s.tmpl", tpdir(tbox) + "helpers.tmpl")
	gwyt := parseTemplates(tbox, tpdir(tbox) + "gateway-k8s.tmpl", tpdir(tbox) + "helpers.tmpl")
	alytt := parseTemplates(tbox, tpdir(tbox) + "analytics-k8s.tmpl", tpdir(tbox) + "helpers.tmpl")
	ptlt := parseTemplates(tbox, tpdir(tbox) + "portal-k8s.tmpl", tpdir(tbox) + "helpers.tmpl")
	valt := parseTemplates(tbox, tpdir(tbox) + "extra-values.tmpl", tpdir(tbox) + "helpers.tmpl")

	// execute templates
	shellExt := subsys.OsEnv.ShellExt

	isManagement := len(subsys.Management.SubsysName) > 0
	isGateway := len(subsys.Gateway.SubsysName) > 0
	isAnalytics := len(subsys.Analytics.SubsysName) > 0
	isPortal := len(subsys.Analytics.SubsysName) > 0

	outpath := ""

	if isManagement && !certsOnly {
		outpath = fileName(outfiles[outdir], tagOutputFileName(outfiles[managementOut], subsys.Tag)) + shellExt
		writeTemplate(mgtt, outpath, subsys.Management)

		if len(subsys.Management.ExtraValuesFile) > 0 {
			outpath = fileName(outfiles[outdir], subsys.Management.ExtraValuesFile)
			writeTemplate(valt, outpath, subsys.Management.ExtraValues)
		}
	}

	if isGateway && !certsOnly {
		outpath = fileName(outfiles[outdir], tagOutputFileName(outfiles[gatewayOut], subsys.Tag)) + shellExt
		writeTemplate(gwyt, outpath, subsys.Gateway)

		if len(subsys.Gateway.ExtraValuesFile) > 0 {
			outpath = fileName(outfiles[outdir], subsys.Gateway.ExtraValuesFile)
			writeTemplate(valt, outpath, subsys.Gateway.ExtraValues)
		}
	}

	if isAnalytics && !certsOnly {
		outpath = fileName(outfiles[outdir], tagOutputFileName(outfiles[analyticsOut], subsys.Tag)) + shellExt
		writeTemplate(alytt, outpath, subsys.Analytics)

		if len(subsys.Analytics.ExtraValuesFile) > 0 {
			outpath = fileName(outfiles[outdir], subsys.Analytics.ExtraValuesFile)
			writeTemplate(valt, outpath, subsys.Analytics.ExtraValues)
		}
	}

	if isPortal && !certsOnly {
		outpath = fileName(outfiles[outdir], tagOutputFileName(outfiles[portalOut], subsys.Tag)) + shellExt
		writeTemplate(ptlt, outpath, subsys.Portal)

		if len(subsys.Portal.ExtraValuesFile) > 0 {
			outpath = fileName(outfiles[outdir], subsys.Portal.ExtraValuesFile)
			writeTemplate(valt, outpath, subsys.Portal.ExtraValues)
		}
	}

	// certs
	if !subsysOnly {
		updateCertSpecs(&subsys.Certs, &subsys.Management, &subsys.Analytics, &subsys.Portal, &subsys.Gateway,
			outfiles[commonCsrOutDir], outfiles[customCsrOutDir])

		outputCerts(&subsys.Certs, outfiles, subsys.Tag, tbox)
	}
}
