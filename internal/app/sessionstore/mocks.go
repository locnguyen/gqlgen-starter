package sessionstore

import (
	"context"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/mock"
)

type MockJetstream struct {
	jetstream.JetStream
	mock.Mock
}

func (js MockJetstream) CreateOrUpdateKeyValue(ctx context.Context, cfg jetstream.KeyValueConfig) (jetstream.KeyValue, error) {
	args := js.Called(ctx, cfg)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(jetstream.KeyValue), args.Error(1)
}

type MockJetstreamKeyValue struct {
	jetstream.KeyValue
	mock.Mock
}

type MockJetstreamKeyValueEntry struct {
	jetstream.KeyValueEntry
	mock.Mock
}

func (jkv MockJetstreamKeyValue) Get(ctx context.Context, key string) (jetstream.KeyValueEntry, error) {
	args := jkv.Called(ctx, key)
	return args.Get(0).(jetstream.KeyValueEntry), nil
}

//	func (jkv *MockJetstreamKeyValue) GetRevision(ctx context.Context, key string, revision uint64) (jetstream.KeyValueEntry, error) {
//		// TODO implement me
//		panic("implement me")
//	}
func (jkv MockJetstreamKeyValue) Put(ctx context.Context, key string, value []byte) (uint64, error) {
	args := jkv.Called(ctx, key, value)
	return args.Get(0).(uint64), args.Error(1)
}

//	func (jkv *MockJetstreamKeyValue) PutString(ctx context.Context, key string, value string) (uint64, error) {
//		// TODO implement me
//		panic("implement me")
//	}
//
//	func (jkv *MockJetstreamKeyValue) Create(ctx context.Context, key string, value []byte) (uint64, error) {
//		// TODO implement me
//		panic("implement me")
//	}
//
//	func (jkv *MockJetstreamKeyValue) Update(ctx context.Context, key string, value []byte, revision uint64) (uint64, error) {
//		// TODO implement me
//		panic("implement me")
//	}
func (jkv MockJetstreamKeyValue) Delete(ctx context.Context, key string, opts ...jetstream.KVDeleteOpt) error {
	args := jkv.Called(ctx, key, opts)
	return args.Error(0)
}

//
// func (jkv *MockJetstreamKeyValue) Purge(ctx context.Context, key string, opts ...jetstream.KVDeleteOpt) error {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (jkv *MockJetstreamKeyValue) Watch(ctx context.Context, keys string, opts ...jetstream.WatchOpt) (jetstream.KeyWatcher, error) {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (jkv *MockJetstreamKeyValue) WatchAll(ctx context.Context, opts ...jetstream.WatchOpt) (jetstream.KeyWatcher, error) {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (jkv *MockJetstreamKeyValue) Keys(ctx context.Context, opts ...jetstream.WatchOpt) ([]string, error) {
// 	args := jkv.Called(ctx, opts)
// 	return args.Get(0).([]string), nil
// }
//
// func (jkv *MockJetstreamKeyValue) ListKeys(ctx context.Context, opts ...jetstream.WatchOpt) (jetstream.KeyLister, error) {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (jkv *MockJetstreamKeyValue) History(ctx context.Context, key string, opts ...jetstream.WatchOpt) ([]jetstream.KeyValueEntry, error) {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (jkv *MockJetstreamKeyValue) Bucket() string {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (jkv *MockJetstreamKeyValue) PurgeDeletes(ctx context.Context, opts ...jetstream.KVPurgeOpt) error {
// 	// TODO implement me
// 	panic("implement me")
// }
//
// func (jkv *MockJetstreamKeyValue) Status(ctx context.Context) (jetstream.KeyValueStatus, error) {
// 	// TODO implement me
// 	panic("implement me")
// }
