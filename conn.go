package stdio

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"time"
)

// DialStdio connects the current process stdio.
func DialStdio() *Conn {
	return &Conn{os.Stdin, os.Stdout, os.Getpid(), os.Getppid()}
}

// DialCommand wraps the stdio of a subprocess.
func DialCommand(cmd *exec.Cmd) (*Conn, error) {
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	writer, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return &Conn{reader, writer, os.Getpid(), cmd.Process.Pid}, nil
}

// Addr implements net.Addr for stdin/stdout.
type Addr struct {
	pid int
}

// Network returns the network type as a string.
func (a Addr) Network() string {
	return "stdio"
}

func (a Addr) String() string {
	return fmt.Sprintf("pid:%d", a.pid)
}

// Conn implements a net.Conn via stdin/stdout.
type Conn struct {
	reader io.ReadCloser
	writer io.WriteCloser

	local  int
	remote int
}

func (s *Conn) Read(p []byte) (int, error) {
	return s.reader.Read(p)
}

func (s *Conn) Write(p []byte) (int, error) {
	return s.writer.Write(p)
}

// Close closes both streams.
func (s *Conn) Close() error {
	err1 := s.reader.Close()
	err2 := s.writer.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

// LocalAddr returns stdio addr.
func (s *Conn) LocalAddr() net.Addr {
	return Addr{s.local}
}

// RemoteAddr returns stdio addr.
func (s *Conn) RemoteAddr() net.Addr {
	return Addr{s.remote}
}

// SetDeadline sets the read/write deadline.
func (s *Conn) SetDeadline(t time.Time) error {
	err1 := s.SetReadDeadline(t)
	err2 := s.SetWriteDeadline(t)
	if err1 != nil {
		return err1
	}
	return err2
}

// SetReadDeadline sets the read/write deadline.
func (s *Conn) SetReadDeadline(t time.Time) error {
	d, ok := s.reader.(interface {
		SetReadDeadline(time.Time) error
	})
	if !ok {
		return &net.OpError{Op: "set", Net: "stdio", Addr: s.LocalAddr(),
			Err: errors.New("read deadline not supported")}
	}
	return d.SetReadDeadline(t)
}

// SetWriteDeadline sets the read/write deadline.
func (s *Conn) SetWriteDeadline(t time.Time) error {
	d, ok := s.writer.(interface {
		SetWriteDeadline(time.Time) error
	})
	if !ok {
		return &net.OpError{Op: "set", Net: "stdio", Addr: s.LocalAddr(),
			Err: errors.New("write deadline not supported")}
	}
	return d.SetWriteDeadline(t)
}
