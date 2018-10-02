package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/session"
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

func showPluginsList() {
	if session.NumPlugins() > 0 {
		rows := [][]string{}
		cols := []string{
			"Name",
			"Path",
		}

		session.EachPlugin(func(p *session.Plugin) {
			rows = append(rows, []string{
				core.Bold(p.Name),
				core.Dim(p.Path),
			})
		})

		fmt.Printf("\n%s\n", core.Bold("plugins"))
		core.AsTable(os.Stdout, cols, rows)
	}
}

func showShellsList() {
	rows := [][]string{}
	cols := []string{
		"Name",
		"Type",
		"Host",
		"Port",
		"Identity",
		"Enabled",
		// "Path",
	}

	for _, sh := range shells {
		en := core.Green("✔")
		if !sh.Enabled {
			en = core.Red("✖")
		}
		rows = append(rows, []string{
			core.Bold(sh.Name),
			core.Dim(sh.Type),
			sh.Host,
			fmt.Sprintf("%d", sh.Port),
			core.Yellow(sh.IdentityName),
			en,
			// core.Dim(sh.Path),
		})
	}

	fmt.Printf("\n%s\n", core.Bold("shells"))
	core.AsTable(os.Stdout, cols, rows)
}

func showList() {
	showIdentsList()
	showPluginsList()
	showShellsList()
}
