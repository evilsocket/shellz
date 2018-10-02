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
	outLock = sync.Mutex{}
	wq      = queue.New(-1, cmdWorker)
)

func processOutput(out []byte, shell models.Shell) string {
	outLock.Lock()
	defer outLock.Unlock()

	outs := core.Dim(" <no output>")
	if out != nil {
		fileName := toOutput
		if fileName == "" {
			outs = fmt.Sprintf("\n\n%s\n", core.Trim(string(out)))
		} else {
			size := len(out)
			buff := bytes.Buffer{}
			if tmpl, err := template.New("filename").Parse(toOutput); err != nil {
				panic(err)
			} else if err := tmpl.Execute(&buff, shell); err != nil {
				panic(err)
			} else {
				fileName = buff.String()
			}

			if file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644); err != nil {
				outs = fmt.Sprintf(" > error while saving to %s: %s", fileName, err)
			} else {
				defer file.Close()
				if wrote, err := file.Write(out); err != nil {
					outs = fmt.Sprintf(" > error while saving to %s: %s", fileName, err)
				} else if wrote != size {
					outs = fmt.Sprintf(" > error while saving to %s: size is %d, wrote %d", fileName, size, wrote)
				} else {
					outs = core.Dim(fmt.Sprintf(" > %d bytes saved to %s", size, fileName))
				}
			}
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

func cmdWorker(job queue.Job) {
	start := time.Now()
	shell := job.(models.Shell)
	name := shell.Name
	err, session := shell.NewSession(timeouts)
	if err != nil {
		if doTest {
			onTestFail(shell, err)
		} else {
			log.Warning("error while creating session for shell %s: %s", name, err)
		}
		return
	}
	defer session.Close()

	out, err := session.Exec(command)
	if doTest {
		if err != nil {
			onTestFail(shell, err)
		}
	} else {
		took := core.Dim(time.Since(start).String())
		host := core.Dim(fmt.Sprintf("%s@%s:%d", shell.Identity.Username, shell.Host, shell.Port))
		outs := processOutput(out, shell)

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
	}
}

func runCommand() {
	if err, on = doShellSelection(onFilter, false); err != nil {
		log.Fatal("%s", err)
	} else if len(on) == 0 {
		log.Fatal("no enabled shell selected by the filter %s", core.Dim(onFilter))
	}

	if doTest {
		log.Info("testing %d shells ...\n", len(on))
	} else {
		log.Info("running %s on %d shells ...\n", core.Dim(command), len(on))
	}

	for name := range on {
		wq.Add(on[name])
	}

	wq.WaitDone()
}
