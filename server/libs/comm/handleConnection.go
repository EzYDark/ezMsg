package comm

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	fb "github.com/ezydark/ezMsg/server/flatbuffers/generated/ezMsg/Communication"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/quic-go/quic-go"
	"github.com/rs/zerolog/log"
)

var Delimiter = []byte("\n\r\n\r")

// A thread-safe map for managing client connections.
type ClientConnectionsList struct {
	sync.RWMutex
	list map[string]quic.Connection
}

// Add safely adds a new client connection to the list.
func (c *ClientConnectionsList) Add(sessionToken string, conn quic.Connection) {
	c.Lock() // Get the "write" lock
	defer c.Unlock()
	c.list[sessionToken] = conn
}

// Remove safely removes a client connection.
func (c *ClientConnectionsList) Remove(sessionToken string) {
	c.Lock() // Get the "write" lock
	defer c.Unlock()
	delete(c.list, sessionToken)
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

// Get safely retrieves a connection for a given session token.
func (c *ClientConnectionsList) Get(sessionToken string) (quic.Connection, bool) {
	c.RLock() // Get a "read" lock
	defer c.RUnlock()
	conn, found := c.list[sessionToken]
	return conn, found
}

// Global instance of our connection list.
var ClientsList = &ClientConnectionsList{
	list: make(map[string]quic.Connection),
}

// handleStream processes incoming data from a single client stream.
// This is the heart of the server's application logic.
func handleStream(conn quic.Connection, stream quic.Stream) error {
	log.Debug().Msg("Client opened stream.")
	reader := bufio.NewReader(stream)

	buf, err := reader.ReadBytes(Delimiter[len(Delimiter)-1])
	if err != nil {
		if err == io.EOF {
			return nil // Clean exit
		}
		return fmt.Errorf("error reading from stream:\n%w", err)
	}

	buf = bytes.TrimSuffix(buf, Delimiter)

	log.Debug().Msg("Client frame received successfully.")

	clientFrame := fb.GetRootAsClientFrame(buf, 0)
	payloadTable := new(flatbuffers.Table)
	if !clientFrame.Payload(payloadTable) {
		return fmt.Errorf("failed to get payload from ClientFrame")
	}

	log.Debug().Msg("Client frame parsed successfully.")

	switch clientFrame.PayloadType() {
	case fb.ClientPayloadChatMessageRequest:
		req := new(fb.ChatMessageRequest)
		req.Init(payloadTable.Bytes, payloadTable.Pos)
		return HandleChatMessageRequest(conn, stream, req)
	default:
		return fmt.Errorf("received unknown frame type:\n%v", clientFrame.PayloadType())
	}
}

func HandleConnection(conn quic.Connection) error {
	log.Debug().Msgf("Client connected from %s.", conn.RemoteAddr())

	defer func() {
		ClientsList.RemoveByConnection(conn)
	}()

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			return fmt.Errorf("client %s closed connection:\n%v", conn.RemoteAddr(), err)
		}
		go handleStream(conn, stream)
	}
}
