package apiv1

import (
	"reflect"
	"time"
	"unsafe"

	"github.com/fabmation-gmbh/briefkasten-go/internal/config"
	"github.com/gofiber/fiber/v2"
)

// storeOAuthCookie stores a cookie used for OAuth2 authentication.
func storeOAuthCookie(c *fiber.Ctx, name, value string) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   config.C.General.SecureCookie,
		HTTPOnly: true,
	})
}

// byteSlice2String converts a byte slice to a string in a performant way.
func byteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// String2bytes converts the given string to a byte slice without memory allocation.
//
// Note it may break if string and/or slice header will change in future go versions.
func String2bytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len

	return b
}
