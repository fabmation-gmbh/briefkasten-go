package redis

import (
	"encoding"
	"fmt"

	"github.com/rueian/rueidis"
)

var c rueidis.Client

// Connect connects to redis and initializes the internal cache.
func Connect(client rueidis.Client) error {
	c = client

	return nil
}

// generateRedisValue returns the string representation of the provided value for use as
// a redis value.
func generateRedisValue(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case encoding.BinaryMarshaler:
		if data, err := v.MarshalBinary(); err == nil {
			return rueidis.BinaryString(data)
		}
	case bool:
		if v {
			return "1"
		}
		return "0"
	}

	// TODO: This is the slowest way to get the string representation => we need to report this so we can improve our code
	return fmt.Sprint(value)
}
