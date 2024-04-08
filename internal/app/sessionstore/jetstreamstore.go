package sessionstore

import (
	"context"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
	"time"
)

// JetStreamSessionStore represents the session store. We can implement the
// [Store interface](https://github.com/alexedwards/scs?tab=readme-ov-file#using-custom-session-stores) using the APIs
// provided by [JetStream](https://docs.nats.io/nats-concepts/jetstream/key-value-store/kv_walkthrough).
type JetStreamSessionStore struct {
	bucket jetstream.KeyValue
}

// New returns a JetStreamSessionStore instance with a bucket containing SCS sessions.
func New(js jetstream.JetStream) (*JetStreamSessionStore, error) {
	bucket, err := js.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket:      "scs-sessions",
		Description: "bucket for persisting SCS sessions",
		TTL:         24 * time.Hour,
	})

	if err != nil {
		return nil, errors.Wrap(err, "getting bucket with SCS sessions")
	}
	return &JetStreamSessionStore{bucket}, nil
}

// Delete should remove the session token and corresponding data from the
// session store. If the token does not exist then Delete should be a no-op
// and return nil (not an error).
// Reference: https://docs.nats.io/using-nats/developer/develop_jetstream/kv#deleting
func (j *JetStreamSessionStore) Delete(token string) error {
	err := j.bucket.Delete(context.Background(), token, nil)
	if err != nil {
		return errors.Wrap(err, "deleting session data by token")
	}
	return nil
}

// Find should return the data for a session token from the store. If the
// session token is not found or is expired, the found return value should
// be false (and the err return value should be nil). Similarly, tampered
// or malformed tokens should result in a found return value of false and a
// nil err value. The err return value should be used for system errors only.
func (j *JetStreamSessionStore) Find(token string) (b []byte, found bool, err error) {
	kve, err := j.bucket.Get(context.Background(), token)
	if err != nil {
		return []byte{}, false, err
	}
	return kve.Value(), true, nil
}

// Commit should add the session token and data to the store, with the given
// expiry time. If the session token already exists, then the data and
// expiry time should be overwritten.
func (j *JetStreamSessionStore) Commit(token string, b []byte, expiry time.Time) error {
	_, err := j.bucket.Put(context.Background(), token, b)
	if err != nil {
		return errors.Wrap(err, "persisting session data by token")
	}
	return nil
}

// All should return a map containing data for all active sessions (i.e.
// sessions which have not expired). The map key should be the session
// token and the map value should be the session data. If no active
// sessions exist this should return an empty (not nil) map.
func (j *JetStreamSessionStore) All() (map[string][]byte, error) {
	keys, err := j.bucket.Keys(context.Background())
	if err != nil {
		return nil, err
	}
	sessions := map[string][]byte{}

	for _, k := range keys {
		val, err := j.bucket.Get(context.Background(), k)
		if err != nil {
			return nil, err
		}
		sessions[k] = val.Value()
	}
	return sessions, nil
}
