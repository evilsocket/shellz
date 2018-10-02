package main

import (
	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/models"
)

func runEnable(filter string, enable bool) {
	word := "enabled"
	if !enable {
		word = "disabled"
	}

	list := models.Shells{}
	if filter == "*" {
		list = shells
	} else {
		for _, name := range core.CommaSplit(filter) {
			if shell, found := shells[name]; !found {
				log.Fatal("can't find shell %s", name)
			} else {
				list[name] = shell
			}
		}
	}

	for _, shell := range list {
		if shell.Enabled == enable {
			log.Error("shell %s is already %s", shell.Name, word)
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
