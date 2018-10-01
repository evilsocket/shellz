package session

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"

	"github.com/reiver/go-telnet"
)

type TelnetSession struct {
	sync.Mutex
	host     string
	client   *telnet.Conn
	timeouts Timeouts
}

func NewTelnet(ctx Context) (error, Session) {
	var err error

	t := &TelnetSession{
		host:     fmt.Sprintf("%s:%d", ctx.Address.String(), ctx.Port),
		timeouts: ctx.Timeouts,
	}

	done := make(chan error)
	timeout := time.After(ctx.Timeouts.Connect)
	go func() {
		t.client, err = telnet.DialTo(t.host)
		done <- err
	}()

	select {
	case <-timeout:
		return fmt.Errorf("timeout while dialing %s", t.host), nil
	case err := <-done:
		if err != nil {
			return err, nil
		}
	}

	if ctx.Username != "" && ctx.Password != "" {
		t.doReadUntil(": ")
		if _, err = t.doWrite([]byte(ctx.Username + "\n")); err != nil {
			return fmt.Errorf("error while sending telnet username: %s", err), nil
		}

		t.doReadUntil(": ")
		if _, err = t.doWrite([]byte(ctx.Password + "\n")); err != nil {
			return fmt.Errorf("error while sending telnet password: %s", err), nil
		}
	}

	return nil, t
}

func (t *TelnetSession) Type() string {
	return "telnet"
}

type rw struct {
	e error
	n int
}

func (t *TelnetSession) doRead(buf []byte) (int, error) {
	r := rw{}
	done := make(chan rw)
	timeout := time.After(t.timeouts.Read)
	go func() {
		n, err := t.client.Read(buf)
		done <- rw{e: err, n: n}
	}()

	select {
	case <-timeout:
		return 0, fmt.Errorf("timeout while reading from %s", t.host)
	case r = <-done:
		if r.e != nil {
			return 0, r.e
		}
	}

	return r.n, r.e
}

func (t *TelnetSession) doWrite(buf []byte) (int, error) {
	w := rw{}
	done := make(chan rw)
	timeout := time.After(t.timeouts.Write)
	go func() {
		n, err := t.client.Write(buf)
		done <- rw{e: err, n: n}
	}()

	select {
	case <-timeout:
		return 0, fmt.Errorf("timeout while writing to %s", t.host)
	case w = <-done:
		if w.e != nil {
			return 0, w.e
		}
	}

	return w.n, w.e
}

func (t *TelnetSession) doReadUntil(s string) (error, string) {
	log.Debug("doReadUntil(%s)", s)
	buff := ""
	for buff = ""; !strings.Contains(buff, s); {
		b := []byte{0}
		if _, err := t.doRead(b); err != nil {
			return err, ""
		}
		log.Debug("  read 0x%x %c", b[0], b[0])
		buff += string(b[0])
	}
	log.Debug("  => '%s'", buff)
	return nil, buff
}

func (t *TelnetSession) Exec(cmd string) ([]byte, error) {
	t.Lock()
	defer t.Unlock()

	cmd = fmt.Sprintf("(%s || echo) && echo PLACEHOLDER", cmd)
	if _, err := t.doWrite([]byte(cmd + "\n")); err != nil {
		return nil, fmt.Errorf("error while sending telnet command: %s", err)
	}
	t.doReadUntil(cmd)

	if err, s := t.doReadUntil("PLACEHOLDER"); err != nil {
		return nil, err
	} else {
		s = strings.Replace(s, "PLACEHOLDER", "", -1)
		s = core.TrimLeft(s)
		return []byte(s), nil
	}
}

func (t *TelnetSession) Close() {
	t.Lock()
	defer t.Unlock()
	t.client.Close()
}
