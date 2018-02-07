package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

type String_Util struct {
}

var StringUtil = String_Util{}

func (s *String_Util) MD5Text(inputText string) string {
	hasher := md5.New()
	hasher.Write([]byte(inputText))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (s *String_Util) MD5Byte(data []byte) []byte {
	hasher := md5.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

func StrToInt64(stringValue string) int64 {

	value, err := strconv.ParseInt(stringValue, 10, 64)
	if err != nil {
		return 0
	}
	return value
}

func (s *String_Util) StrToInt(stringValue string) int {

	value, err := strconv.Atoi(stringValue)
	if err != nil {
		return 0
	}
	return value
}

func (s *String_Util) StrToUInt(stringValue string) uint {

	value, err := strconv.Atoi(stringValue)
	if err != nil {
		return 0
	}
	return uint(value)
}

func (s *String_Util) StrToInt32(stringValue string) int32 {
	value, err := strconv.Atoi(stringValue)
	if err != nil {
		return 0
	}
	return int32(value)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func (s *String_Util) Rand(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
