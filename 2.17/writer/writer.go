package writer

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func Writer(conn net.Conn, quit chan<- struct{}) {
	defer func() {
		quit <- struct{}{}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("enter the text:")
	for scanner.Scan() {
		msg := scanner.Text()
		conn.Write([]byte(msg + "\n"))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
