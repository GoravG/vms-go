package ws

import "encoding/binary"

const opText = 0x1

func (c *Conn) WriteText(msg string) error {
	bw := c.rw.Writer
	payload := []byte(msg)
	if err := bw.WriteByte(0x80 | opText); err != nil {
		return err
	}
	l := len(payload)
	switch {
	case l < 126:
		if err := bw.WriteByte(byte(l)); err != nil {
			return err
		}
	case l <= 65535:
		if err := bw.WriteByte(126); err != nil {
			return err
		}
		var ext [2]byte
		binary.BigEndian.PutUint16(ext[:], uint16(l))
		if _, err := bw.Write(ext[:]); err != nil {
			return err
		}
	default:
		if err := bw.WriteByte(127); err != nil {
			return err
		}
	}

	if l > 0 {
		if _, err := bw.Write(payload); err != nil {
			return err
		}
	}
	return bw.Flush()
}
