package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWSOrderbook(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client to orderbook feed:", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("orderdata -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second * 2)
	}
}
func (s *Server) handleWS(ws *websocket.Conn) {
	//Recieve the connection
	fmt.Printf("new incoming connection from client:", ws.RemoteAddr())

	//Maintain the connection
	s.conns[ws] = true

	s.readLoop(ws)

}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)

	for {
		n, err := ws.Read(buf)

		if err != nil {
			if err == io.EOF {
				break //break when the error is ENd of File
			}
			fmt.Println("read error:", err)
			continue
		}
		msg := buf[:n]

		s.broadcast(msg)
	}
}
func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error:", err)
			}
		}(ws)
	}
}
func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/orderbookfeed", websocket.Handler(server.handleWSOrderbook))
	http.ListenAndServe(":4000", nil)
}
