package session

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"

	"github.com/reiver/go-telnet"
)

const (
	bufferSize = 0xffff
)

type TelnetSession struct {
	sync.Mutex
	client *telnet.Conn
	buffer []byte
}

func NewTelnet(address net.IP, port int, user string, pass string, keyFile string) (error, Session) {
	host := fmt.Sprintf("%s:%d", address.String(), port)
	cli, err := telnet.DialTo(host)
	if err != nil {
		return err, nil
	}

	t := &TelnetSession{
		client: cli,
		buffer: make([]byte, bufferSize),
	}

	if user != "" && pass != "" {
		t.doReadUntil(": ")
		if _, err = t.client.Write([]byte(user + "\n")); err != nil {
			return fmt.Errorf("error while sending telnet username: %s", err), nil
		}

		t.doReadUntil(": ")
		if _, err = t.client.Write([]byte(pass + "\n")); err != nil {
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
