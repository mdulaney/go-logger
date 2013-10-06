package main

import (
    "bufio"
    "flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func CommanderClient(addr string) {
	fmt.Printf("Connecting to %s\n", addr)

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

    in := bufio.NewReader(os.Stdin)
    networkIn := bufio.NewReader(conn)

	defer conn.Close()

    for {
        fmt.Printf("> ")
        input, err := in.ReadString('\n')

        if err != nil {
            log.Fatal(err)
        }

        conn.Write([]byte(input))

        if input == "exit\n" {
            return
        }

        // TODO: enhance to deal with timeouts
        result, err := networkIn.ReadString('\r')

        if err != nil {
            log.Fatal(err)
        }

        result = strings.TrimSuffix(result, "\r")

        fmt.Printf(result)
    }
}

func main() {

    var addr string

	flag.StringVar(&addr, "t", "127.0.0.1:50001", "<ip>:<port>")

    CommanderClient(addr)
}
