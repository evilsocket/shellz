package session

func (p *Plugin) doDefines() error {
	p.vm.Set("log", newLogManager())
	p.vm.Set("tcp", newTcpManager())
	p.vm.Set("http", newHttpClient())
	return nil
}
