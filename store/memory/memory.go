package memory

import (
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/w-h-a/pkg/store"
)

type memoryStore struct {
	options store.StoreOptions
	store   *cache.Cache
}

func (s *memoryStore) Options() store.StoreOptions {
	return s.options
}

func (s *memoryStore) Write(rec *store.Record, opts ...store.WriteOption) error {
	// get the key correct
	key := rec.Key
	if len(s.options.Table) > 0 {
		key = s.options.Table + "/" + key
	}
	if len(s.options.Database) > 0 {
		key = s.options.Database + "/" + key
	}

	// copy the incoming record and then convert the expiry to timestamp
	i := &InternalRecord{
		Key: rec.Key,
	}
	i.Value = make([]byte, len(rec.Value))
	copy(i.Value, rec.Value)
	if rec.Expiry != 0 {
		i.ExpiresAt = time.Now().Add(rec.Expiry)
	}

	// set
	s.store.Set(key, i, rec.Expiry)

	return nil
}

func (s *memoryStore) Read(key string, opts ...store.ReadOption) ([]*store.Record, error) {
	options := store.NewReadOptions(opts...)

	keys := []string{key}

	if options.Prefix || options.Suffix {
		opts := []store.ListOption{}

		if options.Prefix {
			opts = append(opts, store.ListWithPrefix(key))
		}

		if options.Suffix {
			opts = append(opts, store.ListWithSuffix(key))
		}

		// TODO: limit and offset
		ks, err := s.List(opts...)
		if err != nil {
			return nil, err
		}

		keys = ks
	}

	records := []*store.Record{}

	if len(keys) == 0 {
		return records, store.ErrRecordNotFound
	}

	for _, k := range keys {
		record, err := s.read(k)
		if err != nil {
			return records, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (s *memoryStore) read(key string) (*store.Record, error) {
	// get the key correct
	if len(s.options.Table) > 0 {
		key = s.options.Table + "/" + key
	}
	if len(s.options.Database) > 0 {
		key = s.options.Database + "/" + key
	}

	// get the record
	r, found := s.store.Get(key)
	if !found {
		return nil, store.ErrRecordNotFound
	}

	// coerce to internal record
	i, ok := r.(*InternalRecord)
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	// copy and return record
	record := &store.Record{
		Key: i.Key,
	}
	record.Value = make([]byte, len(i.Value))
	copy(record.Value, i.Value)
	if !i.ExpiresAt.IsZero() {
		record.Expiry = time.Until(i.ExpiresAt)
	}

	return record, nil
}

func (s *memoryStore) List(opts ...store.ListOption) ([]string, error) {
	options := store.NewListOptions(opts...)

	allKeys := s.list(options.Limit, options.Offset)

	if len(options.Prefix) > 0 {
		prefixKeys := []string{}
		for _, k := range allKeys {
			if strings.HasPrefix(k, options.Prefix) {
				prefixKeys = append(prefixKeys, k)
			}
		}
		allKeys = prefixKeys
	}

	if len(options.Suffix) > 0 {
		suffixKeys := []string{}
		for _, k := range allKeys {
			if strings.HasSuffix(k, options.Suffix) {
				suffixKeys = append(suffixKeys, k)
			}
		}
		allKeys = suffixKeys
	}

	return allKeys, nil
}

func (s *memoryStore) list(_, _ uint) []string {
	allItems := s.store.Items()

	allKeys := make([]string, len(allItems))

	i := 0

	for k := range allItems {
		if len(s.options.Database) > 0 {
			k = strings.TrimPrefix(k, s.options.Database+"/")
		}
		if len(s.options.Table) > 0 {
			k = strings.TrimPrefix(k, s.options.Table+"/")
		}
		allKeys[i] = k
		i++
	}

	// TODO: handle limit and offset

	return allKeys
}

func (s *memoryStore) Delete(key string, opts ...store.DeleteOption) error {
	// get the key correct
	if len(s.options.Table) > 0 {
		key = s.options.Table + "/" + key
	}
	if len(s.options.Database) > 0 {
		key = s.options.Database + "/" + key
	}

	// delete
	s.store.Delete(key)

	return nil
}

func (s *memoryStore) String() string {
	return "memory"
}

func NewStore(opts ...store.StoreOption) store.Store {
	options := store.NewStoreOptions(opts...)

	s := &memoryStore{
		options: options,
		store:   cache.New(cache.NoExpiration, 5*time.Minute),
	}

	return s
}
