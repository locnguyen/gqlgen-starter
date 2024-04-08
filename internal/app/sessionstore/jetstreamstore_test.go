package sessionstore

import (
	"github.com/nats-io/nats.go/jetstream"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	js := MockJetstream{}
	js.On("CreateOrUpdateKeyValue", mock.Anything, mock.Anything).
		Return(MockJetstreamKeyValue{}, nil)

	sessionStore, err := New(js)
	assert.NoError(t, err)
	assert.NotNil(t, sessionStore)
}

func TestNewWithError(t *testing.T) {
	js := MockJetstream{}
	js.On("CreateOrUpdateKeyValue", mock.Anything, mock.Anything).
		Return(nil, jetstream.ErrCantGetBucket)

	sessionStore, err := New(js)
	assert.ErrorContains(t, err, "getting bucket with SCS sessions")
	assert.Nil(t, sessionStore)
}

func TestCommit(t *testing.T) {
	js := MockJetstream{}
	kvs := MockJetstreamKeyValue{}
	fakeRevision := uint64(1)
	kvs.On("Put", mock.Anything, mock.Anything, mock.Anything).
		Return(fakeRevision, nil)
	js.On("CreateOrUpdateKeyValue", mock.Anything, mock.Anything).
		Return(kvs, nil)

	sessionStore, err := New(js)
	assert.NoError(t, err)
	err = sessionStore.Commit(ulid.Make().String(), ulid.Make().Bytes(), time.Now().Add(time.Hour))
	assert.NoError(t, err, "error should be nil if session data is committed")
}

func TestCommitWithError(t *testing.T) {
	js := MockJetstream{}
	kvs := MockJetstreamKeyValue{}
	kvs.On("Put", mock.Anything, mock.Anything, mock.Anything).
		Return(uint64(0), jetstream.ErrBadRequest)
	js.On("CreateOrUpdateKeyValue", mock.Anything, mock.Anything).
		Return(kvs, nil)

	sessionStore, err := New(js)
	assert.NoError(t, err)
	err = sessionStore.Commit(ulid.Make().String(), ulid.Make().Bytes(), time.Now().Add(time.Hour))
	assert.ErrorContains(t, err, "persisting session data by token")
}

func TestDelete(t *testing.T) {
	js := MockJetstream{}
	kvs := MockJetstreamKeyValue{}
	kvs.On("Delete", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	js.On("CreateOrUpdateKeyValue", mock.Anything, mock.Anything).
		Return(kvs, nil)

	sessionStore, err := New(js)
	assert.NoError(t, err)
	err = sessionStore.Delete(ulid.Make().String())
	assert.NoError(t, err)
}

func TestDeleteWithError(t *testing.T) {
	js := MockJetstream{}
	kvs := MockJetstreamKeyValue{}
	kvs.On("Delete", mock.Anything, mock.Anything, mock.Anything).
		Return(jetstream.ErrBadRequest)
	js.On("CreateOrUpdateKeyValue", mock.Anything, mock.Anything).
		Return(kvs, nil)

	sessionStore, err := New(js)
	assert.NoError(t, err)
	err = sessionStore.Delete(ulid.Make().String())
	assert.ErrorContains(t, err, "deleting session data by token")
}

func TestFind(t *testing.T) {
	storeData := ulid.Make()
	kve := MockJetstreamKeyValueEntry{}
	kve.On("Value").
		Return(storeData.Bytes())

	kvs := MockJetstreamKeyValue{}
	kvs.On("Get", mock.Anything, mock.Anything).
		Return(kve, nil)

	js := MockJetstream{}
	js.On("CreateOrUpdateKeyValue", mock.Anything, mock.Anything).
		Return(kvs, nil)

	sessionStore, err := New(js)
	assert.NoError(t, err)

	data, found, err := sessionStore.Find(storeData.String())
	assert.NoError(t, err)
	assert.True(t, found)
	assert.NotNil(t, data)
}
