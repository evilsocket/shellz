package models

type Proxy struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p Proxy) Empty() bool {
	return p.Address == ""
}
