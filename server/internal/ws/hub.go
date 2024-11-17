package ws

type Room struct {
	Clients map[string]*Client `json:"clients"`
	ID      string             `json:"id"`
	Name    string             `json:"name"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			r, roomExists := h.Rooms[cl.RoomID]

			if !roomExists {
				return
			}

			if _, clientExist := r.Clients[cl.ID]; clientExist {
				return
			}

			r.Clients[cl.ID] = cl
		case cl := <-h.Unregister:
			r, roomExists := h.Rooms[cl.RoomID]

			if !roomExists {
				return
			}

			if _, clientExist := r.Clients[cl.ID]; !clientExist {
				return
			}

			if len(r.Clients) != 0 {
				h.Broadcast <- &Message{
					Content:  "user left the chat",
					RoomID:   cl.RoomID,
					Username: cl.Username,
				}
			}

			delete(r.Clients, cl.ID)
			close(cl.Message)
		case m := <-h.Broadcast:
			r, roomExists := h.Rooms[m.RoomID]

			if !roomExists {
				return
			}

			for _, cl := range r.Clients {
				cl.Message <- m
			}
		}
	}
}
