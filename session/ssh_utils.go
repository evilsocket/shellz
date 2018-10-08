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
