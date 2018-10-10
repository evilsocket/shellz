package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"text/template"
	"time"

	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/models"
	"github.com/evilsocket/shellz/plugins"
	"github.com/evilsocket/shellz/session"

	"github.com/dustin/go-humanize"

	"github.com/evilsocket/islazy/async"
	"github.com/evilsocket/islazy/str"
	"github.com/evilsocket/islazy/tui"
)

type statistics struct {
	Started       time.Time
	Done          time.Time
	Shells        uint64
	Success       uint64
	Failed        uint64
	FailedConnect uint64
	FailedExec    uint64
	Output        uint64
	AvgTime       uint64
}

var (
	toOutputLock = sync.Mutex{}
	wq           = async.NewQueue(-1, cmdWorker)
	stats        = statistics{}
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
			outs = tui.Dim(fmt.Sprintf(" > %d bytes saved to %s", size, fileName))
		}
	}
	return
}

func processOutput(out []byte, shell models.Shell) string {
	outs := tui.Dim(" <no output>")
	if out != nil {
		if toOutput == "" {
			outs = fmt.Sprintf("\n\n%s\n", str.Trim(string(out)))
		} else {
			outs = toOutputFile(shell, out)
		}
	}
	return outs
}

func onTestFail(sh models.Shell, err error) {
	log.Warning("shell %s failed with: %s", tui.Bold(sh.Name), err)
	sh.Enabled = false
	if err := sh.Save(); err != nil {
		log.Error("error while disabling shell %s: %s", sh.Name, err)
	}
}

func onTestSuccess(sh models.Shell) {
	log.Info("shell %s is working", tui.Bold(sh.Name))
	if !sh.Enabled {
		log.Debug("enabling shell %s", sh.Name)
		sh.Enabled = true
		if err := sh.Save(); err != nil {
			log.Error("error while enabling shell %s: %s", sh.Name, err)
		}
	}
}

func trackSuccess() {
	atomic.AddUint64(&stats.Success, 1)
}

func trackFailure(connect bool) {
	atomic.AddUint64(&stats.Failed, 1)
	if connect {
		atomic.AddUint64(&stats.FailedConnect, 1)
	} else {
		atomic.AddUint64(&stats.FailedExec, 1)
	}
}

func trackOutput(out []byte) {
	if out != nil {
		atomic.AddUint64(&stats.Output, uint64(len(out)))
	}
}

func findSessionFor(sh models.Shell) (err error, sess session.Session) {
	// first try one of the default handlers
	if err, sess = session.For(sh, timeouts); sess == nil && err == nil {
		// try one of the user plugins
		if plugin := plugins.Get(sh); plugin != nil {
			err, sess = plugin.NewSession(sh, timeouts)
		}
	}
	// no error but no session found?
	if err == nil && sess == nil {
		err = fmt.Errorf("session type %s for shell %s is not supported", sh.Type, sh.Name)
	}
	return
}

func cmdWorker(job async.Job) {
	shell := job.(models.Shell)
	start := time.Now()

	err, session := findSessionFor(shell)
	if err != nil {
		trackFailure(true)
		if doTest {
			onTestFail(shell, err)
		} else {
			log.Warning("error while creating session for shell %s: %s", shell.Name, err)
		}
		return
	}
	defer session.Close()

	out, err := session.Exec(command)
	took := tui.Dim(time.Since(start).String())
	trackOutput(out)

	if doTest {
		if err != nil {
			trackFailure(false)
			onTestFail(shell, err)
		} else {
			trackSuccess()
			onTestSuccess(shell)
		}
	} else {
		outs := processOutput(out, shell)
		host := ""
		if shell.Identity.Username != "" {
			host = tui.Dim(fmt.Sprintf("%s@%s", shell.Identity.Username, shell.Host))
		} else {
			host = tui.Dim(shell.Host)
		}

		if !shell.Proxy.Empty() {
			host = tui.Dim(fmt.Sprintf("%s:%d > %s", shell.Proxy.Address, shell.Proxy.Port, host))
		}

		if err != nil {
			trackFailure(false)
			log.Error("%s (%s %s %s) > %s (%s)%s",
				tui.Bold(shell.Name),
				tui.Green(shell.Type),
				host,
				took,
				command,
				tui.Red(err.Error()),
				outs)
		} else {
			trackSuccess()
			log.Info("%s (%s %s %s) > %s%s",
				tui.Bold(shell.Name),
				tui.Green(shell.Type),
				host,
				took,
				tui.Blue(command),
				outs)
		}
	}
}

func viewStats() {
	log.Raw(tui.Dim("_______________________"))
	log.Raw(tui.Bold("Statistics\n"))

	totTime := stats.Done.Sub(stats.Started)
	avgTime := time.Duration(0)
	if stats.Success > 0 {
		avgTime = time.Duration(uint64(totTime) / stats.Success)
	}

	log.Raw("total shells : %d", stats.Shells)
	log.Raw("total time   : %s (%s/shell avg)", totTime, avgTime)
	log.Raw("total output : %s", humanize.Bytes(stats.Output))
	log.Raw(tui.Green("ok           : %d"), stats.Success)
	if stats.Failed > 0 {
		log.Raw(tui.Red("ko           : %d ( %d connect / %d exec )"), stats.Failed, stats.FailedConnect, stats.FailedExec)
	} else {
		log.Raw(tui.Dim("ko           : 0"))
	}
}

func runCommand() {
	log.Debug("onFilter = %s", onFilter)
	if err, onShells = doShellSelection(onFilter, doForce); err != nil {
		log.Fatal("%s", err)
	} else if nShells = len(onShells); nShells == 0 {
		log.Fatal("no enabled shell selected by the filter %s (use the -force argument to select disabled shells)", tui.Dim(onFilter))
	}

	if doTest {
		log.Debug("testing %d shells ...", nShells)
	} else {
		log.Debug("running %s on %d shells ...", tui.Dim(command), nShells)
	}

	stats.Shells = uint64(nShells)
	stats.Started = time.Now()

	for name := range onShells {
		wq.Add(onShells[name])
	}

	wq.WaitDone()

	stats.Done = time.Now()
	if doStats {
		viewStats()
	}
}
