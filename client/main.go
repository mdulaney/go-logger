package main

import (
    "flag"
	"fmt"
	"log"
	"net"
)

func main() {
	var addr = flag.String("t", "127.0.0.1:50000", "<ip>:<port>")

    flag.Parse()

	fmt.Printf("Connecting to %s\n", *addr)

	conn, err := net.Dial("tcp", *addr)

	if err != nil {
		log.Fatal(err)
	}

	bytesWritten, err := conn.Write([]byte("Here's a message\n"))

	if err != nil {
		log.Fatal(err)
	}

	bytesWritten, err = conn.Write([]byte("a\n"))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sent %d bytes\n", bytesWritten)


	conn.Close()

}
