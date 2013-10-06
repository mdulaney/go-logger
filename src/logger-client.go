package main

import (
    "flag"
	"fmt"
	"log"
	"net"
	"os"
)

type client struct {
    id int
    c chan string
    done chan int
}

func FileReader(fileName string) {
    fmt.Printf("Opening file %s\n", fileName)

    f, err := os.Open(fileName)

    if err != nil {
        log.Fatal(err)
    }

    defer f.Close()

}

// TODO: can go routines return values?
func LogReporterClient(c client, addr string) {
	fmt.Printf("Connecting to %s\n", addr)

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

	bytesWritten, err := conn.Write([]byte("Here's a message2\n"))

	if err != nil {
		log.Fatal(err)
	}

	bytesWritten, err = conn.Write([]byte("a\n"))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sent %d bytes\n", bytesWritten)

	conn.Close()

    c.done <- 1
}

func main() {
    var inFileName string
    var addr string

    chan1 := make(chan string)
    done1 := make(chan int)

    cli := client{ 0, chan1 , done1}

	flag.StringVar(&addr, "t", "127.0.0.1:50000", "<ip>:<port>")
    flag.StringVar(&inFileName, "f", "data/infile.txt", "input file")

    flag.Parse()

    go FileReader(inFileName)

    go LogReporterClient(cli, addr)

    <-cli.done
}
