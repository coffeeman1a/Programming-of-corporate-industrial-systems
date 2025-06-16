package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("Connection error: %v", err)
	}
	defer conn.Close()

	serverReader := bufio.NewReader(conn)
	serverWriter := bufio.NewWriter(conn)
	console := bufio.NewReader(os.Stdin)

	// async server reading
	go func() {
		for {
			line, err := serverReader.ReadString('\n')
			if err != nil {
				os.Exit(0)
			}
			fmt.Print(line)
		}
	}()

	for {
		text, _ := console.ReadString('\n')
		text = strings.TrimSpace(text)
		serverWriter.WriteString(text + "\n")
		serverWriter.Flush()
	}
}
