package main

import (
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/models"
	"github.com/evilsocket/shellz/plugins"
)

func main() {
	log.Format = "{level:color}{level:name}{reset} {message}"

	if err := log.Open(); err != nil {
		panic(err)
	}
	defer log.Close()

	if !noBanner {
		log.Raw(core.Banner)
	}

	if err = models.Init(); err != nil {
		log.Fatal("error while initializing models: %s", err)
	} else if err, Idents, Shells, Groups = models.Load(); err != nil {
		log.Fatal("error while loading data: %s", err)
	} else if len(Shells) == 0 {
		log.Fatal("no shells found on the system, start creating json files inside %s", models.Paths["shells"])
	} else if err = plugins.Load(models.Paths["plugins"]); err != nil {
		log.Fatal("error while loading plugins: %s", err)
	} else {
		log.Debug("loaded %d identities and %d shells", len(Idents), len(Shells))
	}

	if doList {
		showList()
	} else if doEnable != "" {
		runEnable(doEnable, true)
	} else if doDisable != "" {
		runEnable(doDisable, false)
	} else if doTest {
		command = "echo 1" // this should run on every OS ^_^
		runCommand()
	} else if command != "" {
		runCommand()
	} else {
		showHelp()
	}
}
