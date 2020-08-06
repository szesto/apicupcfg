package apicupcfg

import "encoding/json"

type Backup struct {
	Sftp SftpBackupConfig
	Objstore ObjstoreBackupConfig
	V10 V10BackupConfig
	Schedule BackupSchedule
	v10Flag bool
}

func (b *Backup) gen(v10, verbose bool) {
	if v10 {
		b.V10.gen()

	} else {
		b.Sftp.gen()
		b.Objstore.gen()
	}

	b.Schedule.gen()

	b.v10Flag = v10
}

func (b *Backup) MarshalJSON() ([]byte, error) {

	if b.v10Flag {
		return json.Marshal(&struct {
			V10 V10BackupConfig
			Schedule BackupSchedule
		}{b.V10, b.Schedule})

	} else {
		return json.Marshal(&struct {
			Sftp SftpBackupConfig
			Objstore ObjstoreBackupConfig
			Schedule BackupSchedule
		}{b.Sftp, b.Objstore, b.Schedule})
	}
}

type BackupSchedule struct {
	Min059 string
	Hour023 string
	DayOfMonth131 string
	Month112 string
	DayOfWeek06 string
}

func (b *BackupSchedule) gen() {
	b.Min059 = "*"
	b.Hour023 = "*"
	b.DayOfMonth131 = "*"
	b.Month112 = "*"
	b.DayOfWeek06 = "*"
}

type V10BackupConfig struct {
	Enabled bool
	S3Provider string // v10: aws/ibm
	Host string
	Path string
	Retries int // v10: number of retries
	Credentials string // v10: kube secret name
}

func (b *V10BackupConfig) gen() {
	b.Enabled = false
	b.S3Provider = "ibm/aws"
	b.Host = "backup-host-more-detail-please"
	b.Path = "bucket/subfolder"
	b.Retries = 0
	b.Credentials = "k8s-s3-auth-secret-name"
}

// objstore://s3-secret-key-id@s3-secret-access-key/endpoint/region/bucket/subfolder
type SftpBackupConfig struct {
	Enabled bool
	BackupAuthUser string
	BackupAuthPass string
	BackupHost string
	BackupPort int
	BackupPath string
}

func (b *SftpBackupConfig) gen() {
	b.Enabled = false
	b.BackupAuthUser = "admin"
	b.BackupAuthPass = "password"
	b.BackupHost = "backup-host-fqdn"
	b.BackupPort = 2222
	b.BackupPath = "/backup"
}

type ObjstoreBackupConfig struct {
	Enabled bool
	ObjstoreS3SecretKeyId string // -> auth-user
	ObjstoreS3SecretAccessKey string // -> auth-pass
	ObjstoreEndpointRegion string // endpoint/region -> host, v10 reuse
	ObjstoreBucketSubfolder string // bucket|bucket/subfolder -> backup-path, v10 reuse
}

func (b *ObjstoreBackupConfig) gen() {
	b.Enabled = false
	b.ObjstoreS3SecretKeyId = "s3-secret-key-id"
	b.ObjstoreS3SecretAccessKey = "s3-secret-access-key"
	b.ObjstoreEndpointRegion = "endpoint/region"
	b.ObjstoreBucketSubfolder = "bucket/subfolder"
}
