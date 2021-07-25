package pkt

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/klintcheng/kim/wire"
	"github.com/stretchr/testify/assert"
)

func TestReadPkt(t *testing.T) {
	seq := wire.Seq.Next()

	packet := New("auth.login.aa", WithSeq(seq), WithStatus(Status_Success))
	assert.Equal(t, "auth", packet.ServiceName())
	// assert.Equal(t, "login.aa", packet.CommandPath())

	packet = New("auth.login", WithSeq(seq), WithStatus(Status_Success))
	packet.WriteBody(&LoginReq{
		Token: "test token",
	})
	packet.AddMeta(&Meta{
		Key:   "test",
		Value: "test",
	}, &Meta{
		Key:   wire.MetaDestServer,
		Value: "test",
	}, &Meta{
		Key:   wire.MetaDestChannels,
		Value: "test1,test2",
	})
	buf := new(bytes.Buffer)
	_ = packet.Encode(buf)

	t.Log(buf.Bytes())
	// r := bytes.NewBuffer(Marshal(packet))
	//

	got, err := Read(buf)
	p := got.(*LogicPkt)
	assert.Nil(t, err)
	assert.Equal(t, "auth.login", p.Command)
	assert.Equal(t, seq, p.Sequence)
	assert.Equal(t, Status_Success, p.Status)
	assert.Equal(t, Status_Success, p.Status)

	assert.Equal(t, 3, len(packet.Meta))

	packet.DelMeta(wire.MetaDestServer)
	assert.Equal(t, 2, len(packet.Meta))
	assert.Equal(t, wire.MetaDestChannels, packet.Meta[1].Key)

	packet.DelMeta(wire.MetaDestChannels)
	assert.Equal(t, 1, len(packet.Meta))
}

func Test_Encode(t *testing.T) {
	var pkt = struct {
		Source   uint32
		Sequence uint64
		Data     []byte
	}{
		Source:   0x010201,
		Sequence: 2<<60 + 3,
		Data:     []byte("hello world"),
	}

	// 为了方便观看，使用大端序
	endian := binary.BigEndian

	buf := make([]byte, 1024) // buffer
	i := 0
	endian.PutUint32(buf[i:i+4], pkt.Source)
	i += 4
	endian.PutUint64(buf[i:i+8], pkt.Sequence)
	i += 8
	// 由于data长度不确定，必须先把长度写入buf, 这样在反序列化时就可以正确的解析出data
	dataLen := len(pkt.Data)
	endian.PutUint32(buf[i:i+4], uint32(dataLen))
	i += 4
	// 写入数据data
	copy(buf[i:i+dataLen], pkt.Data)
	i += dataLen
	t.Log(buf[0:i])
	t.Log("length", i)

}

func Test_decode(t *testing.T) {
	var pkt struct {
		Source   uint32
		Sequence uint64
		Data     []byte
	}

	recv := []byte{0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 11, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}
	endian := binary.BigEndian
	i := 0
	pkt.Source = endian.Uint32(recv[i : i+4])
	i += 4
	pkt.Sequence = endian.Uint64(recv[i : i+8])
	i += 8
	dataLen := endian.Uint32(recv[i : i+4])
	i += 4
	pkt.Data = make([]byte, dataLen)
	copy(pkt.Data, recv[i:i+int(dataLen)])
	t.Logf("Src:%d Seq:%d Data:%s", pkt.Source, pkt.Sequence, pkt.Data)
}

// func Test_PktEncode(t *testing.T) {
// 	p := Pkt{
// 		Source:   10000000,
// 		Sequence: 2<<60 + 3,
// 		Data:     []byte("hello world"),
// 	}
// 	bts, err := proto.Marshal(&p)
// 	assert.Nil(t, err)
// 	t.Log(bts)
// 	t.Log("length ", len(bts))

// 	bts, err = json.Marshal(&p)
// 	assert.Nil(t, err)
// 	t.Log(bts)
// 	t.Log(string(bts))
// 	t.Log("length ", len(bts))

// 	t.Log(0x0201)

// }

func Benchmark_Encode(t *testing.B) {
	var pkt = struct {
		Source   uint32
		Sequence uint64
		Data     []byte
	}{
		Source:   10000000,
		Sequence: 2<<60 + 3,
		Data:     []byte("hello world"),
	}
	// t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		// 为了方便观看，使用大端序
		endian := binary.BigEndian

		buf := make([]byte, 30) // buffer
		i := 0
		endian.PutUint32(buf[i:i+4], pkt.Source)
		i += 4
		endian.PutUint64(buf[i:i+8], pkt.Sequence)
		i += 8
		// 由于data长度不确定，必须先把长度写入buf, 这样在反序列化时就可以正确的解析出data
		dataLen := len(pkt.Data)
		endian.PutUint32(buf[i:i+4], uint32(dataLen))
		i += 4
		// 写入数据data
		copy(buf[i:i+dataLen], pkt.Data)
		i += dataLen
	}
}

// func Benchmark_Protobuf(t *testing.B) {
// 	p := Pkt{
// 		Source:   10000000,
// 		Sequence: 2<<60 + 3,
// 		Data:     []byte("hello world"),
// 	}
// 	for i := 0; i < t.N; i++ {
// 		bts, err := proto.Marshal(&p)
// 		assert.Nil(t, err)
// 		assert.NotEmpty(t, bts)
// 	}
// }

// func Benchmark_Json(t *testing.B) {
// 	p := Pkt{
// 		Source:   10000000,
// 		Sequence: 2<<60 + 3,
// 		Data:     []byte("hello world"),
// 	}
// 	for i := 0; i < t.N; i++ {
// 		bts, err := json.Marshal(&p)
// 		assert.Nil(t, err)
// 		assert.NotEmpty(t, bts)
// 	}
// }
