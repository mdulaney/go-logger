package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

const (
    PID_FILENAME = "logger-server.pid"
)

type clientInfo struct {
    addr string
}

var gClientMapLocker sync.Mutex
var gClientMap map[string]clientInfo = make(map[string]clientInfo)

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

    gClientMapLocker.Lock()
    delete(gClientMap, c.RemoteAddr().String())
    gClientMapLocker.Unlock()
}

func makePidFile() int {
    var pidFileName = PID_FILENAME

    pid := os.Getpid()

    f, err := os.Create(pidFileName)

    if err != nil {
        log.Fatal(err)
    }

    defer f.Close()

    fmt.Fprintf(f, "%d\n", pid)

    return 0
}

func handleCommandConnection(conn net.Conn) {

    r := bufio.NewReader(conn)

    for {
        cmdStr, err := r.ReadString('\n')

        if err != nil {
            return
        }

        switch {
        case cmdStr == "exit\n":
            return
        case cmdStr == "history\n":
            conn.Write([]byte("Displaying history\n\r"))
        case cmdStr == "clients\n":
            conn.Write([]byte(fmt.Sprintf("%d\n\r", len(gClientMap))))
        }
    }
}

func acceptCommandConnections(addr string) {
    l, err := net.Listen("tcp", addr)

    if err != nil {
        log.Fatal(err)
    }

    defer l.Close()

    for {
        conn, err := l.Accept()

        if err != nil {
            log.Fatal(err)
        }

        go handleCommandConnection(conn)
    }
}

func main() {
	var mainId = 0
	var idx = 1
	var addr = flag.String("l", ":50000", "<ip>:<port>")
    var cmdAddr string

    flag.StringVar(&cmdAddr, "c", ":50001", "<ip>:<port>")

    flag.Parse()

	printLogMsg(mainId, "Listening for connections\n")

    go acceptCommandConnections(cmdAddr)

	l, err := net.Listen("tcp", *addr)

	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

    makePidFile()

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Fatal(err)
		}

		printLogMsg(mainId, "Received a new connection, handling it\n")

        gClientMapLocker.Lock()
        clientAddrString := conn.RemoteAddr().String()
        gClientMap[clientAddrString] = clientInfo { clientAddrString }
        gClientMapLocker.Unlock()

		go handleConnection(conn, idx)
		idx += 1
	}
}


