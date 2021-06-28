// Copyright (c) 2013-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package endian

import (
	"encoding/binary"
	"testing"
)

func TestReadUint32(t *testing.T) {
	a := uint32(0x01020304)
	arr := make([]byte, 4)
	binary.BigEndian.PutUint32(arr, a)
	t.Log(arr) //[1 2 3 4]

	binary.LittleEndian.PutUint32(arr, a)
	t.Log(arr) //[4 3 2 1]
}

func TestSerial(t *testing.T) {
	var pkt = struct {
		Source         uint16
		Destination    uint16
		Sequence       uint32
		Acknowledgment uint32 //
		Data           []byte
	}{
		Source:         4000,
		Destination:    80,
		Sequence:       100,
		Acknowledgment: 1,
		Data:           []byte("hello world"),
	}

	// 为了方便观看，使用大端序
	endian := binary.BigEndian

	buf := make([]byte, 1024) // buffer
	i := 0
	endian.PutUint16(buf[i:i+2], pkt.Source)
	i += 2 // 移动指针2个字节
	endian.PutUint16(buf[i:i+2], pkt.Destination)
	i += 2
	endian.PutUint32(buf[i:i+4], pkt.Sequence)
	i += 4
	endian.PutUint32(buf[i:i+4], pkt.Acknowledgment)
	i += 4
	// 由于data长度不确定，必须先把长度写入buf, 这样在反序列化时就可以正确的解析出data
	dataLen := len(pkt.Data)
	endian.PutUint32(buf[i:i+4], uint32(dataLen))
	i += 4
	// 写入数据data
	copy(buf[i:i+dataLen], pkt.Data)
	i += dataLen
	t.Log(buf[0:i])
}

func TestDecode(t *testing.T) {
	var pkt struct {
		Source         uint16
		Destination    uint16
		Sequence       uint32
		Acknowledgment uint32 //
		Data           []byte
	}

	recv := []byte{15, 160, 0, 80, 0, 0, 0, 100, 0, 0, 0, 1, 0, 0, 0, 11, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}
	// 为了方便观看，使用大端序
	endian := binary.BigEndian
	i := 0
	pkt.Source = endian.Uint16(recv[i : i+2])
	i += 2 // 移动指针2个字节
	pkt.Destination = endian.Uint16(recv[i : i+2])
	i += 2
	pkt.Sequence = endian.Uint32(recv[i : i+4])
	i += 4
	pkt.Acknowledgment = endian.Uint32(recv[i : i+4])
	i += 4
	dataLen := endian.Uint32(recv[i : i+4])
	i += 4
	pkt.Data = make([]byte, dataLen)
	copy(pkt.Data, recv[i:i+int(dataLen)])
	t.Logf("Src:%d Dest:%d Seq:%d Ack:%d Data:%s", pkt.Source, pkt.Destination, pkt.Sequence, pkt.Acknowledgment, pkt.Data)
	// t.Log(pkt)
}
