package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/session"
)

const (
	defaultHost     = "localhost"
	defaultPort     = 22
	defaultIdentity = "default"
	defaultType     = "ssh"
)

type Shell struct {
	Name         string    `json:"name"`
	Host         string    `json:"host"`
	Address      net.IP    `json:"address"`
	Port         int       `json:"port"`
	IdentityName string    `json:"identity"`
	Type         string    `json:"type"`
	Identity     *Identity `json:"-"`
	Path         string    `json:"-"`
}

func LoadShell(path string, idents Identities) (err error, shell Shell) {
	shell = Shell{
		Path:         path,
		Host:         defaultHost,
		Port:         defaultPort,
		Type:         defaultType,
		IdentityName: defaultIdentity,
		Address:      net.IP{0},
	}

	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	raw, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	if err = json.Unmarshal(raw, &shell); err != nil {
		return fmt.Errorf("error decoding '%s': %s", path, err), shell
	} else if ident, found := idents[shell.IdentityName]; !found {
		return fmt.Errorf("shell '%s' referenced an unknown identity '%s'", path, shell.IdentityName), shell
	} else {
		shell.Identity = &ident
	}

	return
}

func (sh Shell) NewSession() (error, session.Session) {
	if sh.Address[0] == 0 {
		if addrs, err := net.LookupIP(sh.Host); err != nil {
			return fmt.Errorf("could not resolve host '%s' for shell '%s'", sh.Host, sh.Name), nil
		} else {
			sh.Address = addrs[0]
			log.Debug("host %s resolved to %s", sh.Host, sh.Address)
		}
	}

	if mgr, found := session.Manager[sh.Type]; found {
		return mgr(sh.Address, sh.Port, sh.Identity.Username, sh.Identity.Password, sh.Identity.KeyFile)
	}
	return fmt.Errorf("session type %s for shell %s is not supported", sh.Type, sh.Name), nil
}
