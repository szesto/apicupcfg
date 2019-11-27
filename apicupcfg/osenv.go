package apicupcfg

import (
	"fmt"
	"os"
)

type OsEnv struct {
	IsLinux bool
	IsWindows bool
	PathSeparator string
	BinApicup string
	ShellExt string
	ScriptInvoke string
}

func (env *OsEnv) init() {
	env.init2("", false)
}

func (env *OsEnv) init2(version string, useVersion bool) {
	if os.PathSeparator == '/' {
		env.IsLinux = true
		env.IsWindows = false
		if len(version) > 0 && useVersion {
			// apicup-[linux_lts_v2018.4.1.8-ifix2.0]
			// todo: check for macos... apicup-[mac_lts_v2018.4.1.8-ifix2.0]
			env.BinApicup = fmt.Sprintf("../bin/apicup-%s", version)
		} else {
			env.BinApicup = "../bin/apicup"
		}
		env.ShellExt = ".sh"
		env.ScriptInvoke = "/bin/bash"

	} else {
		env.IsLinux = false
		env.IsWindows = true
		if len(version) > 0 && useVersion {
			// apicup-[windows_lts_v2018.4.1.8-ifix2.0].exe
			env.BinApicup = fmt.Sprintf("..\\bin\\apicup-%s.exe", version)
		} else {
			env.BinApicup = "..\\bin\\apicup.exe"
		}
		env.ShellExt = ".bat"
		env.ScriptInvoke = "call"
	}

	env.PathSeparator = string(os.PathSeparator)
}

func (env *OsEnv) copyDefaults(from OsEnv) {
	env.IsLinux = from.IsLinux
	env.IsWindows = from.IsWindows
	env.PathSeparator = from.PathSeparator
	env.BinApicup = from.BinApicup
	env.ShellExt = from.ShellExt
	env.ScriptInvoke = from.ScriptInvoke
}