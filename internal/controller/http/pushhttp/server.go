package pushhttp

import (
	"log"
	"net/http"

	"github.com/cr00z/gocqrs/internal/controller/websock"
)

type httpServer struct {
	hub *websock.Hub
}

func NewHttpServer(hub *websock.Hub) *httpServer {
	server := &httpServer{
		hub: hub,
	}

	http.HandleFunc("/pusher", server.webSocketHandler)

	return server
}

func (s *httpServer) Start() error {
	log.Println("push server started at :8080")
	return http.ListenAndServe(":8080", nil)
}

func (s *httpServer) webSocketHandler(w http.ResponseWriter, r *http.Request) {
	socket, err := websock.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not upgrade", http.StatusInternalServerError)
		return
	}

	client := websock.NewClient(s.hub, socket)
	s.hub.Register(client)

	go client.Write()
}
