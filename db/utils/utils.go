package utils

import "encoding/binary"

func I2B(num int64) []byte {
	value := make([]byte, 8)
	binary.BigEndian.PutUint64(value, uint64(num))
	return value
}

func B2I(value []byte) int64 {
	return int64(binary.BigEndian.Uint64(value))
}
