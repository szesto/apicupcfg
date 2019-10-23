package apicupcfg

import "os"

type OsEnv struct {
	IsLinux bool
	IsWindows bool
	PathSeparator string
	BinApicup string
	ShellExt string
	ScriptInvoke string
}

func (env *OsEnv) init() {
	if os.PathSeparator == '/' {
		env.IsLinux = true
		env.IsWindows = false
		env.BinApicup = "../bin/apicup"
		env.ShellExt = ".sh"
		env.ScriptInvoke = "/bin/bash"
	} else {
		env.IsLinux = false
		env.IsWindows = true
		env.BinApicup = "..\\bin\\apicup"
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