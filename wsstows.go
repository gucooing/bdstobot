package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

func main() {
	upstreamURL := "wss://gateway.discord.gg" // 上游（wss）URL
	downstreamPath := "/"               // 下游（ws）路径

	u, err := url.Parse(upstreamURL)
	if err != nil {
		log.Fatal(err)
	}

	proxy := NewProxy(u)

	http.HandleFunc(downstreamPath, proxy.Handle)

	log.Println("WebSocket proxy server started on :9966")
	err = http.ListenAndServe(":9966", nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Proxy struct {
	upstreamURL *url.URL
}

func NewProxy(upstreamURL *url.URL) *Proxy {
	return &Proxy{
		upstreamURL: upstreamURL,
	}
}

func (p *Proxy) Handle(w http.ResponseWriter, r *http.Request) {
	upstreamConn, _, err := websocket.DefaultDialer.Dial(p.upstreamURL.String(), nil)
	if err != nil {
		log.Println("failed to connect to upstream:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer upstreamConn.Close()

	downstreamConn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println("failed to upgrade downstream connection:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer downstreamConn.Close()

	go p.proxyMessages(upstreamConn, downstreamConn)
	p.proxyMessages(downstreamConn, upstreamConn)
}

func (p *Proxy) proxyMessages(src *websocket.Conn, dst *websocket.Conn) {
	for {
		messageType, message, err := src.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure) {
				log.Println("failed to read message:", err)
			}
			break
		}

		err = dst.WriteMessage(messageType, message)
		if err != nil {
			log.Println("failed to write message:", err)
			break
		}
	}

	src.Close()
	dst.Close()
}
