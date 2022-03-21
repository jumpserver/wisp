package forward

import (
	"io"
	"net"

	"golang.org/x/crypto/ssh"

	"github.com/jumpserver/wisp/pkg/logger"
)

type SSHForward struct {
	Client  *ssh.Client
	DstAddr string

	ln   net.Listener
	addr *net.TCPAddr
}

func (s *SSHForward) Start() error {
	ln, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return err
	}
	s.ln = ln
	go s.run()
	return nil
}

func (s *SSHForward) GetTCPAddr() *net.TCPAddr {
	return s.ln.Addr().(*net.TCPAddr)
}

func (s *SSHForward) Stop() {
	if s.ln != nil {
		if err := s.ln.Close(); err != nil {
			logger.Error(err)
		}
	}
	if err := s.Client.Close(); err != nil {
		logger.Error(err)
	}
}

func (s *SSHForward) String() string {
	return s.DstAddr
}

func (s *SSHForward) run() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if err != net.ErrClosed {
				logger.Errorf("listen %s accept failed: %v", s.ln.Addr(), err)
			}
			return
		}
		go s.forward(conn)
	}
}

func (s *SSHForward) forward(conn net.Conn) {
	defer conn.Close()
	proxyCon, err := s.Client.Dial("tcp", s.DstAddr)
	if err != nil {
		logger.Errorf("ssh.Dial %s failed: %s\n", s.DstAddr, err)
		return
	}
	go func() {
		defer proxyCon.Close()
		if _, err = io.Copy(proxyCon, conn); err != nil && err != io.EOF {
			logger.Errorf("io.Copy local-> proxy err: %s\n", err)
		}
	}()
	if _, err = io.Copy(conn, proxyCon); err != nil && err != io.EOF {
		logger.Errorf("io.Copy proxy -> local err: %s\n", err)
	}
}
