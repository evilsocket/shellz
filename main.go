package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/models"
	"github.com/evilsocket/shellz/queue"
)

var (
	command  = ""
	onFilter = "*"
	onNames  = []string{}
	on       = models.Shells{}
	doList   = false
	err      = error(nil)
	idents   = models.Identities(nil)
	shells   = models.Shells(nil)

	wq = queue.New(-1, func(job queue.Job) {
		start := time.Now()
		shell := job.(models.Shell)
		name := shell.Name
		err, session := shell.NewSession()
		if err != nil {
			log.Warning("error while creating session for shell %s: %s", name, err)
			return
		}
		defer session.Close()

		out, err := session.Exec(command)

		took := core.Dim(time.Since(start).String())
		host := core.Dim(fmt.Sprintf("%s %s:%d", shell.Identity.Username, shell.Address, shell.Port))
		outs := core.Dim("<no output>")
		if out != nil {
			outs = core.Trim(string(out))
			outs = fmt.Sprintf("\n\n%s\n", outs)
		}

		if err != nil {
			log.Error("%s (%s %s %s) > %s (%s)%s",
				core.Bold(name),
				core.Green(shell.Type),
				host,
				took,
				command,
				core.Red(err.Error()),
				outs)
		} else {
			log.Info("%s (%s %s %s) > %s%s",
				core.Bold(name),
				core.Green(shell.Type),
				host,
				took,
				core.Blue(command),
				outs)
		}
	})
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

func runCommand() {
	if onFilter == "*" {
		on = shells
	} else {
		for _, name := range core.CommaSplit(onFilter) {
			if shell, found := shells[name]; !found {
				log.Fatal("can't find shell %s", name)
			} else {
				on[name] = shell
			}
		}
	}

	if len(on) == 0 {
		log.Fatal("no shell selected by the filter %s", core.Dim(onFilter))
	}

	log.Info("running %s on %d shells ...\n", core.Dim(command), len(on))

	for name, _ := range on {
		wq.Add(on[name])
	}

	wq.WaitDone()
}

func showHelp() {
	log.Info("none of the -run or -list parameters have been specified")

	fmt.Println()
	fmt.Printf("Usage:\n\n")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Printf("Examples:\n\n")

	examples := []struct {
		cmd  string
		help string
	}{
		{"shellz -list", "list available identities and shells"},
		{"shellz -run id", "run the command 'id' on each shell"},
		{"shellz -run id -on machineA", "run the command 'id' on a single shell named 'machineA'"},
		{"shellz -run id -on 'machineA, machineB'", "run the command 'id' on machineA and machineB"},
	}

	for _, e := range examples {
		fmt.Printf("  %s\n", core.Dim("# "+e.help))
		fmt.Printf("  %s\n", core.Bold(e.cmd))
		fmt.Println()
	}

	os.Exit(1)
}

func init() {
	flag.StringVar(&command, "run", command, "Command to run on the selected shells.")
	flag.StringVar(&onFilter, "on", onFilter, "Comma separated list of shell names to select or * for all.")
	flag.BoolVar(&log.DebugMessages, "debug", log.DebugMessages, "Enable debug messages.")
	flag.BoolVar(&doList, "list", doList, "List available shells and exit.")
	flag.Parse()
}

func main() {
	log.Raw(core.Banner)

	err, idents, shells = models.Load()
	if err != nil {
		log.Fatal("error while loading identities and shells: %s", err)
	} else if len(shells) == 0 {
		log.Fatal("no shells found on the system, start creating json files inside %s", models.Paths["shells"])
	} else {
		log.Debug("loaded %d identities and %d shells", len(idents), len(shells))
	}

	if doList {
		showList()
	} else if command != "" {
		runCommand()
	} else {
		showHelp()
	}
}
