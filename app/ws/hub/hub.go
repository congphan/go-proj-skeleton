package hub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
	"github.com/richard-xtek/go-grpc-micro-kit/log"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type (
	TopicRegister struct {
		stream string
		client *client
	}

	// keep connections subscribed on specific stream
	Hub struct {
		logger log.Factory

		conns map[*client]bool // connected clients

		// Register requests from the clients.
		register chan *client

		// Unregister requests from clients.
		unregister chan *client

		multicast chan messageWrapper

		topics map[string]map[*client]bool // keep connected clients register on specific topics

		topicRegister chan TopicRegister

		topicUnRegister chan TopicRegister

		tokenValidator *auth.TokenValidator
	}

	client struct {
		logger log.Factory

		hub  *Hub
		m    *sync.Mutex
		conn *websocket.Conn
		sub  map[string]bool // subscribed stream

		// Buffered channel of outbound messages.
		send chan interface{}

		// event will be sent to client in form of combined or not
		combinedStream bool

		userID string
	}

	messageWrapper struct {
		Stream string      `json:"stream"`
		Data   interface{} `json:"data"`
		UserID string      `json:"user_id,omitempty"` // if not empty then only route the message to specific user
	}
)

func New(logger log.Factory, tokenValidator *auth.TokenValidator) *Hub {
	h := &Hub{
		logger,
		map[*client]bool{},
		make(chan *client),
		make(chan *client),
		make(chan messageWrapper),
		map[string]map[*client]bool{},
		make(chan TopicRegister),
		make(chan TopicRegister),
		tokenValidator,
	}

	go h.run()

	return h
}

func newClient(c *websocket.Conn, hub *Hub, combinedStream bool) *client {
	return &client{
		hub.logger,
		hub,
		&sync.Mutex{},
		c,
		map[string]bool{},
		make(chan interface{}, 100),
		combinedStream,
		"",
	}
}

// stream:symbol@<event> eg: btcusdt@aggTrade
func (c *client) Subscribe(streams ...string) {
	c.m.Lock()
	defer c.m.Unlock()

	for _, s := range streams {
		if s != "" {
			c.hub.topicRegister <- TopicRegister{
				s, c,
			}
			c.sub[s] = true
		}
	}
}

func (c *client) subscriptions() []string {
	c.m.Lock()
	defer c.m.Unlock()

	streams := []string{}
	for stream := range c.sub {
		// hide user stream
		if isHiddenStream(stream) {
			continue
		}

		streams = append(streams, stream)
	}

	return streams
}

// stream:symbol@<event> eg: btcusdt@aggTrade
func (c *client) Unsubscribe(streams ...string) {
	c.m.Lock()
	defer c.m.Unlock()

	for _, s := range streams {
		if s != "" {
			c.hub.topicUnRegister <- TopicRegister{
				s, c,
			}
			delete(c.sub, s)
		}
	}
}

func (h *Hub) Add(conn *websocket.Conn, combinedStream bool) *client {
	c := newClient(conn, h, combinedStream)

	go c.readPump()
	go c.writePump()

	h.register <- c

	return c
}

// This is currently support for testing only
func (h *Hub) StartMessage(message interface{}, stream string, targetClientID string) {
	for {
		h.multicast <- messageWrapper{stream, message, targetClientID}
		time.Sleep(time.Second * 5)
	}
}

// if targetUserID is empty then broadcast to all clients listen on this tream
func (h *Hub) Send(message interface{}, stream string, targetUserID string) {
	h.multicast <- messageWrapper{
		Stream: stream,
		Data:   message,
		UserID: targetUserID,
	}
}

