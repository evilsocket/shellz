package session

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/models"
)

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

func authFromFile(sh models.Shell) (err error, auth ssh.AuthMethod) {
	log.Debug("loading ssh key from %s ...", sh.Identity.KeyFile)

	if sh.Identity.KeyFile, err = core.ExpandPath(sh.Identity.KeyFile); err != nil {
		err = fmt.Errorf("error while expanding path '%s': %s", sh.Identity.KeyFile, err)
	} else if key, err := ioutil.ReadFile(sh.Identity.KeyFile); err != nil {
		err = fmt.Errorf("error while reading key file %s: %s", sh.Identity.KeyFile, err)
	} else if signer, err := ssh.ParsePrivateKey(key); err != nil {
		err = fmt.Errorf("error while parsing key file %s: %s", sh.Identity.KeyFile, err)
	} else {
		auth = ssh.PublicKeys(signer)
	}
	return
}

func sh2ClientConfig(sh models.Shell, timeouts core.Timeouts) (error, *ssh.ClientConfig) {
	authMethods := []ssh.AuthMethod{}
	if sh.Identity.Password != "" {
		authMethods = append(authMethods, ssh.Password(sh.Identity.Password))
	}

	if sh.Identity.KeyFile == SSHAgentKey {
		if err, auth := authFromAgent(); err != nil {
			return err, nil
		} else {
			authMethods = append(authMethods, auth)
		}
	} else if sh.Identity.KeyFile != "" {
		if err, auth := authFromFile(sh); err != nil {
			return err, nil
		} else {
			authMethods = append(authMethods, auth)
		}
	}

	return nil, &ssh.ClientConfig{
		Config: ssh.Config{
			Ciphers: sh.Ciphers,
		},
		User:            sh.Identity.Username,
		Auth:            authMethods,
		Timeout:         timeouts.Connect,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}
