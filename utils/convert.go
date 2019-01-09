package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Uint128 [2]uint64

func Uint16ToBytes(u uint16) []byte {
	ret := make([]byte, 2)
	binary.BigEndian.PutUint16(ret, u)
	return ret
}

func Uint32ToBytes(u uint32) []byte {
	ret := make([]byte, 4)
	binary.BigEndian.PutUint32(ret, u)
	return ret
}

func Uint128ToBytes(u Uint128) []byte {
	part1 := make([]byte, 8)
	part2 := make([]byte, 8)
	binary.BigEndian.PutUint64(part1, uint64(u[0]))
	binary.BigEndian.PutUint64(part2, uint64(u[1]))
	return append(part1, part2...)
}

func BytesToUint16(b []byte) uint16 {
	var ret uint16
	if err := binary.Read(bytes.NewBuffer(b), binary.BigEndian, &ret); err != nil {
		fmt.Println(err.Error())
		ret = 0
	}
	return ret
}

func BytesToInt(b []byte) int {
	var ret int
	if err := binary.Read(bytes.NewBuffer(b), binary.BigEndian, &ret); err != nil {
		fmt.Println(err.Error())
		ret = 0
	}
	return ret
}
