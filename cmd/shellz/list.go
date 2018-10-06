package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/plugins"
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

	keys := []string{}
	for k := range Idents {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		i := Idents[name]
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
	if plugins.Number() > 0 {
		rows := [][]string{}
		cols := []string{
			"Name",
			"Path",
		}

		plugins.Each(func(p *plugins.Plugin) {
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
		"Groups",
		"Type",
		"Host",
		"Port",
		"Identity",
		"Enabled",
		// "Path",
	}

	keys := []string{}
	for k := range Shells {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		sh := Shells[name]
		en := core.Green("✔")
		if !sh.Enabled {
			en = core.Red("✖")
		}
		row := []string{
			core.Bold(sh.Name),
			core.Blue(strings.Join(sh.Groups, ", ")),
			core.Dim(sh.Type),
			sh.Host,
			fmt.Sprintf("%d", sh.Port),
			core.Yellow(sh.IdentityName),
			en,
			// core.Dim(sh.Path),
		}

		if !sh.Enabled {
			for i := range row {
				row[i] = core.Dim(row[i])
			}
		}

		rows = append(rows, row)
	}

	fmt.Printf("\n%s\n", core.Bold("shells"))
	core.AsTable(os.Stdout, cols, rows)
}

func showList() {
	showIdentsList()
	showPluginsList()
	showShellsList()
}
