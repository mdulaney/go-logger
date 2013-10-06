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

func FileReader(fileName string, clients []client, frDone chan bool) {
    fmt.Printf("Opening file %s\n", fileName)

    f, err := os.Open(fileName)

    if err != nil {
        log.Fatal(err)
    }

    defer f.Close()

    r := bufio.NewReader(f)

    totalClients := len(clients)
    curCli := 0

    for {
        line, err := r.ReadString('\n')

        if err == io.EOF {
            break
        } else if err != nil {
            log.Fatal(err)
        }

        clients[curCli].ch <-line
        curCli = (curCli + 1) % totalClients
    }

    fmt.Printf("Finished reading file\n")

    for i := 0; i < len(clients); i++ {
        frDone<-true
    }
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
    fmt.Printf("[%d] Notifying main\n", c.id)
    c.done <- 1
}

func main() {
    var inFileName string
    var addr string

	flag.StringVar(&addr, "t", "127.0.0.1:50000", "<ip>:<port>")
    flag.StringVar(&inFileName, "f", "data/infile.txt", "input file")

    flag.Parse()

    clients := make([]client, 0)
    frDone := make(chan bool)

    chan1 := make(chan string)
    done1 := make(chan int)
    clients = append(clients, client{ 0, chan1 , done1})

    chan2 := make(chan string)
    done2 := make(chan int)
    clients = append(clients, client{ 0, chan2 , done2})

    go LogReporterClient(clients[0], addr, frDone)
    go LogReporterClient(clients[1], addr, frDone)

    go FileReader(inFileName, clients, frDone)

    for _, c := range clients {
        <-c.done
    }
}
