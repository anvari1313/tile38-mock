package server

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/tidwall/resp"
	"io"
	"log"
	"net"
)

type MockServer struct {
	listener net.Listener
	Mock     map[string]func(resp.Value) (resp.Value, error)
}

func (s *MockServer) Init(address string) error {
	var err error = nil
	s.listener, err = net.Listen("tcp4", address)
	if err != nil {
		return err
	}
	go func() {
		defer s.listener.Close()

		for {
			conn, err := s.listener.Accept()
			fmt.Println("New connection accepted")
			if err != nil {
				log.Fatalf("Ffailed to accpect connection, %s", err)
			}
			go s.handle(conn)
		}
	}()

	return err
}

func (s *MockServer) handle(conn net.Conn) {
	for {
		message := readCommand(conn)

		rd := resp.NewReader(bytes.NewBufferString(string(message)))
		for {
			var buf bytes.Buffer
			wr := resp.NewWriter(&buf)

			v, _, err := rd.ReadValue()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			if v.Type() == resp.Array {
				if handler, ok := s.Mock[v.Array()[0].String()]; ok {
					value, err := handler(v)
					if err != nil {
						_ = wr.WriteError(err)
					} else {
						_ = wr.WriteValue(value)
					}
				} else {
					_ = wr.WriteError(errors.New("command not found"))
				}

			} else {
				_ = wr.WriteError(errors.New("not implemented"))
				_, _ = conn.Write(buf.Bytes())
			}
			_, _ = conn.Write(buf.Bytes())
		}
	}
}

func readCommand(conn net.Conn) []byte {
	reader := bufio.NewReader(conn)

	var message []byte
	for {
		buf := make([]byte, 10)
		n, err := reader.Read(buf)
		if err != nil {
			break
		}
		message = append(message, buf[:n]...)
		if n < len(buf) {
			break
		}
	}

	return message
}
