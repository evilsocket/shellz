package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/models"
	"github.com/evilsocket/shellz/queue"
)

var (
	toOutputLock = sync.Mutex{}
	wq           = queue.New(-1, cmdWorker)
)

func toOutputFilename(shell models.Shell) string {
	buff := bytes.Buffer{}
	if tmpl, err := template.New("filename").Parse(toOutput); err != nil {
		log.Fatal("error while parsing '%s': %s", toOutput, err)
	} else if err := tmpl.Execute(&buff, shell); err != nil {
		log.Fatal("error while running '%s' on shell %s: %s", toOutput, shell.Name, err)
	}
	return buff.String()
}

func toOutputFile(shell models.Shell, out []byte) (outs string) {
	toOutputLock.Lock()
	defer toOutputLock.Unlock()

	fileName := toOutputFilename(shell)
	if file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644); err != nil {
		outs = fmt.Sprintf(" > error while saving to %s: %s", fileName, err)
	} else {
		defer file.Close()
		size := len(out)
		if wrote, err := file.Write(out); err != nil {
			outs = fmt.Sprintf(" > error while saving to %s: %s", fileName, err)
		} else if wrote != size {
			outs = fmt.Sprintf(" > error while saving to %s: size is %d, wrote %d", fileName, size, wrote)
		} else {
			outs = core.Dim(fmt.Sprintf(" > %d bytes saved to %s", size, fileName))
		}
	}
	return
}

func processOutput(out []byte, shell models.Shell) string {
	outs := core.Dim(" <no output>")
	if out != nil {
		if toOutput == "" {
			outs = fmt.Sprintf("\n\n%s\n", core.Trim(string(out)))
		} else {
			outs = toOutputFile(shell, out)
		}
	}
	return outs
}

func onTestFail(sh models.Shell, err error) {
	log.Warning("shell %s failed with: %s", core.Bold(sh.Name), err)
	sh.Enabled = false
	if err := sh.Save(); err != nil {
		log.Error("error while disabling shell %s: %s", sh.Name, err)
	}
}

func onTestSuccess(sh models.Shell) {
	log.Info("shell %s is working", core.Bold(sh.Name))
	if !sh.Enabled {
		log.Debug("enabling shell %s", sh.Name)
		sh.Enabled = true
		if err := sh.Save(); err != nil {
			log.Error("error while enabling shell %s: %s", sh.Name, err)
		}
	}
}

func cmdWorker(job queue.Job) {
	start := time.Now()
	shell := job.(models.Shell)

	err, session := shell.NewSession(timeouts)
	if err != nil {
		if doTest {
			onTestFail(shell, err)
		} else {
			log.Warning("error while creating session for shell %s: %s", shell.Name, err)
		}
		return
	}
	defer session.Close()

	out, err := session.Exec(command)
	if doTest {
		if err != nil {
			onTestFail(shell, err)
		} else {
			onTestSuccess(shell)
		}
	} else {
		took := core.Dim(time.Since(start).String())
		outs := processOutput(out, shell)
		host := ""
		if shell.Identity.Username != "" {
			host = core.Dim(fmt.Sprintf("%s@%s", shell.Identity.Username, shell.Host))
		} else {
			host = core.Dim(shell.Host)
		}

		if err != nil {
			log.Error("%s (%s %s %s) > %s (%s)%s",
				core.Bold(shell.Name),
				core.Green(shell.Type),
				host,
				took,
				command,
				core.Red(err.Error()),
				outs)
		} else {
			log.Info("%s (%s %s %s) > %s%s",
				core.Bold(shell.Name),
				core.Green(shell.Type),
				host,
				took,
				core.Blue(command),
				outs)
		}
	}
}

func runCommand() {
	log.Debug("onFilter = %s", onFilter)
	if err, onShells = doShellSelection(onFilter, doForce); err != nil {
		log.Fatal("%s", err)
	} else if nShells = len(onShells); nShells == 0 {
		log.Fatal("no enabled shell selected by the filter %s (use the -force argument to select disabled shells)", core.Dim(onFilter))
	}

	if doTest {
		log.Debug("testing %d shells ...", nShells)
	} else {
		log.Debug("running %s on %d shells ...", core.Dim(command), nShells)
	}

	for name := range onShells {
		wq.Add(onShells[name])
	}

	wq.WaitDone()
}
