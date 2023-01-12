package helper

import (
	"crypto/rand"
	"unsafe"

	"github.com/fabmation-gmbh/briefkasten-go/internal/log"
	"go.uber.org/zap"
)

// RandString generates a random string with n characters.
// If the generation fails, the function panics to prevent security problems.
func RandString(n uint) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	data := make([]byte, n)

	if _, err := rand.Read(data); err != nil {
		log.Panic("Unable to read random data", zap.Error(err))
	}

	for i, b := range data {
		data[i] = letterBytes[b%byte(len(letterBytes))]
	}

	return *(*string)(unsafe.Pointer(&data))
}
