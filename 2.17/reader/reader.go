package reader

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func Reader(conn net.Conn, quit <-chan struct{}) {
	for {
		select {
		case <-quit:
			return
		default:
			response, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(response)
		}
	}
}
