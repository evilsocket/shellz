package session

import (
	"net"
	"strconv"
	"sync"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/models"

	"github.com/evilsocket/islazy/async"
)

const (
	SSHAgentKey = "@agent"
	SSHAuthSock = "SSH_AUTH_SOCK"
)

type SSHSession struct {
	sync.Mutex
	host     string
	proxy    models.Proxy
	config   *ssh.ClientConfig
	client   *ssh.Client
	session  *ssh.Session
	timeouts core.Timeouts
}

func NewSSH(sh models.Shell, timeouts core.Timeouts) (error, Session) {
	err, cfg := sh2ClientConfig(sh, timeouts)
	if err != nil {
		return err, nil
	}

	sshs := &SSHSession{
		host:     net.JoinHostPort(sh.Host, strconv.Itoa(sh.Port)),
		config:   cfg,
		proxy:    sh.Proxy,
		timeouts: timeouts,
	}

	err, _ = async.WithTimeout(sshs.timeouts.Connect, func() interface{} {
		if sshs.proxy.Empty() {
			log.Debug("dialing ssh %s ...", sshs.host)
			if sshs.client, err = ssh.Dial("tcp", sshs.host, sshs.config); err != nil {
				return err
			}
		} else {
			log.Debug("dialing ssh %s via socks5://%s ...", sshs.host, sshs.proxy.String())

			if dialer, err := proxy.SOCKS5("tcp", sshs.proxy.String(), nil, proxy.Direct); err != nil {
				return err
			} else if conn, err := dialer.Dial("tcp", sshs.host); err != nil {
				return err
			} else if c, chans, reqs, err := ssh.NewClientConn(conn, sshs.host, sshs.config); err != nil {
				conn.Close()
				return err
			} else {
				sshs.client = ssh.NewClient(c, chans, reqs)
			}
		}
		return nil
	})

	if err != nil {
		return err, nil
	} else if sshs.session, err = sshs.client.NewSession(); err != nil {
		sshs.client.Close()
		return err, nil
	}

	return nil, sshs
}

func (s *SSHSession) Type() string {
	return "ssh"
}

type cmdResult struct {
	out []byte
	err error
}

func (s *SSHSession) Exec(cmd string) ([]byte, error) {
	s.Lock()
	defer s.Unlock()

	err, obj := async.WithTimeout(s.timeouts.Write+s.timeouts.Read, func() interface{} {
		out, err := s.session.CombinedOutput(cmd)
		return cmdResult{out: out, err: err}
	})

	if err != nil {
		return nil, err
	}

	res := obj.(cmdResult)
	return res.out, res.err
}

func (s *SSHSession) Close() {
	s.Lock()
	defer s.Unlock()
	s.session.Close()
	s.client.Close()
}
