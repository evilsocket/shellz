package session

var (
	managers = map[string]Handler{
		"ssh":    NewSSH,
		"telnet": NewTelnet,
	}
)

func Get(name string) Handler {
	return managers[name]
}
