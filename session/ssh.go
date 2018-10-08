package session

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
)

const (
	SSHAgentKey = "@agent"
	SSHAuthSock = "SSH_AUTH_SOCK"
)

type SSHSession struct {
	sync.Mutex
	host     string
	config   *ssh.ClientConfig
	client   *ssh.Client
	session  *ssh.Session
	timeouts Timeouts
}

func authFromAgent() (err error, auth ssh.AuthMethod) {
	log.Debug("asking for ssh key to %s", SSHAuthSock)

	if socket := os.Getenv(SSHAuthSock); socket == "" {
		err = fmt.Errorf("error while connecting to ssh-agent (cant find %s variable)", SSHAuthSock)
	} else if conn, err := net.Dial("unix", socket); err != nil {
		err = fmt.Errorf("error while connecting to ssh-agent '%s': %s", socket, err)
	} else {
		auth = ssh.PublicKeysCallback(agent.NewClient(conn).Signers)
	}
	return
}

func authFromFile(ctx Context) (err error, auth ssh.AuthMethod) {
	log.Debug("loading ssh key from %s ...", ctx.KeyFile)

	if ctx.KeyFile, err = core.ExpandPath(ctx.KeyFile); err != nil {
		err = fmt.Errorf("error while expanding path '%s': %s", ctx.KeyFile, err)
	} else if key, err := ioutil.ReadFile(ctx.KeyFile); err != nil {
		err = fmt.Errorf("error while reading key file %s: %s", ctx.KeyFile, err)
	} else if signer, err := ssh.ParsePrivateKey(key); err != nil {
		err = fmt.Errorf("error while parsing key file %s: %s", ctx.KeyFile, err)
	} else {
		auth = ssh.PublicKeys(signer)
	}
	return
}

func ctx2ClientConfig(ctx Context) (error, *ssh.ClientConfig) {
	authMethods := []ssh.AuthMethod{}
	if ctx.Password != "" {
		authMethods = append(authMethods, ssh.Password(ctx.Password))
	}

	if ctx.KeyFile == SSHAgentKey {
		if err, auth := authFromAgent(); err != nil {
			return err, nil
		} else {
			authMethods = append(authMethods, auth)
		}
	} else if ctx.KeyFile != "" {
		if err, auth := authFromFile(ctx); err != nil {
			return err, nil
		} else {
			authMethods = append(authMethods, auth)
		}
	}

	return nil, &ssh.ClientConfig{
		Config: ssh.Config{
			Ciphers: ctx.Ciphers,
		},
		User:            ctx.Username,
		Auth:            authMethods,
		Timeout:         ctx.Timeouts.Connect,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func NewSSH(ctx Context) (error, Session) {
	err, cfg := ctx2ClientConfig(ctx)
	if err != nil {
		return err, nil
	}

	sshs := &SSHSession{
		host:     fmt.Sprintf("%s:%d", ctx.Address.String(), ctx.Port),
		config:   cfg,
		timeouts: ctx.Timeouts,
	}

	if sshs.client, err = ssh.Dial("tcp", sshs.host, sshs.config); err != nil {
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
