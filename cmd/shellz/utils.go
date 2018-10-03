package main

import (
	"strings"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/models"
)

func doFilterSelection(expr string) models.Shells {
	// first match by group, full match
	for group, shells := range Groups {
		if expr == group {
			return shells
		}
	}

	// then by shell name, full or prefix
	found := models.Shells{}
	for _, sh := range Shells {
		if sh.Name == expr || strings.HasPrefix(sh.Name, expr) {
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

func doShellSelection(csFilters string, includeDisabled bool) (error, models.Shells) {
	sel := models.Shells{}
	for _, filter := range core.CommaSplit(csFilters) {
		found := doFilterSelection(filter)
		for name, shell := range doEnabledSelection(found, includeDisabled) {
			sel[name] = shell
		}
	}

	return nil, sel
}
