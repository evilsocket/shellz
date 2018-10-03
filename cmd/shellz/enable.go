package main

import (
	"github.com/evilsocket/shellz/log"
)

func runEnable(filter string, enable bool) {
	if err, onShells = doShellSelection(filter, true); err != nil {
		log.Fatal("%s", err)
	} else if nShells = len(onShells); nShells == 0 {
		log.Fatal("no shell selected by the filter %s", filter)
	}

	word := "enabled"
	if !enable {
		word = "disabled"
	}
	for _, shell := range onShells {
		if shell.Enabled == enable {
			log.Debug("shell %s already %s", shell.Name, word)
		} else {
			shell.Enabled = enable
			if err := shell.Save(); err != nil {
				log.Error("error while setting shell %s to %s: %s", shell.Name, word, err)
			} else {
				log.Info("shell %s succesfully %s", shell.Name, word)
			}
		}
	}
}
