package session

import (
	"io"
	"net"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/models"

	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/islazy/tui"
)

type SSHTun struct {
	*SSHSession
	listener net.Listener
	tunnel   models.Tunnel
}

func NewSSHTun(sh models.Shell, timeouts core.Timeouts) (error, Session) {
	err, s := NewSSH(sh, timeouts)
	if err != nil {
		return err, nil
	}

	tun := &SSHTun{
		SSHSession: s.(*SSHSession),
		tunnel:     sh.Tunnel,
	}

	return tun.start(), tun
}

func (t *SSHTun) Type() string {
	return "ssh.tunnel"
}

func (t *SSHTun) start() (err error) {
	log.Info("tunnel for %s (via %s) available at %s ...",
		tui.Bold(t.tunnel.Remote.String()),
		tui.Dim(t.host),
		tui.Bold(t.tunnel.Local.String()))

	if t.listener, err = net.Listen("tcp", t.tunnel.Local.String()); err != nil {
		return
	}
	defer t.listener.Close()

	for {
		if conn, err := t.listener.Accept(); err != nil {
			log.Error("error while accepting connection: %s", err)
			continue
		} else {
			go t.forward(conn)
		}
	}
}

func (t *SSHTun) pipe(writer, reader net.Conn) {
	if _, err := io.Copy(writer, reader); err != nil {
		log.Warning("pipe error: %s", err)
	}
}

func (t *SSHTun) forward(loc net.Conn) {
	log.Debug("forwarding %s -> %s", loc.LocalAddr().String(), loc.RemoteAddr().String())
	if rem, err := t.client.Dial("tcp", t.tunnel.Remote.String()); err != nil {
		log.Error("error while dialing %s: %s", t.tunnel.Remote.String(), err)
	} else {
		go t.pipe(loc, rem)
		go t.pipe(rem, loc)
	}
}

func (t *SSHTun) Close() {
	t.Lock()
	defer t.Unlock()
	t.listener.Close()
	t.session.Close()
	t.client.Close()
}
