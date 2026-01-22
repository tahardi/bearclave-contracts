package foundry

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const (
	HexBase         = 16
	HexStringPrefix = "0x"
	NumSize         = 64
)

func BytesToHexString(b []byte) string {
	return HexStringPrefix + strings.ToLower(hex.EncodeToString(b))
}

func Uint64ToHexString(n uint64) string {
	return HexStringPrefix + fmt.Sprintf("%x", n)
}

func ParseBytesFromHexString(s string) ([]byte, error) {
	return hex.DecodeString(
		strings.ToLower(
			strings.TrimPrefix(s, HexStringPrefix),
		),
	)
}

func ParseUint64FromHexString(s string) (uint64, error) {
	return strconv.ParseUint(
		strings.ToLower(strings.TrimPrefix(s, HexStringPrefix)),
		HexBase,
		NumSize,
	)
}
