package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
)

var (
	examples = []struct {
		cmd  string
		help string
	}{
		{"shellz -list", "list available identities and shells"},

		{"shellz -enable \"machineA, machineB\"", "enable the shells named machineA and machineB"},
		{"shellz -disable machineA", "disable the shell named machineA (commands won't be executed on it)"},

		{"shellz -run id", "run the command 'id' on each shell"},
		{"shellz -run id -on machineA", "run the command 'id' on a single shell named 'machineA'"},
		{"shellz -run id -on 'machineA, machineB'", "run the command 'id' on machineA and machineB"},
		{"shellz -run uptime -to all.txt", "run the command 'uptime' on every shell and append all outputs to the 'all.txt' file"},
		{"shellz -run uptime -to \"{{.Identity.Username}}_{{.Name}}.txt\"", "run the command 'uptime' on every shell and save each outputs to a different file using per-shell data."},
	}
)

func showHelp() {
	log.Info("none of the -run or -list parameters have been specified")

	fmt.Println()
	fmt.Printf("Usage:\n\n")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Printf("Examples:\n\n")

	for _, e := range examples {
		fmt.Printf("  %s\n", core.Dim("# "+e.help))
		fmt.Printf("  %s\n", core.Bold(e.cmd))
		fmt.Println()
	}

	os.Exit(1)
}
