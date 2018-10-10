package plugins

func (p *Plugin) doDefines() {
	p.VM.Set("log", getLOG())
	p.VM.Set("tcp", getTCP())
	p.VM.Set("http", getHTTP())
}
