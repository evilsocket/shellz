package models

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
)

type Identity struct {
	Name     string           `json:"name"`
	Username string           `json:"username"`
	KeyFile  string           `json:"key"`
	Signer   ssh.Signer       `json:"-"`
	Config   ssh.ClientConfig `json:"-"`
	Password string           `json:"password"`
	Path     string           `json:"-"`
}

func LoadIdent(path string) (err error, ident Identity) {
	ident = Identity{Path: path}
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	raw, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	if err = json.Unmarshal(raw, &ident); err != nil {
		return
	}

	authMethods := []ssh.AuthMethod{}
	if ident.Password != "" {
		authMethods = append(authMethods, ssh.Password(ident.Password))
	}

	if ident.KeyFile != "" {
		if ident.KeyFile, err = core.ExpandPath(ident.KeyFile); err != nil {
			log.Fatal("error while expanding path '%s': %s", ident.KeyFile, err)
		}

		log.Debug("loading ssh key from %s ...", ident.KeyFile)
		key, err := ioutil.ReadFile(ident.KeyFile)
		if err != nil {
			return fmt.Errorf("error while reading key file %s: %s", ident.KeyFile, err), ident
		} else if ident.Signer, err = ssh.ParsePrivateKey(key); err != nil {
			return fmt.Errorf("error while parsing key file %s: %s", ident.KeyFile, err), ident
		} else {
			authMethods = append(authMethods, ssh.PublicKeys(ident.Signer))
		}
	}

	// TODO: support host specific verification keys
	ident.Config = ssh.ClientConfig{
		User:            ident.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return
}
