package ws

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

const wsGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

var ErrBadRequest = errors.New("bad websocket upgrade request")

type Conn struct {
	conn net.Conn
	rw   *bufio.ReadWriter
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	if !strings.EqualFold(r.Header.Get("Connection"), "Upgrade") || !strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		return nil, ErrBadRequest
	}
	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		return nil, ErrBadRequest
	}

	h := sha1.New()
	_, _ = io.WriteString(h, key+wsGUID)

	accept := base64.StdEncoding.EncodeToString(h.Sum(nil))

	hj, ok := w.(http.Hijacker)

	if !ok {
		return nil, errors.New("server does not support hijacking")
	}

	netc, buf, err := hj.Hijack()

	if err != nil {
		return nil, err
	}
	resp := fmt.Sprintf(
		"HTTP/1.1 101 Switching Protocols\r\n"+
			"Upgrade: websocket\r\n"+
			"Connection: Upgrade\r\n"+
			"Sec-WebSocket-Accept: %s\r\n\r\n", accept)
	if _, err := buf.WriteString(resp); err != nil {
		_ = netc.Close()
		return nil, err
	}
	if err := buf.Flush(); err != nil {
		_ = netc.Close()
		return nil, err
	}
	return &Conn{
		conn: netc,
		rw:   bufio.NewReadWriter(buf.Reader, buf.Writer),
	}, nil
}
