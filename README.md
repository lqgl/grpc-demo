# gRPC Demo
This repository contains a small demo application showing how to use gRPC in Go. 

## Installation and setup
To run the demo, you need to have Go installed on your machine. You can then clone the repository and run the following commands:
```
$ cd grpc-server
$ go mod tidy
$ ./grpc-server.exe
```

This should start the gRPC server, which listens on port 5001 by default.

## Using the demo
Once the server is running, you can use a gRPC client to interact with it. The demo includes a simple command-line client that you can use to test the service. To use it, run the following commands:
```
$ cd grpc-client
$ go mod tidy
$ go run main.go
```

## License
This code is released under the [MIT License](https://opensource.org/license/mit/).

## Contributing
If you find a bug or have a suggestion for improvement, feel free to open an issue or submit a pull request. We welcome contributions from anyone!
