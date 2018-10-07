package main

import (
	"time"

	"github.com/evilsocket/shellz/models"
	"github.com/evilsocket/shellz/session"
)

var (
	Idents = models.Identities(nil)
	Shells = models.Shells(nil)
	Groups = models.Groups(nil)

	command   = ""
	onFilter  = "all"
	onShells  = models.Shells{}
	nShells   = 0
	toOutput  = ""
	doForce   = false
	doList    = false
	doTest    = false
	doEnable  = ""
	doDisable = ""
	doStats   = false
	noBanner  = false
	err       = error(nil)

	timeouts = session.Timeouts{
		Connect: 5 * time.Second,
		Read:    500 * time.Millisecond,
		Write:   500 * time.Millisecond,
	}
)
