package main

import (
    "bufio"
	"encoding/json"
    "flag"
	"fmt"
    "io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

type client struct {
    id int
    ch chan string
    done chan int
}

type config struct {
    numOfClients int `json:"numofclients"`
    server string `json:"server"`
    inFile string `json:"inputfile"`
    delay int
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
func LogReporterClient(c client, addr string, delay time.Duration, frDone chan bool) {
	fmt.Printf("Connecting to %s\n", addr)

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

    isDone := false

    for ; !isDone; {
        time.Sleep(delay * time.Second)

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

func parseConfigFile(fileName string, cfg config) {

    f, err := os.Open(fileName)

    if err != nil {
        log.Fatal(err)
    }

    data, err := ioutil.ReadAll(f)

    if err != nil {
        log.Fatal(err)
    }

    //fmt.Printf("JSON data: %s\n", data)

    err = json.Unmarshal(data, &cfg)

    if err != nil {
        log.Fatal(err)
    }
}

func main() {

    var cfgFileName string

    cfg := config { }
    cfgFile := config { }

	flag.StringVar(&cfg.server, "t", "127.0.0.1:50000", "<ip>:<port>")
    flag.StringVar(&cfg.inFile, "f", "data/infile.txt", "input file")
    flag.IntVar(&cfg.numOfClients, "n", 3, "number of clients")
    flag.IntVar(&cfg.delay, "d", 1, "delay between log events")

    flag.StringVar(&cfgFileName, "c", "data/sample.cfg", "config file")

    flag.Parse()

    if cfgFileName != "" {
        parseConfigFile(cfgFileName, cfgFile)
    }

    clients := make([]client, 0)
    frDone := make(chan bool)

    for i := 0; i < cfg.numOfClients; i++ {
        clients = append(clients, client{i, make(chan string), make(chan int)})
        go LogReporterClient(clients[i], cfg.server, time.Duration(cfg.delay), frDone)
    }

    go FileReader(cfg.inFile, clients, frDone)

    for _, c := range clients {
        <-c.done
    }
}
