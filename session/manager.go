package session

var (
	Manager = map[string]Handler{
		"ssh":    NewSSH,
		"telnet": NewTelnet,
	}
)
