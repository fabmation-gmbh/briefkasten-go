package redis

import (
	"encoding/json"
)

// Key represents a cache key prefix.
type Key string

// MarshalBinary implements the Binary marshaller to be compatible with redis.
func (k Key) MarshalBinary() ([]byte, error) {
	return json.Marshal(string(k))
}

const (
	// cacheKeySessionID is the session-id cache key prefix.
	// It resolves a session ID to a user.
	cacheKeySessionID = Key("session:")
)

// SessionID returns the cache [Key] for User ID lookups.
func SessionID(id string) Key {
	return cacheKeySessionID + Key(id)
}
