package main

import (
	"github.com/evilsocket/shellz/log"
)

func runEnable(name string, enable bool) {
	word := "enabled"
	if !enable {
		word = "disabled"
	}

	if shell, found := shells[name]; !found {
		log.Fatal("can't find shell %s", name)
	} else if shell.Enabled == enable {
		log.Fatal("shell %s is already %s", name, word)
	} else {
		shell.Enabled = enable
		if err := shell.Save(); err != nil {
			log.Fatal("error while setting shell %s to %s: %s", name, word, err)
		} else {
			log.Info("shell %s succesfully %s", name, word)
		}
	}
}
