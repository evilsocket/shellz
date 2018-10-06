package plugins

import (
	"bufio"
	"net"
)

type tcpPackage struct{}

var tcp = tcpPackage{}

func getTCP() tcpPackage {
	return tcp
}

type tcpClient struct {
	con    net.Conn
	reader *bufio.Reader
}

type tcpResponse struct {
	Error error
	Size  int
	Raw   []byte
}

func (m tcpPackage) Connect(host string) *tcpClient {
	con, err := net.Dial("tcp", host)
	if err != nil {
		return nil
	}
	return &tcpClient{
		con:    con,
		reader: bufio.NewReader(con),
	}
}

func (c *tcpClient) Read(size int) tcpResponse {
	buf := make([]byte, size)
	read, err := c.con.Read(buf)
	return tcpResponse{
		Error: err,
		Size:  read,
		Raw:   buf[0:read],
	}
}

func (c *tcpClient) ReadUntil(delim rune) tcpResponse {
	message, err := c.reader.ReadString(byte(delim))
	return tcpResponse{
		Error: err,
		Size:  len(message),
		Raw:   []byte(message),
	}
}

func (c *tcpClient) Write(buf []byte) tcpResponse {
	wrote, err := c.con.Write(buf)
	return tcpResponse{
		Error: err,
		Size:  wrote,
	}
}

func (c *tcpClient) Close() error {
	return c.con.Close()
}
