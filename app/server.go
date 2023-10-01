package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Request struct {
	Method, Path, HTTPVersion string
	Headers                   map[string]string
}

func NewRequest(buf []byte) Request {
	s := strings.SplitN(string(buf), " ", 3)
	sH := strings.Split(s[2], "\r\n")
	res := Request{
		Method:      s[0],
		Path:        s[1],
		HTTPVersion: sH[0],
		Headers:     map[string]string{},
	}
	for _, v := range sH {
		kargs := strings.Split(v, ": ")
		if len(kargs) > 1 {
			res.Headers[kargs[0]] = kargs[1]
		}

	}
	return res
}

func WriteTextRes(conn net.Conn, body string, status int) {
	res := fmt.Sprintf("HTTP/1.1 %d OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", status, len(body), body)
	conn.Write([]byte(res))
}

func WriteFileRes(conn net.Conn, body []byte, status int) {
	res := fmt.Sprintf("HTTP/1.1 %d OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", status, len(body), body)
	conn.Write([]byte(res))
}

func HandleRequest(conn net.Conn) {
	b := make([]byte, 1024)
	bytesRead, err := conn.Read(b)
	if err != nil {
		fmt.Println("Error reading data: ", err.Error())
		os.Exit(1)
	}
	if len(b) > 0 {
		res := NewRequest(b[:bytesRead])
		path := res.Path
		paths := strings.Split(path, "/")
		if paths[1] == "echo" {
			body := path[6:]
			WriteTextRes(conn, body, 200)
		} else if paths[1] == "user-agent" {
			agent := res.Headers["User-Agent"]
			WriteTextRes(conn, agent, 200)
		} else if paths[1] == "files" {
			filename := path[6:]
			dir := os.Args[2]
			filePath := fmt.Sprintf("%s/%s", dir, filename)
			file, err := os.ReadFile(filePath)
			if err != nil {
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			} else {
				WriteFileRes(conn, file, 200)
			}
		} else if path == "/" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
	}
	conn.Close()
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go HandleRequest(conn)
	}
}
