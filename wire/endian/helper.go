// Copyright (c) 2013-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package endian

import (
	"encoding/binary"
	"io"
)

var Default = binary.BigEndian

// ReadUint8 从 reader 中读取一个 uint8
func ReadUint8(r io.Reader) (uint8, error) {
	var bytes = make([]byte, 1)
	if _, err := io.ReadFull(r, bytes); err != nil {
		return 0, err
	}
	return uint8(bytes[0]), nil
}

// ReadUint32 从 reader 中读取一个 uint32
func ReadUint32(r io.Reader) (uint32, error) {
	var bytes = make([]byte, 4)
	if _, err := io.ReadFull(r, bytes); err != nil {
		return 0, err
	}
	return Default.Uint32(bytes), nil
}

// ReadUint16 从 reader 中读取一个 uint16
func ReadUint16(r io.Reader) (uint16, error) {
	var bytes = make([]byte, 2)
	if _, err := io.ReadFull(r, bytes); err != nil {
		return 0, err
	}
	return Default.Uint16(bytes), nil
}

// ReadUint64 从 reader 中读取一个 uint64
func ReadUint64(r io.Reader) (uint64, error) {
	var bytes = make([]byte, 8)
	if _, err := io.ReadFull(r, bytes); err != nil {
		return 0, err
	}
	return Default.Uint64(bytes), nil
}

// ReadString 从 reader 中读取一个 string
func ReadString(r io.Reader) (string, error) {
	buf, err := ReadBytes(r)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// ReadBytes 从 reader 中读取一个 []byte, reader中前4byte 必须是[]byte 的长度
func ReadBytes(r io.Reader) ([]byte, error) {
	bufLen, err := ReadUint32(r)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, bufLen)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

//ReadFixedBytes 读取固定长度的字节
func ReadFixedBytes(len int, r io.Reader) ([]byte, error) {
	buf := make([]byte, len)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// WriteUint8 写一个 uint8到 writer 中
func WriteUint8(w io.Writer, val uint8) error {
	buf := []byte{byte(val)}
	if _, err := w.Write(buf); err != nil {
		return err
	}
	return nil
}

// WriteUint16 写一个 int16到 writer 中
func WriteUint16(w io.Writer, val uint16) error {
	buf := make([]byte, 2)
	Default.PutUint16(buf, val)
	if _, err := w.Write(buf); err != nil {
		return err
	}
	return nil
}

// WriteUint32 写一个 int32到 writer 中
func WriteUint32(w io.Writer, val uint32) error {
	buf := make([]byte, 4)
	Default.PutUint32(buf, val)
	if _, err := w.Write(buf); err != nil {
		return err
	}
	return nil
}

// WriteUint64 写一个 int64到 writer 中
func WriteUint64(w io.Writer, val uint64) error {
	buf := make([]byte, 8)
	Default.PutUint64(buf, val)
	if _, err := w.Write(buf); err != nil {
		return err
	}
	return nil
}

// WriteString 写一个 string 到 writer 中
func WriteString(w io.Writer, str string) error {
	if err := WriteBytes(w, []byte(str)); err != nil {
		return err
	}
	return nil
}

// WriteBytes 写一个 buf []byte 到 writer 中
func WriteBytes(w io.Writer, buf []byte) error {
	bufLen := len(buf)

	if err := WriteUint32(w, uint32(bufLen)); err != nil {
		return err
	}
	if _, err := w.Write(buf); err != nil {
		return err
	}
	return nil
}

func WriteShortBytes(w io.Writer, buf []byte) error {
	bufLen := len(buf)

	if err := WriteUint16(w, uint16(bufLen)); err != nil {
		return err
	}
	if _, err := w.Write(buf); err != nil {
		return err
	}
	return nil
}

func ReadShortBytes(r io.Reader) ([]byte, error) {
	bufLen, err := ReadUint16(r)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, bufLen)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func ReadShortString(r io.Reader) (string, error) {
	buf, err := ReadShortBytes(r)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
