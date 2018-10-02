package main

import (
	"fmt"
	"strings"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/models"
)

func findShells(name string) models.Shells {
	found := models.Shells{}
	for _, sh := range shells {
		if sh.Name == name || strings.HasPrefix(sh.Name, name) {
			found[sh.Name] = sh
		}
	}
	return found
}

func doEnabledSelection(m models.Shells, includeDisabled bool) models.Shells {
	sel := models.Shells{}
	for _, sh := range m {
		if includeDisabled || sh.Enabled {
			sel[sh.Name] = sh
		} else {
			log.Debug("skipping disabled shell %s", sh.Name)
		}
	}
	return sel
}

func doShellSelection(filter string, includeDisabled bool) (error, models.Shells) {
	sel := models.Shells{}

	if filter == "*" {
		return nil, doEnabledSelection(shells, includeDisabled)
	}

	for _, name := range core.CommaSplit(filter) {
		if found := findShells(name); len(found) == 0 {
			return fmt.Errorf("can't find shell %s", name), nil
		} else {
			sel = doEnabledSelection(found, includeDisabled)
		}
	}

	return nil, sel
}
