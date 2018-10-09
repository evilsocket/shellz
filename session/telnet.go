package session

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/models"

	"github.com/reiver/go-telnet"
)

type TelnetSession struct {
	sync.Mutex
	host     string
	client   *telnet.Conn
	timeouts core.Timeouts
}

func NewTelnet(sh models.Shell, timeouts core.Timeouts) (error, Session) {
	var err error

	t := &TelnetSession{
		host:     net.JoinHostPort(sh.Host, strconv.Itoa(sh.Port)),
		timeouts: timeouts,
	}

	err, _ = core.WithTimeout(timeouts.Connect, func() interface{} {
		t.client, err = telnet.DialTo(t.host)
		return err
	})
	if err != nil {
		return err, nil
	}

	if sh.Identity.Username != "" && sh.Identity.Password != "" {
		t.doReadUntil(": ")
		if _, err = t.doWrite([]byte(sh.Identity.Username + "\n")); err != nil {
			return fmt.Errorf("error while sending telnet username: %s", err), nil
		}

		t.doReadUntil(": ")
		if _, err = t.doWrite([]byte(sh.Identity.Password + "\n")); err != nil {
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
	err, obj := core.WithTimeout(t.timeouts.Read, func() interface{} {
		n, err := t.client.Read(buf)
		return rw{e: err, n: n}
	})
	if err != nil {
		return -1, err
	}

	r := obj.(rw)
	return r.n, r.e
}

func (t *TelnetSession) doWrite(buf []byte) (int, error) {
	err, obj := core.WithTimeout(t.timeouts.Write, func() interface{} {
		n, err := t.client.Write(buf)
		return rw{e: err, n: n}
	})
	if err != nil {
		return -1, err
	}

	w := obj.(rw)
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
		// log.Debug("  read 0x%x %c", b[0], b[0])
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
