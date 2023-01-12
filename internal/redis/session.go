package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fabmation-gmbh/briefkasten-go/internal/log"
	"github.com/markbates/goth"
	"github.com/pkg/errors"
	"github.com/rueian/rueidis"
	"go.uber.org/zap"
)

// UserSession is wrapper arround [goth.Session].
type UserSession struct {
	Session  goth.Session
	Provider string

	Prov goth.Provider
}

// MarshalBinary implements the Binary marshaller to be compatible with redis.
func (u UserSession) MarshalBinary() (data []byte, err error) {
	type rawUserSession struct {
		Session  string
		Provider string
	}

	ses := rawUserSession{
		Session:  u.Session.Marshal(),
		Provider: u.Provider,
	}

	bytes, err := json.Marshal(ses)
	return bytes, err
}

// UnmarshalBinary implements the Binary unmarshaller to be compatible with redis.
func (u *UserSession) UnmarshalBinary(data []byte) error {
	type rawUserSession struct {
		Session  string
		Provider string
	}

	var us rawUserSession

	if err := json.Unmarshal(data, &us); err != nil {
		return errors.Wrap(err, "unable to unmarshal raw session")
	}

	ses, err := u.Prov.UnmarshalSession(us.Session)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal session")
	}

	u.Session = ses
	u.Provider = us.Provider

	return nil
}

// StoreUserSession stores the user session.
func StoreUserSession(ctx context.Context, sesID string, ses UserSession) error {
	const sessionTTL = 24 * time.Hour

	return ESetTTL(ctx, SessionID(sesID), ses, sessionTTL)
}

// GetUserSession returns the stored user session.
func GetUserSession(ctx context.Context, sesID string, provider goth.Provider) (UserSession, bool) {
	const ttl = 24 * time.Hour

	cmd := c.B().
		Get().
		Key(string(SessionID(sesID))).
		Cache()

	u := UserSession{
		Prov: provider,
	}

	err := c.DoCache(ctx, cmd, ttl).DecodeJSON(&u)
	if err != nil {
		if !rueidis.IsRedisNil(err) {
			log.Error("Unable to retrieve user session from redis", zap.Error(err))
		}

		return u, false
	}

	return u, true
}

// DeleteUserSession deletes the user session.
func DeleteUserSession(ctx context.Context, sesID string) {
	const sessionTTL = 24 * time.Hour

	err := EDel(ctx, SessionID(sesID))
	if err != nil {
		log.Error("An error occurred while deleting a session key from redis", zap.String("session_id", sesID))
	}
}

// EDel deletes the key-value pair.
func EDel(ctx context.Context, key Key) error {
	cmd := c.B().Del().
		Key(string(key)).
		Build()

	log.Debug("Deleting redis key", zap.String("key", string(key)))

	err := c.Do(ctx, cmd).Error()
	if err != nil && rueidis.IsRedisNil(err) {
		err = nil
	}

	return errors.Wrap(err, "unable to delete key")
}

// ESetTTL stores the provided key-value pair with a TTL.
func ESetTTL(ctx context.Context, key Key, val any, ttl time.Duration) error {
	cmds := make(rueidis.Commands, 0, 2)

	cmds = append(cmds,
		c.B().Set().
			Key(string(key)).
			Value(generateRedisValue(val)).
			Build(),
		c.B().Expire().
			Key(string(key)).
			Seconds(int64(ttl.Seconds())).
			Nx(). // NOTE: 'NX' is cruical to prevent that the tags will be deleted BEFORE this key is removed
			Build(),
	)

	for _, resp := range c.DoMulti(ctx, cmds...) {
		if err := resp.Error(); err != nil {
			return errors.Wrap(err, "unable to set cache entry")
		}
	}

	return nil
}
