package session

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"

	"github.com/evilsocket/shellz/log"
)

const (
	SSHAgentKey = "@agent"
	SSHAuthSock = "SSH_AUTH_SOCK"
)

type SSHSession struct {
	sync.Mutex
	host     string
	proxy    Proxy
	config   *ssh.ClientConfig
	client   *ssh.Client
	session  *ssh.Session
	timeouts Timeouts
}

func NewSSH(ctx Context) (error, Session) {
	err, cfg := ctx2ClientConfig(ctx)
	if err != nil {
		return err, nil
	}

	sshs := &SSHSession{
		host:     net.JoinHostPort(ctx.Host, strconv.Itoa(ctx.Port)),
		config:   cfg,
		proxy:    ctx.Proxy,
		timeouts: ctx.Timeouts,
	}

	if sshs.proxy.Address == "" {
		log.Debug("dialing ssh %s ...", sshs.host)
		if sshs.client, err = ssh.Dial("tcp", sshs.host, sshs.config); err != nil {
			return err, nil
		}
	} else {
		log.Debug("dialing ssh %s via socks5://%s ...", sshs.host, sshs.proxy.String())

		if dialer, err := proxy.SOCKS5("tcp", sshs.proxy.String(), nil, proxy.Direct); err != nil {
			return err, nil
		} else if conn, err := dialer.Dial("tcp", sshs.host); err != nil {
			return err, nil
		} else if c, chans, reqs, err := ssh.NewClientConn(conn, sshs.host, sshs.config); err != nil {
			conn.Close()
			return err, nil
		} else {
			sshs.client = ssh.NewClient(c, chans, reqs)
		}
	}

	if sshs.session, err = sshs.client.NewSession(); err != nil {
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

	// horrible, but there's no other way around
	// with this ssh client library :/
	res := cmdResult{}
	done := make(chan cmdResult)
	timeout := time.After(s.timeouts.Write + s.timeouts.Read)
	go func() {
		out, err := s.session.CombinedOutput(cmd)
		done <- cmdResult{out: out, err: err}
	}()

	select {
	case <-timeout:
		return nil, fmt.Errorf("timeout while sending ssh command to %s", s.host)
	case res = <-done:
		if res.err != nil {
			return res.out, res.err
		}
	}

	return res.out, res.err
}

func (s *SSHSession) Close() {
	s.Lock()
	defer s.Unlock()
	s.session.Close()
	s.client.Close()
}
