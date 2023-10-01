package main

import (
	"fmt"
	"strings"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	 l, err := net.Listen("tcp", "0.0.0.0:4221")
	 if err != nil {
	 	fmt.Println("Failed to bind to port 4221")
	 	os.Exit(1)
	 }
	conn, err := l.Accept()
	if err != nil {
	 	fmt.Println("Error accepting connection: ", err.Error())
	 	os.Exit(1)
	 }
	b := make([]byte, 1024)
	_, err = conn.Read(b)
	if err != nil {
	 	fmt.Println("Error reading data: ", err.Error())
	 	os.Exit(1)
	 }
	if len(b) > 0 {
		data := string(b)
		path := strings.Split(data, " ")[1]
		if path != "/" {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		} else {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		}
	}
	conn.Close()
	
}
