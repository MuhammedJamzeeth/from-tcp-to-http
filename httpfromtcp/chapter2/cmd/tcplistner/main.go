package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(file io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer file.Close()
		defer close(out)

		str := ""
		for {
			data := make([]byte, 8)
			// cursor automatically move 8byte
			n, err := file.Read(data)

			if err != nil {
				// the last line has no \n, so flush what is left over
				if str != "" {
					out <- str
				}
				break
			}

			// one chunk can hold more than one \n, so keep splitting until none left
			chunk := data[:n]
			for {
				i := bytes.IndexByte(chunk, '\n')
				if i == -1 {
					break
				}
				// get bytes up to \n, sometime in middle of the data we can have \n
				str += string(chunk[:i])
				out <- str
				str = ""
				// and continue from the byte right after the \n
				chunk = chunk[i+1:]
			}
			str += string(chunk)

		}

	}()

	return out
}

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

		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}

	}
}
