package main

import (
	"time"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/models"
)

var (
	Idents = models.Identities(nil)
	Shells = models.Shells(nil)
	Groups = models.Groups(nil)

	numWorkers = -1
	command    = ""
	onFilter   = "all"
	onShells   = models.Shells{}
	nShells    = 0
	toOutput   = ""
	doForce    = false
	doList     = false
	doTest     = false
	doTunnel   = false
	doEnable   = ""
	doDisable  = ""
	doStats    = false
	noBanner   = false
	err        = error(nil)

	timeouts = core.Timeouts{
		Connect: 5 * time.Second,
		Read:    500 * time.Millisecond,
		Write:   500 * time.Millisecond,
	}
)
