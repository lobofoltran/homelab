package stream

type Hub struct {
	clients    map[chan []byte]bool
	register   chan chan []byte
	unregister chan chan []byte
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[chan []byte]bool),
		register:   make(chan chan []byte),
		unregister: make(chan chan []byte),
		broadcast:  make(chan []byte, 10),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			delete(h.clients, client)
			close(client)
		case msg := <-h.broadcast:
			for client := range h.clients {
				select {
				case client <- msg:
				default:
				}
			}
		}
	}
}

func (h *Hub) Broadcast(frame []byte) {
	h.broadcast <- frame
}

func (h *Hub) Register(ch chan []byte) chan []byte {
	h.register <- ch
	return ch
}

func (h *Hub) Unregister(ch chan []byte) {
	h.unregister <- ch
}
