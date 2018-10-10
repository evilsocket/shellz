package models

import (
	"fmt"

	"github.com/evilsocket/islazy/log"

	"github.com/evilsocket/islazy/fs"
)

type Identities map[string]Identity
type Shells map[string]Shell
type Groups map[string]Shells

func Load() (error, Identities, Shells, Groups) {
	idents := make(Identities)
	shells := make(Shells)
	groups := make(Groups)

	log.Debug("loading identities from %s ...", Paths["idents"])
	err := fs.Glob(Paths["idents"], "*.json", func(fileName string) error {
		if err, ident := LoadIdent(fileName); err != nil {
			return fmt.Errorf("error while loading identity '%s': %s", fileName, err)
		} else if taken, found := idents[ident.Name]; found {
			return fmt.Errorf("identity '%s' has name %s which is already taken by '%s'", fileName, ident.Name, taken.Path)
		} else {
			idents[ident.Name] = ident
		}
		return nil
	})
	if err != nil {
		return err, nil, nil, nil
	}

	log.Debug("loading shells from %s ...", Paths["shells"])
	err = fs.Glob(Paths["shells"], "*.json", func(fileName string) error {
		if err, shell := LoadShell(fileName, idents); err != nil {
			return fmt.Errorf("error while loading shell '%s': %s", fileName, err)
		} else if taken, found := shells[shell.Name]; found {
			return fmt.Errorf("shell '%s' has name %s which is already taken by '%s'", fileName, shell.Name, taken.Path)
		} else {
			shells[shell.Name] = shell
		}
		return nil
	})
	if err != nil {
		return err, nil, nil, nil
	}

	log.Debug("creating groups ...")
	for _, sh := range shells {
		for _, group := range append(sh.Groups, "all") {
			ref, found := groups[group]
			if found == false {
				ref = make(Shells)
				groups[group] = ref
			}
			ref[sh.Name] = sh
		}
	}

	return nil, idents, shells, groups
}
