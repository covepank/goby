package store

import (
	"fmt"

	"github.com/sanbsy/goby/internal/idx"
	bolt "go.etcd.io/bbolt"
)

type (
	BoltStorage struct {
		dbPath        string
		db            *bolt.DB
		defaultBucket *BoltBucket
	}
	BoltBucket struct {
		*BoltStorage
		name []byte
	}
)

func NewBolt(dbPath string) (*BoltStorage, error) {
	if dbPath == "" {
		dbPath = "temp.db"
	}
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	b := &BoltStorage{
		dbPath: dbPath,
		db:     db,
	}
	defaultBucket, err := b.NewBucket("default")
	if err != nil {
		return nil, err
	}
	b.defaultBucket = defaultBucket
	return b, nil
}

func (b *BoltStorage) Set(key, value []byte) error {
	return b.defaultBucket.Set(key, value)
}

func (b *BoltStorage) Get(key []byte) ([]byte, error) {
	return b.defaultBucket.Get(key)
}

func (b *BoltStorage) Delete(key []byte) error {
	return b.defaultBucket.Delete(key)
}
func (b *BoltStorage) All() (map[string][]byte, error) {
	return b.defaultBucket.All()
}

func (b *BoltStorage) BatchUpdate(key, value []byte) error {
	return b.defaultBucket.BatchUpdate(key, value)
}

func (b *BoltStorage) Close() error {
	return b.db.Close()
}

// Sync  数据刷盘
func (b *BoltStorage) Sync() error {
	return b.db.Sync()
}

func (b *BoltStorage) Write(data []byte) (int, error) {
	err := b.Set([]byte(idx.NewID().String()), data)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

// NewBucket 创建一个新的 Bucket
func (b *BoltStorage) NewBucket(name string) (*BoltBucket, error) {
	err := b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &BoltBucket{
		BoltStorage: b,
		name:        []byte(name),
	}, nil
}

// DeleteBucket 删除指定 Bucket
func (b *BoltStorage) DeleteBucket(name string) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(name))
		if err != nil {
			return fmt.Errorf("delete bucket: %v", err)
		}
		return nil
	})
	return err
}

// Set 在指定 Bucket 中存储数据
func (bb *BoltBucket) Set(key, value []byte) error {
	return bb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bb.name)
		if b == nil {
			return fmt.Errorf("this bucket does not exist: %s", string(bb.name))
		}
		return b.Put(key, value)
	})
}

// Get 在指定 Bucket 中查询相关值
func (bb *BoltBucket) Get(key []byte) ([]byte, error) {
	var value []byte
	err := bb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bb.name)
		if b == nil {
			return fmt.Errorf("this bucket does not exist: %s", string(bb.name))
		}
		value = b.Get(key)
		if value == nil {
			return fmt.Errorf("this key does not exist: %s", string(key))
		}
		return nil
	})
	return value, err
}

// Delete 删除指定 KEY-VALUE 数据
func (bb *BoltBucket) Delete(key []byte) error {
	return bb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bb.name)
		if b == nil {
			return fmt.Errorf("this bucket does not exist: %s", string(bb.name))
		}
		return b.Delete(key)
	})
}

// All 获取指定 Bucket 中所有数据
func (bb *BoltBucket) All() (map[string][]byte, error) {
	result := make(map[string][]byte)
	err := bb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bb.name)
		if b == nil {
			return fmt.Errorf("this bucket does not exist: %s", string(bb.name))
		}
		err := b.ForEach(func(k, v []byte) error {
			result[string(k)] = v
			return nil
		})
		return err
	})
	return result, err
}

// Update 数据， 当使用多线程/协程处理时，用此方法
func (bb *BoltBucket) BatchUpdate(key, value []byte) error {
	return bb.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bb.name)
		if b == nil {
			return fmt.Errorf("this bucket does not exist: %s", string(bb.name))
		}
		return b.Put(key, value)
	})
}
