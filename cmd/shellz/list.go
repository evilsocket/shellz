package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/evilsocket/shellz/core"
)

func showIdentsList() {
	rows := [][]string{}
	cols := []string{
		"Name",
		"Username",
		"Key",
		"Password",
		// "Path",
	}

	for _, i := range idents {
		key := i.KeyFile
		if key == "" {
			key = core.Dim("<empty>")
		}
		pass := strings.Repeat("*", len(i.Password))
		if pass == "" {
			pass = core.Dim("<empty>")
		}

		rows = append(rows, []string{
			core.Bold(i.Name),
			i.Username,
			key,
			pass,
			// core.Dim(i.Path),
		})
	}

	fmt.Printf("\n%s\n", core.Bold("identities"))
	core.AsTable(os.Stdout, cols, rows)
}

func showShellsList() {
	rows := [][]string{}
	cols := []string{
		"Name",
		"Type",
		"Host",
		"Address",
		"Port",
		"Identity",
		// "Path",
	}

	for _, sh := range shells {
		rows = append(rows, []string{
			core.Bold(sh.Name),
			core.Dim(sh.Type),
			sh.Host,
			core.Blue(sh.Address.String()),
			fmt.Sprintf("%d", sh.Port),
			core.Yellow(sh.IdentityName),
			// core.Dim(sh.Path),
		})
	}

	fmt.Printf("\n%s\n", core.Bold("shells"))
	core.AsTable(os.Stdout, cols, rows)
}

func showList() {
	showIdentsList()
	showShellsList()
}
