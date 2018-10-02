package main

import (
	"flag"
	"time"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/models"
	"github.com/evilsocket/shellz/session"
)

var (
	command   = ""
	onFilter  = "*"
	onNames   = []string{}
	on        = models.Shells{}
	toOutput  = ""
	doList    = false
	doEnable  = ""
	doDisable = ""
	err       = error(nil)
	idents    = models.Identities(nil)
	shells    = models.Shells(nil)

	timeouts = session.Timeouts{
		Connect: 5 * time.Second,
		Read:    500 * time.Millisecond,
		Write:   500 * time.Millisecond,
	}
)

func init() {
	flag.BoolVar(&doList, "list", doList, "List available shells and exit.")

	flag.StringVar(&doEnable, "enable", "", "Enable the specified shells.")
	flag.StringVar(&doDisable, "disable", "", "Disable the specified shells.")

	flag.StringVar(&command, "run", command, "Command to run on the selected shells.")
	flag.StringVar(&onFilter, "on", onFilter, "Comma separated list of shell names to select or * for all.")
	flag.StringVar(&toOutput, "to", toOutput, "If filled, commands output will be saved to this file instead of being printed on the standard output.")

	flag.DurationVar(&timeouts.Connect, "connection-timeout", timeouts.Connect, "Connection timeout.")
	flag.DurationVar(&timeouts.Read, "read-timeout", timeouts.Read, "Read timeout.")
	flag.DurationVar(&timeouts.Write, "write-timeout", timeouts.Write, "Write timeout.")

	flag.BoolVar(&log.DebugMessages, "debug", log.DebugMessages, "Enable debug messages.")
	flag.Parse()
}

func main() {
	log.Raw(core.Banner)

	if err, idents, shells = models.Load(); err != nil {
		log.Fatal("error while loading data: %s", err)
	} else if len(shells) == 0 {
		log.Fatal("no shells found on the system, start creating json files inside %s", models.Paths["shells"])
	} else if err = session.LoadPlugins(models.Paths["plugins"]); err != nil {
		log.Fatal("error while loading plugins: %s", err)
	} else {
		log.Debug("loaded %d identities and %d shells", len(idents), len(shells))
	}

	if doList {
		showList()
	} else if doEnable != "" {
		runEnable(doEnable, true)
	} else if doDisable != "" {
		runEnable(doDisable, false)
	} else if command != "" {
		runCommand()
	} else {
		showHelp()
	}
}
