package tile38_mock

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/tidwall/resp"
	"io"
	"log"
	"net"
	"strings"
)

type MockServer struct {
	listener    net.Listener
	cmdResponse map[string]resp.Value
}

func CreateMockServer() *MockServer {
	s := MockServer{}
	s.cmdResponse = make(map[string]resp.Value)
	return &s
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
			if err != nil {
				log.Fatalf("failed to accpect connection, %s", err)
			}
			go s.handle(conn)
		}
	}()

	return err
}

func (s *MockServer) SetStringResponse(cmd []string, res []string) {
	var builder strings.Builder
	for _, c := range cmd {
		builder.WriteString(strings.ToUpper(c))
		builder.WriteString(" ")
	}

	var r resp.Value
	if len(res) == 1 {
		r = resp.StringValue(res[0])
	} else {
		arr := make([]resp.Value, len(res))
		for i, r := range res {
			arr[i] = resp.StringValue(r)
		}
		r = resp.ArrayValue(arr)
	}
	s.cmdResponse[builder.String()] = r
}

func (s *MockServer) handle(conn net.Conn) {
	for {
		message := read(conn)

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
				var builder strings.Builder
				for _, v := range v.Array() {
					builder.WriteString(v.String())
					builder.WriteString(" ")
				}

				if r, ok := s.cmdResponse[strings.ToUpper(builder.String())]; ok {
					_ = wr.WriteValue(r)
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

func read(conn net.Conn) []byte {
	reader := bufio.NewReader(conn)

	var message []byte
	for {
		buf := make([]byte, 1024)
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
