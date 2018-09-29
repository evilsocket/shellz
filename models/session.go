package models

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"sync"
)

type Session struct {
	sync.Mutex

	client  *ssh.Client
	session *ssh.Session
}

func NewSessionFor(sh Shell) (err error, s *Session) {
	s = &Session{}
	if s.client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", sh.Address.String(), sh.Port), &sh.Identity.Config); err != nil {
		return
	} else if s.session, err = s.client.NewSession(); err != nil {
		s.client.Close()
		return
	}
	return
}

func (s *Session) Close() {
	s.Lock()
	defer s.Unlock()
	s.session.Close()
	s.client.Close()
}

func (s *Session) Exec(cmd string) ([]byte, error) {
	s.Lock()
	defer s.Unlock()
	return s.session.CombinedOutput(cmd)
}
