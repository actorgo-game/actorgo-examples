package main

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	cfacade "github.com/actorgo-game/actorgo/facade"
	"github.com/actorgo-game/actorgo/net/parser"
	cproto "github.com/actorgo-game/actorgo/net/proto"
	"google.golang.org/protobuf/proto"
)

type agpResponse struct {
	body []byte
	err  error
}

// AGPClient implements the AGP/1 TCP framing used by the demo gate.
type AGPClient struct {
	conn           net.Conn
	codec          *cproto.PacketCodec
	framer         *parser.TCPPacketFramer
	requestTimeout time.Duration

	requestID atomic.Uint32
	writeMu   sync.Mutex
	pendingMu sync.Mutex
	pending   map[uint32]chan agpResponse
	done      chan struct{}
	closeOnce sync.Once
}

func NewAGPClient(requestTimeout time.Duration) *AGPClient {
	if requestTimeout <= 0 {
		requestTimeout = 10 * time.Second
	}
	return &AGPClient{
		codec:          cproto.NewPacketCodec(cproto.DefaultLimits()),
		framer:         parser.NewTCPPacketFramer(cproto.DefaultMaxPacketSize),
		requestTimeout: requestTimeout,
		pending:        make(map[uint32]chan agpResponse),
		done:           make(chan struct{}),
	}
}

func (c *AGPClient) Connect(address string) error {
	conn, err := net.DialTimeout("tcp", address, c.requestTimeout)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.readLoop()

	handshakeBody, err := c.request(
		cproto.SystemMethodHandshake,
		&cproto.HandshakeRequest{SupportedVersions: []uint32{1}},
		cfacade.CodecProtobuf,
	)
	if err != nil {
		c.Close()
		return fmt.Errorf("AGP handshake: %w", err)
	}
	handshake := new(cproto.HandshakeResponse)
	if err = proto.Unmarshal(handshakeBody, handshake); err != nil {
		c.Close()
		return fmt.Errorf("decode AGP handshake: %w", err)
	}
	if handshake.ProtocolVersion != 1 {
		c.Close()
		return fmt.Errorf("unsupported AGP protocol version %d", handshake.ProtocolVersion)
	}
	if handshake.HeartbeatIntervalMs > 0 {
		go c.heartbeatLoop(time.Duration(handshake.HeartbeatIntervalMs) * time.Millisecond)
	}
	return nil
}

func (c *AGPClient) Request(methodID uint32, request, response proto.Message) error {
	body, err := c.request(methodID, request, cfacade.CodecProtobuf)
	if err != nil {
		return err
	}
	if response == nil || len(body) == 0 {
		return nil
	}
	return proto.Unmarshal(body, response)
}

func (c *AGPClient) request(methodID uint32, request proto.Message, codec int32) ([]byte, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("AGP client is not connected")
	}
	body, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	requestID := c.requestID.Add(1)
	if requestID == 0 {
		requestID = c.requestID.Add(1)
	}
	responseCh := make(chan agpResponse, 1)
	c.pendingMu.Lock()
	c.pending[requestID] = responseCh
	c.pendingMu.Unlock()

	packet := &cproto.Packet{
		Kind: &cproto.Packet_Request{Request: &cproto.Request{
			RequestId: requestID,
			MethodId:  methodID,
			TimeoutMs: uint32(c.requestTimeout / time.Millisecond),
			Body:      body,
		}},
		Codec: codec,
	}
	if err = c.writePacket(packet); err != nil {
		c.removePending(requestID)
		return nil, err
	}

	timer := time.NewTimer(c.requestTimeout)
	defer timer.Stop()
	select {
	case response := <-responseCh:
		return response.body, response.err
	case <-timer.C:
		c.removePending(requestID)
		return nil, fmt.Errorf("AGP request timeout: methodID=%d", methodID)
	case <-c.done:
		return nil, fmt.Errorf("AGP connection closed")
	}
}

func (c *AGPClient) writePacket(packet *cproto.Packet) error {
	data, err := c.codec.Encode(packet)
	if err != nil {
		return err
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if err = c.conn.SetWriteDeadline(time.Now().Add(c.requestTimeout)); err != nil {
		return err
	}
	return c.framer.WritePacketBytes(c.conn, data)
}

func (c *AGPClient) readLoop() {
	for {
		data, err := c.framer.ReadPacketBytes(c.conn)
		if err != nil {
			c.Close()
			return
		}
		packet, err := c.codec.Decode(data)
		if err != nil {
			c.Close()
			return
		}
		response := packet.GetResponse()
		if response == nil {
			continue
		}

		responseCh := c.removePending(response.RequestId)
		if responseCh == nil {
			continue
		}
		if response.Code != 0 {
			responseCh <- agpResponse{err: fmt.Errorf("AGP method failed: code=%d message=%s", response.Code, response.Message)}
		} else {
			responseCh <- agpResponse{body: response.Body}
		}
	}
}

func (c *AGPClient) heartbeatLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			_, err := c.request(
				cproto.SystemMethodHeartbeat,
				&cproto.HeartbeatRequest{ClientTimeMs: time.Now().UnixMilli()},
				cfacade.CodecProtobuf,
			)
			if err != nil {
				c.Close()
				return
			}
		case <-c.done:
			return
		}
	}
}

func (c *AGPClient) removePending(requestID uint32) chan agpResponse {
	c.pendingMu.Lock()
	defer c.pendingMu.Unlock()
	responseCh := c.pending[requestID]
	delete(c.pending, requestID)
	return responseCh
}

func (c *AGPClient) Close() {
	c.closeOnce.Do(func() {
		close(c.done)
		if c.conn != nil {
			_ = c.conn.Close()
		}
		c.pendingMu.Lock()
		for requestID, responseCh := range c.pending {
			delete(c.pending, requestID)
			responseCh <- agpResponse{err: fmt.Errorf("AGP connection closed")}
		}
		c.pendingMu.Unlock()
	})
}
