package main

import (
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shellz/models"
	"github.com/evilsocket/shellz/session"

	"github.com/evilsocket/islazy/async"
	"github.com/evilsocket/islazy/tui"
)

var (
	tunQueue = (*async.WorkQueue)(nil)
)

func tunWorker(job async.Job) {
	shell := job.(models.Shell)

	shell.Type = "ssh.tunnel"
	err, tun := session.For(shell, timeouts)
	if err != nil {
		log.Error("error while tunneling %s via shell %s: %s", shell.Tunnel.String(), shell.Name, err)
		return
	}
	defer tun.Close()
}

func runTunnel() {
	log.Debug("onFilter = %s", onFilter)
	if err, onShells = doShellSelection(onFilter, doForce); err != nil {
		log.Fatal("%s", err)
	} else if nShells = len(onShells); nShells == 0 {
		log.Fatal("no enabled shell selected by the filter %s (use the -force argument to select disabled shells)", tui.Dim(onFilter))
	}

	tmp := models.Shells{}
	for name, sh := range onShells {
		if sh.Type != "ssh" {
			log.Error("shell %s is of type %s which does not support tunneling, skipping.", sh.Name, sh.Type)
			continue
		} else if sh.Tunnel.Empty() {
			log.Error("shell %s does not provide tunneling information, skipping.", sh.Name)
			continue
		}
		tmp[name] = sh
	}
	onShells = tmp
	nShells = len(onShells)

	if nShells == 0 || nShells > 1 {
		log.Info("starting %d tunnels, please wait ...\n", nShells)
	} else {
		log.Info("starting %d tunnel, please wait ...\n", nShells)
	}

	tunQueue = async.NewQueue(numWorkers, tunWorker)

	for name := range onShells {
		tunQueue.Add(onShells[name])
	}

	tunQueue.WaitDone()
}
