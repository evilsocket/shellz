package session

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"sync"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
)

type SSHSession struct {
	sync.Mutex
	host    string
	config  *ssh.ClientConfig
	client  *ssh.Client
	session *ssh.Session
}

func ctx2ClientConfig(ctx Context) (error, *ssh.ClientConfig) {
	var err error

	authMethods := []ssh.AuthMethod{}
	if ctx.Password != "" {
		authMethods = append(authMethods, ssh.Password(ctx.Password))
	}

	if ctx.KeyFile != "" {
		if ctx.KeyFile, err = core.ExpandPath(ctx.KeyFile); err != nil {
			return fmt.Errorf("error while expanding path '%s': %s", ctx.KeyFile, err), nil
		}

		log.Debug("loading ssh key from %s ...", ctx.KeyFile)

		if key, err := ioutil.ReadFile(ctx.KeyFile); err != nil {
			return fmt.Errorf("error while reading key file %s: %s", ctx.KeyFile, err), nil
		} else if signer, err := ssh.ParsePrivateKey(key); err != nil {
			return fmt.Errorf("error while parsing key file %s: %s", ctx.KeyFile, err), nil
		} else {
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		}
	}

	return nil, &ssh.ClientConfig{
		User:            ctx.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func NewSSH(ctx Context) (error, Session) {
	err, cfg := ctx2ClientConfig(ctx)
	if err != nil {
		return err, nil
	}

	sshs := &SSHSession{
		host:   fmt.Sprintf("%s:%d", ctx.Address.String(), ctx.Port),
		config: cfg,
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

func (s *SSHSession) Exec(cmd string) ([]byte, error) {
	s.Lock()
	defer s.Unlock()
	return s.session.CombinedOutput(cmd)
}

func (s *SSHSession) Close() {
	s.Lock()
	defer s.Unlock()
	s.session.Close()
	s.client.Close()
}
