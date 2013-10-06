package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

func printLogMsg(id int, s string) {
	fmt.Printf("[%d] %s", id, s)
}

func handleConnection(c net.Conn, id int) {
	
	r := bufio.NewReader(c)
	for {
		logStr, err := r.ReadString('\n')

		if err != nil {
			break
		}

		printLogMsg(id, logStr)

	}
}

func main() {
	var mainId = 0
	var idx = 1
	var addr = flag.String("-l", ":50000", "<ip>:<port>")

    flag.Parse()

	printLogMsg(mainId, "Listening for connections\n")	

	l, err := net.Listen("tcp", *addr)

	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Fatal(err)
		}

		printLogMsg(mainId, "Received a new connection, handling it\n")	
		go handleConnection(conn, idx)
		idx += 1
	}
}


