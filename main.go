package main

import (
	"flag"
	"fmt"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/models"
)

var (
	runCommand = "uptime"
	onFilter   = "*"
	onNames    = []string{}
	on         = models.Shells{}
)

func init() {
	flag.StringVar(&runCommand, "run", runCommand, "Command to run on the selected shells.")
	flag.StringVar(&onFilter, "on", onFilter, "Comma separated list of shell names to select or * for all.")
	flag.BoolVar(&log.DebugMessages, "debug", log.DebugMessages, "Enable debug messages.")
	flag.Parse()
}

func main() {
	log.Raw(core.Banner)

	err, idents, shells := models.Load()
	if err != nil {
		log.Fatal("error while loading identities and shells: %s", err)
	} else if len(shells) == 0 {
		log.Fatal("no shells found on the system, start creating json files inside %s", models.Paths["shells"])
	} else {
		log.Debug("loaded %d identities and %d shells", len(idents), len(shells))
	}

	if onFilter == "*" {
		on = shells
	} else {
		for _, name := range core.CommaSplit(onFilter) {
			if shell, found := shells[name]; !found {
				log.Fatal("can't find shell %s", name)
			} else {
				on[name] = shell
			}
		}
	}

	if len(on) == 0 {
		log.Fatal("no shell selected by the filter %s", core.Dim(onFilter))
	}

	log.Info("running %s on %d shells ...\n", core.Dim(runCommand), len(on))

	for name, shell := range on {
		err, session := shell.NewSession()
		if err != nil {
			log.Warning("error while creating session for shell %s: %s", name, err)
			continue
		}
		defer session.Close()

		out, err := session.Exec(runCommand)
		if err != nil {
			log.Error("%s (%s) > %s\n\n%s", core.Bold(name), core.Dim(fmt.Sprintf("%s:%d", shell.Address, shell.Port)), core.Red(runCommand), out)
		} else {
			log.Info("%s (%s) > %s\n\n%s", core.Bold(name), core.Dim(fmt.Sprintf("%s:%d", shell.Address, shell.Port)), core.Blue(runCommand), out)
		}
	}
}
