package main

import (
	"time"

	"github.com/evilsocket/shellz/models"
	"github.com/evilsocket/shellz/session"
)

var (
	command   = ""
	onFilter  = "*"
	onShells  = models.Shells{}
	nShells   = 0
	toOutput  = ""
	doList    = false
	doTest    = false
	doEnable  = ""
	doDisable = ""
	noBanner  = false
	err       = error(nil)
	idents    = models.Identities(nil)
	shells    = models.Shells(nil)

	timeouts = session.Timeouts{
		Connect: 5 * time.Second,
		Read:    500 * time.Millisecond,
		Write:   500 * time.Millisecond,
	}
)
