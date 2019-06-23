package server

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/tidwall/resp"
	"io"
	"log"
	"net"
)

type MockServer struct {
	listener net.Listener
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
		fmt.Print(string(message))
		rd := resp.NewReader(bytes.NewBufferString(string(message)))
		for {
			v, _, err := rd.ReadValue()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Read %s\n", v.Type())
			if v.Type() == resp.Array {
				for i, v := range v.Array() {
					fmt.Printf("  #%d %s, value: '%s'\n", i, v.Type(), v)
				}
			}
		}
		var buf bytes.Buffer
		wr := resp.NewWriter(&buf)
		_ = wr.WriteString("Ok")

		_, err := conn.Write(buf.Bytes())
		if err != nil {
			_ = fmt.Errorf("error in writing to connection %s", err)
			break
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
