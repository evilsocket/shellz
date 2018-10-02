package main

import (
	"fmt"
	"strings"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/models"
)

func findShells(name string) []models.Shell {
	found := []models.Shell{}
	for _, sh := range shells {
		if sh.Name == name || strings.HasPrefix(sh.Name, name) {
			found = append(found, sh)
		}
	}
	return found
}

func doShellSelection(filter string, includeDisabled bool) (error, models.Shells) {
	sel := models.Shells{}

	if filter == "*" {
		for _, sh := range shells {
			if includeDisabled || sh.Enabled {
				sel[sh.Name] = sh
			} else {
				log.Debug("skipping disabled shell %s", sh.Name)
			}
		}
		return nil, sel
	}

	for _, name := range core.CommaSplit(filter) {
		if found := findShells(name); len(found) == 0 {
			return fmt.Errorf("can't find shell %s", name), nil
		} else {
			for _, sh := range found {
				if includeDisabled || sh.Enabled {
					sel[sh.Name] = sh
				} else {
					log.Debug("skipping disabled shell %s", sh.Name)
				}
			}
		}
	}

	return nil, sel
}
