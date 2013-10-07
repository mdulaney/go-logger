package main

import (
	"bufio"
	"container/list"
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

// Global declarations
var gClientMapLocker sync.Mutex
var gClientMap map[string]clientInfo = make(map[string]clientInfo)

var gHistoryLock sync.Mutex
var gLogHistory = list.New()
var gMaxHistoryLen = 10

func printLogMsg(id int, s string) {
	fmt.Printf("[%d] %s", id, s)
}

func updateHistory(log string) {

    gHistoryLock.Lock()
    gLogHistory.PushFront(log)

    if gLogHistory.Len() > gMaxHistoryLen {
        gLogHistory.Remove(gLogHistory.Back())
    }
    gHistoryLock.Unlock()
}

func getHistoryString() string {

    history := ""

    gHistoryLock.Lock()
    for l := gLogHistory.Front(); l != nil; l = l.Next() {
        history = history + l.Value.(string)
    }

    gHistoryLock.Unlock()

    return history
}

func logAggregator(logChan chan string) {
    var logStr string

    fmt.Printf("Started logAggregator()\n")
    for {
        select {
        case logStr = <-logChan:
            printLogMsg(0, logStr)
            updateHistory(logStr)
        }
    }
}

func handleConnection(c net.Conn, id int, logChan chan string) {

	r := bufio.NewReader(c)
	for {
		logStr, err := r.ReadString('\n')

		if err != nil {
			break
		}

        logChan <- logStr
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

func writeResponse(c net.Conn, response string) {
    c.Write([]byte(response + "\r"))
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
            writeResponse(conn, getHistoryString())
        case cmdStr == "clients\n":
            writeResponse(conn, fmt.Sprintf("%d\n", len(gClientMap)))
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

    logChan := make(chan string)

    go logAggregator(logChan)
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

		go handleConnection(conn, idx, logChan)
		idx += 1
	}
}


