# Tile38 Mock

This project will make up a mock server for [Tile38](https://tile38.com/) that you can define your desired responses with respect of passed command. It can run on a port and you can communicate with it by redis protocol. Best for testing your application with low overhead.

## Using

To start using Tile38 Mock install go and get it with:
```
go get -u github.com/anvari1313/tile38-mock
```

## How to use it

Each instance of ```tile38_mock.MockServer``` is responsible for an instance of tile38 server.

1. First of all import the package
2. Create a new instance  with ```CreateMockServer```.
3. Set commands and their responses with ```SetStringResponse```.
4. Init the server with address and port.

Sample Code:

```go
import "github.com/anvari1313/tile38-mock"
```
...
```go
s := tile38_mock.CreateMockServer() // Create a new mock server instance

// Set command and responses
s.SetStringResponse([]string{"set", "key", "value"}, []string{"OK"})        // SET KEY VALUE -> OK
s.SetStringResponse([]string{"get", "key"}, []string{"VALUE"})              // GET KEY -> VALUE 
s.SetStringResponse([]string{"got", "key"}, []string{"VALUE1", "VALUE1"})   // GOT KEY -> 1)VALUE1 2)VALUE2

_ = s.Init(":1024") // Run the server on port 1024
```

**NOTE:** Be aware that tcp server will run in a separate go routine and the init function is non blocking.