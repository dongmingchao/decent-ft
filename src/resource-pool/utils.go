package resourcePool

import (
	"bytes"
	"encoding/binary"
)

func IntTo3Bytes(i int) [3]byte {
	var buf [3]byte
	u32 := uint32(i)
	buf[2] = uint8(u32)
	buf[1] = uint8(u32 >> 8)
	buf[0] = uint8(u32 >> 16)
	return buf
}

func BytesToUInt32(b []byte) uint32 {
	var buf bytes.Buffer
	for i:=len(b); i < 4; i++ {
		buf.WriteByte(0)
	}
	buf.Write(b)
	var x uint32
	binary.Read(&buf, binary.BigEndian, &x)
	return x
}
