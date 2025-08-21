package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"tcp-client/reader"
	"tcp-client/writer"
	"time"
)

func main() {

	host := flag.String("host", "tcpbin.com", "Configure the host of TCP server")
	port := flag.String("port", "4242", "Configure the port of TCP server")
	timeout := flag.Duration("timeout", 10*time.Second, "Configure the timeout of connection")

	flag.Parse()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", *host, *port), *timeout)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connection has established")
	defer conn.Close()

	quit := make(chan struct{}, 1)

	go writer.Writer(conn, quit)
	go reader.Reader(conn, quit)

	<-quit
	fmt.Println("program has finished with EOF")
}
