package main

import (
	"fmt"
	"log"
	"net"

	"chapter4/internal/request"
)

func main() {

	// instead of reading data from file, now we are going to read it from tcp
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", err)
		}

		// tcp instead listen from http
		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Println("error parsing request:", err)
			conn.Close()
			continue
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

		conn.Close()
	}
}
