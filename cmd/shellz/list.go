package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/plugins"

	"github.com/evilsocket/islazy/tui"
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
			key = tui.Dim("<empty>")
		}
		pass := strings.Repeat("*", len(i.Password))
		if pass == "" {
			pass = tui.Dim("<empty>")
		}

		rows = append(rows, []string{
			tui.Bold(i.Name),
			i.Username,
			key,
			pass,
			// tui.Dim(i.Path),
		})
	}

	fmt.Printf("\n%s\n", tui.Bold("identities"))
	tui.Table(os.Stdout, cols, rows)
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
				tui.Bold(p.Name),
				tui.Dim(p.Path),
			})
		})

		fmt.Printf("\n%s\n", tui.Bold("plugins"))
		tui.Table(os.Stdout, cols, rows)
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

	if err, onShells = doShellSelection(onFilter, true); err != nil {
		log.Fatal("%s", err)
	} else if nShells = len(onShells); nShells == 0 {
		log.Fatal("no shell selected by the filter %s", tui.Dim(onFilter))
	}

	keys := []string{}
	for k := range onShells {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		sh := onShells[name]
		en := tui.Green("✔")
		if !sh.Enabled {
			en = tui.Red("✖")
		}
		row := []string{
			tui.Bold(sh.Name),
			tui.Blue(strings.Join(sh.Groups, ", ")),
			tui.Dim(sh.Type),
			sh.Host,
			fmt.Sprintf("%d", sh.Port),
			tui.Yellow(sh.IdentityName),
			en,
			// tui.Dim(sh.Path),
		}

		if !sh.Enabled {
			for i := range row {
				row[i] = tui.Dim(row[i])
			}
		}

		rows = append(rows, row)
	}

	fmt.Printf("\n%s\n", tui.Bold("shells"))
	tui.Table(os.Stdout, cols, rows)
}

func showList() {
	showIdentsList()
	showPluginsList()
	showShellsList()
}
