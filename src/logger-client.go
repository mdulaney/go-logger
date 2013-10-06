package main

import (
    "bufio"
    "io"
    "flag"
	"fmt"
	"log"
	"net"
	"os"
)

type client struct {
    id int
    ch chan string
    done chan int
}

func FileReader(fileName string, outChan chan string, frDone chan bool) {
    fmt.Printf("Opening file %s\n", fileName)

    f, err := os.Open(fileName)

    if err != nil {
        log.Fatal(err)
    }

    defer f.Close()

    r := bufio.NewReader(f)

    for {
        line, err := r.ReadString('\n')

        if err == io.EOF {
            break
        } else if err != nil {
            log.Fatal(err)
        }

        outChan <-line
    }

    fmt.Printf("Finished reading file\n")
    frDone<-true
}

// TODO: can go routines return values?
func LogReporterClient(c client, addr string, frDone chan bool) {
	fmt.Printf("Connecting to %s\n", addr)

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

    isDone := false

    for ; !isDone; {
        select {
        case data := <-c.ch:
	        bytesWritten, err := conn.Write([]byte(data))

	        if err != nil {
                log.Fatal(err)
	        }
	        fmt.Printf("Sent %d bytes\n", bytesWritten)
        case <- frDone:
            fmt.Printf("Received terminate signal\n")
            isDone = true
        }

    }
    fmt.Printf("Notifying main\n")
    c.done <- 1
}

func main() {
    var inFileName string
    var addr string

    chan1 := make(chan string)
    done1 := make(chan int)
    frDone := make(chan bool)

    cli := client{ 0, chan1 , done1}

	flag.StringVar(&addr, "t", "127.0.0.1:50000", "<ip>:<port>")
    flag.StringVar(&inFileName, "f", "data/infile.txt", "input file")

    flag.Parse()

    go FileReader(inFileName, chan1, frDone)

    go LogReporterClient(cli, addr, frDone)

    <-cli.done
}
