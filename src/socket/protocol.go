package socket

import (
	"bytes"
	"fmt"
	"github.com/nlpodyssey/gopickle/pickle"
)

const (
	commandTerminator    = "<F2B_END_COMMAND>"
	pingCommand          = "ping"
	socketReadBufferSize = 10000
)

func (s *Fail2BanSocket) sendCommand(command []string) (interface{}, error) {
	err := s.write(command)
	if err != nil {
		return nil, err
	}
	return s.read()
}

func (s *Fail2BanSocket) write(command []string) error {
	err := s.encoder.Encode(command)
	if err != nil {
		return err
	}
	_, err = s.socket.Write([]byte(commandTerminator))
	if err != nil {
		return err
	}
	return nil
}

func (s *Fail2BanSocket) read() (interface{}, error) {
	buf := make([]byte, socketReadBufferSize)
	_, err := s.socket.Read(buf)
	if err != nil {
		return nil, err
	}

	bufReader := bytes.NewReader(buf)
	unpickler := pickle.NewUnpickler(bufReader)

	unpickler.FindClass = func(module, name string) (interface{}, error) {
		if module == "builtins" && name == "str" {
			return &Py_builtins_str{}, nil
		}
		return nil, fmt.Errorf("class not found: " + module + " : " + name)
	}

	return unpickler.Load()
}