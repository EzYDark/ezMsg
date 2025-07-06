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

// Define specific types for identifiers to improve code clarity and safety.
type ClientAddr string
type PublicKey string

// A thread-safe map for managing client connections.
type ClientConnectionsList[T ClientAddr | PublicKey] struct {
	sync.RWMutex
	list map[T]quic.Connection
}

// Add safely adds a new client connection to the list.
func (c *ClientConnectionsList[T]) Add(identifier T, conn quic.Connection) {
	c.Lock()
	defer c.Unlock()
	c.list[identifier] = conn
}

// Remove safely removes a client connection from the list.
func (c *ClientConnectionsList[T]) Remove(identifier T) {
	c.Lock()
	defer c.Unlock()
	delete(c.list, identifier)
}

// RemoveByConnection safely removes a client by their connection object.
func (c *ClientConnectionsList[T]) RemoveByConnection(conn quic.Connection) {
	c.Lock()
	defer c.Unlock()
	for id, list_conn := range c.list {
		if list_conn == conn {
			delete(c.list, id)
			break
		}
	}
}

// Global instances of our connection lists with specific types.
var (
	// Connections that are not yet authenticated.
	GuestClients = &ClientConnectionsList[ClientAddr]{
		list: make(map[ClientAddr]quic.Connection),
	}
	// Connections that have successfully authenticated, keyed by their public key.
	AuthenticatedClients = &ClientConnectionsList[PublicKey]{
		list: make(map[PublicKey]quic.Connection),
	}
)

// handleStream processes incoming data from a single client stream.
// This is the heart of the server's application logic.
func handleStream(conn quic.Connection, stream quic.Stream) error {
	log.Debug().Msgf("Client '%s' opened stream.", conn.RemoteAddr())
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

	// =================================================================
	// TODO: Implement Authentication Logic Here
	//
	// If authentication is successful:
	//   1. guestID := GuestIdentifier(conn.RemoteAddr().String())
	//   2. userKey := PublicKey("the_user_public_key")
	//   3. GuestClients.Remove(guestID)
	//   4. AuthenticatedClients.Add(userKey, conn)
	// =================================================================

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
	log.Debug().Msgf("Client connected from '%s'.", conn.RemoteAddr())

	defer func() {
		GuestClients.RemoveByConnection(conn)
		AuthenticatedClients.RemoveByConnection(conn)
		log.Debug().Msgf("Client '%s' disconnected.", conn.RemoteAddr())
	}()

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			return fmt.Errorf("client '%s' closed connection:\n%v", conn.RemoteAddr(), err)
		}
		go handleStream(conn, stream)
	}
}
