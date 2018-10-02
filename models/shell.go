package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"

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
	Name         string `json:"name"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	IdentityName string `json:"identity"`
	Type         string `json:"type"`
	Enabled      bool   `json:"enabled"`

	Address  net.IP    `json:"-"`
	Identity *Identity `json:"-"`
	Path     string    `json:"-"`
}

func LoadShell(path string, idents Identities) (err error, shell Shell) {
	shell = Shell{
		Enabled:      true,
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

func (sh Shell) Save() error {
	if data, err := json.Marshal(sh); err != nil {
		return err
	} else if err = ioutil.WriteFile(sh.Path, data, 0644); err != nil {
		return err
	}
	return nil
}

func (sh Shell) NewSession(timeouts session.Timeouts) (error, session.Session) {
	if !strings.Contains(sh.Host, "://") && sh.Address[0] == 0 {
		if addrs, err := net.LookupIP(sh.Host); err != nil {
			return fmt.Errorf("could not resolve host '%s' for shell '%s'", sh.Host, sh.Name), nil
		} else {
			sh.Address = addrs[0]
			log.Debug("host %s resolved to %s", sh.Host, sh.Address)
		}
	}

	ctx := session.Context{
		Host:     sh.Host,
		Address:  sh.Address,
		Port:     sh.Port,
		Username: sh.Identity.Username,
		Password: sh.Identity.Password,
		KeyFile:  sh.Identity.KeyFile,
		Timeouts: timeouts,
	}

	// check the built in session managers
	if mgr := session.GetManager(sh.Type); mgr != nil {
		return mgr(ctx)
	}
	// check for user provided plugins
	if plugin := session.GetPlugin(sh.Type); plugin != nil {
		return session.NewPluginSession(plugin.Clone(), ctx)
	}

	return fmt.Errorf("session type %s for shell %s is not supported", sh.Type, sh.Name), nil
}
