package session

import (
	"fmt"
	"strings"
	"sync"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"

	"github.com/reiver/go-telnet"
)

type TelnetSession struct {
	sync.Mutex
	host   string
	client *telnet.Conn
}

func NewTelnet(ctx Context) (error, Session) {
	var err error

	t := &TelnetSession{
		host: fmt.Sprintf("%s:%d", ctx.Address.String(), ctx.Port),
	}

	t.client, err = telnet.DialTo(t.host)
	if err != nil {
		return err, nil
	}

	if ctx.Username != "" && ctx.Password != "" {
		t.doReadUntil(": ")
		if _, err = t.client.Write([]byte(ctx.Username + "\n")); err != nil {
			return fmt.Errorf("error while sending telnet username: %s", err), nil
		}

		t.doReadUntil(": ")
		if _, err = t.client.Write([]byte(ctx.Password + "\n")); err != nil {
			return fmt.Errorf("error while sending telnet password: %s", err), nil
		}
	}

	return nil, t
}

func (t *TelnetSession) Type() string {
	return "telnet"
}

func (t *TelnetSession) doReadUntil(s string) (error, string) {
	log.Debug("doReadUntil(%s)", s)
	buff := ""
	for buff = ""; !strings.Contains(buff, s); {
		b := []byte{0}
		if _, err := t.client.Read(b); err != nil {
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
	if _, err := t.client.Write([]byte(cmd + "\n")); err != nil {
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
