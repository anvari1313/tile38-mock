package tile38_mock

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/go-redis/redis"
)

func TestCreateMockServer(t *testing.T) {
	s := CreateMockServer()
	if s == nil {
		t.Error("Create mock server not return the mock server object")
	}
}

func TestMockServer_Init1(t *testing.T) {
	// Find an open port that can the server listen on it up and
	// Simply check that the server is up and can
	// accept connections.

	// Choose a random free port
	l, err := net.Listen("tcp4", ":0")
	if err != nil {
		t.Fatal(err)
	}
	address := l.Addr().String()
	l.Close() // Close the listener to free up the selected port

	s := CreateMockServer()
	err = s.Init(address)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := net.Dial("tcp4", address)
	defer conn.Close()
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
}

func TestMockServer_SetStringResponse(t *testing.T) {
	// Test scenario is:
	// 1. Running an instance of tile38 mock server
	// 2. Set commands and responses
	// 3. Run commands with go-redis client

	// Step 1
	s := CreateMockServer()
	// Choose a random free port
	l, err := net.Listen("tcp4", ":0")
	if err != nil {
		t.Fatal(err)
	}
	address := l.Addr().String()
	l.Close() // Close the listener to free up the selected port

	err = s.Init(address)

	if err != nil {
		t.Errorf("Error in initializing server, %s", err)
	}

	// Step 2
	s.SetStringResponse([]string{"set", "key", "value"}, []string{"OK"})
	s.SetStringResponse([]string{"get", "key"}, []string{"VALUE"})
	s.SetStringResponse([]string{"got", "key"}, []string{"VALUE1", "VALUE1"})

	// Step 3
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	cmd := redis.NewStringCmd("GET", "key")
	err = client.Process(cmd)
	if err != nil {
		fmt.Printf("error is %s\n", err)
	}
	result, err := cmd.Result()
	if err != nil {
		fmt.Printf("error is %s\n", err)
	}

	if !strings.EqualFold("value", result) {
		t.Error("gotten value is not as the same with expected value")
	}
}
