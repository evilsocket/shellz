package sessions

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"sync"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
)

type SSHSession struct {
	sync.Mutex
	config  *ssh.ClientConfig
	client  *ssh.Client
	session *ssh.Session
}

func NewSSH(address net.IP, port int, user string, pass string, keyFile string) (error, Session) {
	var err error

	authMethods := []ssh.AuthMethod{}
	if pass != "" {
		authMethods = append(authMethods, ssh.Password(pass))
	}

	if keyFile != "" {
		if keyFile, err = core.ExpandPath(keyFile); err != nil {
			return fmt.Errorf("error while expanding path '%s': %s", keyFile, err), nil
		}

		log.Debug("loading ssh key from %s ...", keyFile)
		key, err := ioutil.ReadFile(keyFile)
		if err != nil {
			return fmt.Errorf("error while reading key file %s: %s", keyFile, err), nil
		} else if signer, err := ssh.ParsePrivateKey(key); err != nil {
			return fmt.Errorf("error while parsing key file %s: %s", keyFile, err), nil
		} else {
			authMethods = append(authMethods, ssh.PublicKeys(signer))
		}
	}

	host := fmt.Sprintf("%s:%d", address.String(), port)
	sshs := &SSHSession{
		config: &ssh.ClientConfig{
			User:            user,
			Auth:            authMethods,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}

	if sshs.client, err = ssh.Dial("tcp", host, sshs.config); err != nil {
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
