package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("message.txt")

	if err != nil {
		log.Fatal("error", "error", err)
	}


	str := ""
	for {
		data := make([]byte, 8)
		// cursor automatically move 8byte 
		n , err := file.Read(data)

		if err != nil {
			break
		}

		i := bytes.IndexByte(data, '\n')
		if i == -1 {
			str += string(data[:n])
		} else {
			// get bytes up to \n, sometime in middle of the data we can have \n
			str += string(data[:i])
			fmt.Printf("read: %s\n", str)
			
			// and store rest of the byte data which had \n middle in order to start the next line from it
			str = string(data[i+1:])
		}
	}
}