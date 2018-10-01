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

func cmdWorker(job queue.Job) {
	start := time.Now()
	shell := job.(models.Shell)
	name := shell.Name
	err, session := shell.NewSession(timeouts)
	if err != nil {
		log.Warning("error while creating session for shell %s: %s", name, err)
		return
	}
	defer session.Close()

	out, err := session.Exec(command)

	took := core.Dim(time.Since(start).String())
	host := core.Dim(fmt.Sprintf("%s@%s:%d", shell.Identity.Username, shell.Host, shell.Port))
	outs := core.Dim(" <no output>")

	if out != nil {
		fileName := toOutput
		if fileName == "" {
			outs = core.Trim(string(out))
			outs = fmt.Sprintf("\n\n%s\n", outs)
		} else {
			buff := bytes.Buffer{}
			if tmpl, err := template.New("filename").Parse(toOutput); err != nil {
				panic(err)
			} else if err := tmpl.Execute(&buff, shell); err != nil {
				panic(err)
			} else {
				fileName = buff.String()
			}

			outLock.Lock()
			defer outLock.Unlock()

			size := len(out)
			file, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
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

	for name := range on {
		wq.Add(on[name])
	}

	wq.WaitDone()
}