func (h *Hub) run() {
	tickerClients := time.NewTicker(time.Second * 10)
	defer tickerClients.Stop()

	for {
		select {
		case <-tickerClients.C:
			h.logger.Bg().Info(fmt.Sprintf("number of connected clients: %d", len(h.conns)))
		case client := <-h.register:
			h.conns[client] = true
		case client := <-h.unregister:
			if _, ok := h.conns[client]; ok {
				// remove all subscribed stream for this client otherwise we got error sending on closed channel
				for _, clients := range h.topics {
					if _, ok := clients[client]; ok {
						delete(clients, client)
					}
				}
				delete(h.conns, client)
				close(client.send)
			}
		case message := <-h.multicast:
			clients, ok := h.topics[message.Stream]
			if !ok {
				break
			}

			for client := range clients {
				// if message specific to a userid then skip for those invalid client
				if message.UserID != "" && message.UserID != client.userID {
					continue
				}

				var data interface{} = message.Data
				if client.combinedStream {
					data = message
				}

				select {
				case client.send <- data:
				default:
					// remove all subscribed stream for this client otherwise we got error sending on closed channel
					for _, cls := range h.topics {
						if _, ok := cls[client]; ok {
							delete(cls, client)
						}
					}
					close(client.send)
					delete(clients, client)
					delete(h.conns, client)
				}
			}

		case reg := <-h.topicRegister:
			if _, ok := h.conns[reg.client]; !ok {
				break
			}
			clients, ok := h.topics[reg.stream]
			if !ok {
				clients = map[*client]bool{}
			}
			clients[reg.client] = true
			h.topics[reg.stream] = clients
		case unreg := <-h.topicUnRegister:
			clients, ok := h.topics[unreg.stream]
			if !ok {
				break
			}
			if _, ok := clients[unreg.client]; ok {
				delete(clients, unreg.client)
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Bg().Error("got unexpected error", zap.Error(err))
			}
			break
		}

		if err := c.handleMessage(messageType, message); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Bg().Error("handleMessage got unexpected error", zap.Error(err))
			}
			break
		}
	}
}

func (c *client) handleMessage(messageType int, p []byte) error {
	comReq := CommonRequest{}
	if err := json.NewDecoder(bytes.NewReader(p)).Decode(&comReq); err != nil {
		return err
	}

	if comReq.Method == methodAuth {
		authReq := AuthRequest{
			Method: comReq.Method,
			Id:     comReq.Id,
		}
		if err := json.NewDecoder(bytes.NewReader(comReq.Params)).Decode(&authReq.Params); err != nil {
			return err
		}

		if authReq.Params.Type != "token" {
			c.send <- ErrorMessage{
				codeInvalidAuth,
				fmt.Sprintf("invalid type: %s", authReq.Params.Type),
			}
			return nil
		}

		principal, err := c.hub.tokenValidator.Validate(authReq.Params.Value)
		if err != nil {
			c.send <- ErrorMessage{
				codeInvalidAuth,
				fmt.Sprintf("invalid token: %s", err.Error()),
			}
			return nil
		}

		// register authenticated user
		c.userID = principal.ID

		// subscribe user event stream
		c.Subscribe(setting.StreamUserEvent)

		c.send <- SubRespone{
			nil,
			authReq.Id,
		}

		return nil
	}

	// nomal request
	req := SubRequest{
		Method: comReq.Method,
		Id:     comReq.Id,
	}
	if err := json.NewDecoder(bytes.NewReader(comReq.Params)).Decode(&req.Params); err != nil {
		return err
	}

	switch true {
	case req.Method == methodSubscribe:
		c.Subscribe(FilterStreamNotAllow(req.Params)...)
		c.send <- SubRespone{
			nil,
			req.Id,
		}
		// write result
	case req.Method == methodUnsubscribe:
		c.Unsubscribe(FilterStreamNotAllow(req.Params)...)
		c.send <- SubRespone{
			nil,
			req.Id,
		}
	case req.Method == methodListSubscriptions:
		c.send <- SubRespone{
			c.subscriptions(),
			req.Id,
		}
	default:
		c.send <- ErrorMessage{
			codeInvalidMethod,
			fmt.Sprintf("method expected one of %s, %s, %s", methodSubscribe, methodUnsubscribe, methodListSubscriptions),
		}
	}

	return nil
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.hub.unregister <- c
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				msg := <-c.send
				if err := c.conn.WriteJSON(msg); err != nil {
					c.logger.Bg().Error("error sending message to client", zap.Error(err))
					return
				}
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.Bg().Error("error sending ping message to client", zap.Error(err))
				return
			}
		}
	}
}

func FilterStreamNotAllow(streams []string) []string {
	out := []string{}

	for _, s := range streams {
		if isHiddenStream(s) {
			continue
		}

		out = append(out, s)
	}

	return out
}

func isHiddenStream(stream string) bool {
	if stream == setting.StreamUserEvent {
		return true
	}

	return false
}
