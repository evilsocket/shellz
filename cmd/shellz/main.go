package main

import (
	"flag"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/models"
)

var (
	command  = ""
	onFilter = "*"
	onNames  = []string{}
	on       = models.Shells{}
	toOutput = ""
	doList   = false
	err      = error(nil)
	idents   = models.Identities(nil)
	shells   = models.Shells(nil)
)

func init() {
	flag.BoolVar(&doList, "list", doList, "List available shells and exit.")

	flag.StringVar(&command, "run", command, "Command to run on the selected shells.")
	flag.StringVar(&onFilter, "on", onFilter, "Comma separated list of shell names to select or * for all.")
	flag.StringVar(&toOutput, "to", toOutput, "If filled, commands output will be saved to this file instead of being printed on the standard output.")

	flag.BoolVar(&log.DebugMessages, "debug", log.DebugMessages, "Enable debug messages.")
	flag.Parse()
}

func main() {
	log.Raw(core.Banner)

	if err, idents, shells = models.Load(); err != nil {
		log.Fatal("error while loading identities and shells: %s", err)
	} else if len(shells) == 0 {
		log.Fatal("no shells found on the system, start creating json files inside %s", models.Paths["shells"])
	} else {
		log.Debug("loaded %d identities and %d shells", len(idents), len(shells))
	}

	if doList {
		showList()
	} else if command != "" {
		runCommand()
	} else {
		showHelp()
	}
}
