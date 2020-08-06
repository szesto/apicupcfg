package apicupcfg

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestBackup_MarshalJSON(t *testing.T) {

	backup := Backup{}
	backup.v10Flag = true

	jb, _ := json.Marshal(&struct{Backup Backup}{backup})
	var out bytes.Buffer
	_ = json.Indent(&out, jb, "", "\t")
	_, _ = out.WriteTo(os.Stdout)
}
