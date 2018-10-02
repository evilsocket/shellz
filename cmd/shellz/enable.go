package main

import (
	"github.com/evilsocket/shellz/log"
)

func runEnable(filter string, enable bool) {
	err, list := doShellSelection(filter, true)
	if err != nil {
		log.Fatal("%s", err)
	}

	word := "enabled"
	if !enable {
		word = "disabled"
	}
	for _, shell := range list {
		if shell.Enabled == enable {
			log.Info("shell %s already %s", shell.Name, word)
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
