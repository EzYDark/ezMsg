package comm

import (
	"context"
	"io"
	"sync"

	fb "github.com/ezydark/ezMsg/server/flatbuffers/generated/ezMsg/Communication"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/quic-go/quic-go"
	"github.com/rs/zerolog/log"
)

// A thread-safe map for managing client connections.
type ClientConnectionsList struct {
	sync.RWMutex
	list map[string]quic.Connection
}

// Add safely adds a new client connection to the list.
func (c *ClientConnectionsList) Add(userID string, conn quic.Connection) {
	c.Lock() // Get the "write" lock
	defer c.Unlock()
	c.list[userID] = conn
}

// Remove safely removes a client connection.
func (c *ClientConnectionsList) Remove(userID string) {
	c.Lock() // Get the "write" lock
	defer c.Unlock()
	delete(c.list, userID)
}

func (c *ClientConnectionsList) RemoveByConnection(conn quic.Connection) {
	c.Lock() // Get the "write" lock
	defer c.Unlock()
	for id, list_conn := range c.list {
		if list_conn == conn {
			delete(c.list, id)
			break
		}
	}
}

// Remove safely removes all connections.
func (c *ClientConnectionsList) RemoveAll() {
	c.Lock() // Get the "write" lock
	defer c.Unlock()
	c.list = make(map[string]quic.Connection)
}

// Get safely retrieves a connection for a given user ID.
func (c *ClientConnectionsList) Get(userID string) (quic.Connection, bool) {
	c.RLock() // Get a "read" lock
	defer c.RUnlock()
	conn, found := c.list[userID]
	return conn, found
}

// Global instance of our connection list.
var ClientsList = &ClientConnectionsList{
	list: make(map[string]quic.Connection),
}

// handleStream processes incoming data from a single client stream.
// This is the heart of the server's application logic.
func handleStream(conn quic.Connection, stream quic.Stream) {
	buf, err := io.ReadAll(stream)
	if err != nil {
		log.Error().Msgf("Error reading from stream: %v", err)
		return
	}

	clientFrame := fb.GetRootAsClientFrame(buf, 0)
	payloadTable := new(flatbuffers.Table)
	if !clientFrame.Payload(payloadTable) {
		log.Error().Msg("Failed to get payload from ClientFrame")
		return
	}

	switch clientFrame.FrameType() {
	case fb.ClientFrameTypeLoginRequest:
		req := new(fb.LoginRequest)
		req.Init(payloadTable.Bytes, payloadTable.Pos)
		handleLogin(conn, stream, req)
	case fb.ClientFrameTypeSendMessageRequest:
		req := new(fb.SendMessageRequest)
		req.Init(payloadTable.Bytes, payloadTable.Pos)
		handleSendMessage(conn, stream, req)
	default:
		log.Error().Msgf("Received unknown frame type: %v", clientFrame.FrameType())
	}
}

func HandleConnection(conn quic.Connection) {
	log.Debug().Msgf("Client connected from %s.", conn.RemoteAddr())

	defer func() {
		// Clean up connection map on disconnect
		ClientsList.RemoveByConnection(conn)
	}()

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			log.Error().Msgf("Client %s closed connection: %v", conn.RemoteAddr(), err)
			return
		}
		go handleStream(conn, stream)
	}
}
