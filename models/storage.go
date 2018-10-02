package models

import (
	"fmt"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
)

type Identities map[string]Identity
type Shells map[string]Shell

func Load() (error, Identities, Shells) {
	idents := make(Identities)
	shells := make(Shells)

	log.Debug("loading identities from %s ...", Paths["idents"])
	err := core.Glob(Paths["idents"], "*.json", func(fileName string) error {
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
		return err, nil, nil
	}

	log.Debug("loading shells from %s ...", Paths["shells"])
	err = core.Glob(Paths["shells"], "*.json", func(fileName string) error {
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
		return err, nil, nil
	}

	return nil, idents, shells
}
