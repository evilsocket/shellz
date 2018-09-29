package models

import (
	"fmt"
	"path/filepath"

	"github.com/evilsocket/shellz/log"
)

type Identities map[string]Identity
type Shells map[string]Shell

func Load() (error, Identities, Shells) {
	idents := make(Identities)
	shells := make(Shells)

	log.Info("loading identities from %s ...", Paths["idents"])

	if files, err := filepath.Glob(filepath.Join(Paths["idents"], "*.json")); err != nil {
		return err, nil, nil
	} else {
		for _, fileName := range files {
			if err, ident := LoadIdent(fileName); err != nil {
				return fmt.Errorf("error while loading identity '%s': %s", fileName, err), nil, nil
			} else if taken, found := idents[ident.Name]; found {
				return fmt.Errorf("identity '%s' has name %s which is already taken by '%s'", fileName, ident.Name, taken.Path), nil, nil
			} else {
				idents[ident.Name] = ident
			}
		}
	}

	log.Info("loading shells from %s ...", Paths["shells"])

	if files, err := filepath.Glob(filepath.Join(Paths["shells"], "*.json")); err != nil {
		return err, nil, nil
	} else {
		for _, fileName := range files {
			if err, shell := LoadShell(fileName, idents); err != nil {
				return fmt.Errorf("error while loading shell '%s': %s", fileName, err), nil, nil
			} else if taken, found := shells[shell.Name]; found {
				return fmt.Errorf("shell '%s' has name %s which is already taken by '%s'", fileName, shell.Name, taken.Path), nil, nil
			} else {
				shells[shell.Name] = shell
			}
		}
	}

	return nil, idents, shells
}
