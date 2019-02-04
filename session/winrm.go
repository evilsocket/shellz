package session

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/models"

	"github.com/evilsocket/islazy/async"

	"github.com/masterzen/winrm"
)

type WinRMSession struct {
	sync.Mutex
	host     string
	port     int
	endpoint *winrm.Endpoint
	client   *winrm.Client
	timeouts core.Timeouts
}

func NewWinRM(sh models.Shell, timeouts core.Timeouts) (error, Session) {
	var err error

	t := &WinRMSession{
		host:     sh.Host,
		port:     sh.Port,
		timeouts: timeouts,
		endpoint: winrm.NewEndpoint(
			sh.Host,
			sh.Port,
			sh.HTTPS,
			sh.Insecure,
			nil,
			nil,
			nil,
			timeouts.Total()),
	}

	_, err = async.WithTimeout(timeouts.Connect, func() interface{} {
		t.client, err = winrm.NewClient(t.endpoint, sh.Identity.Username, sh.Identity.Password)
		return err
	})
	if err != nil {
		return err, nil
	}

	return nil, t
}

func (w *WinRMSession) Type() string {
	return "winrm"
}

func (w *WinRMSession) Exec(cmd string) ([]byte, error) {
	w.Lock()
	defer w.Unlock()

	obj, err := async.WithTimeout(w.timeouts.RW(), func() interface{} {
		outWriter := bytes.Buffer{}
		errWriter := bytes.Buffer{}
		if _, err := w.client.Run(cmd, &outWriter, &errWriter); err != nil {
			return cmdResult{err: err}
		} else {
			return cmdResult{out: outWriter.Bytes(), err: fmt.Errorf("%s", errWriter.String())}
		}
	})
	if err != nil {
		return nil, err
	}

	res := obj.(cmdResult)
	return res.out, res.err
}

func (w *WinRMSession) Close() {
	w.Lock()
	defer w.Unlock()
}
