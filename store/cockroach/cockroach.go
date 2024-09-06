package cockroach

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/lib/pq"
	"github.com/w-h-a/pkg/store"
	"github.com/w-h-a/pkg/telemetry/log"
)

type cockroachStore struct {
	options    store.StoreOptions
	client     *sql.DB
	write      *sql.Stmt
	readOne    *sql.Stmt
	readMany   *sql.Stmt
	readOffset *sql.Stmt
	list       *sql.Stmt
	delete     *sql.Stmt
}

func (s *cockroachStore) Options() store.StoreOptions {
	return s.options
}

// TODO: retry
/* for loop with max retries
**	go call a function that (a) calls exponential function (b) sleeps for the duration and (c) calls actual write that sends error to chan
**  if err from chan is nil, stop
**  else, keep looping
 */
func (s *cockroachStore) Write(rec *store.Record, opts ...store.WriteOption) error {
	var err error

	if rec.Expiry != 0 {
		_, err = s.write.Exec(rec.Key, rec.Value, time.Now().Add(rec.Expiry))
	} else {
		_, err = s.write.Exec(rec.Key, rec.Value, nil)
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *cockroachStore) Read(key string, opts ...store.ReadOption) ([]*store.Record, error) {
	options := store.NewReadOptions(opts...)

	// read many; otherwise, read one
	if options.Prefix || options.Suffix {
		return s.read(key, options)
	}

	records := []*store.Record{}

	var timehelper pq.NullTime

	row := s.readOne.QueryRow(key)

	record := &store.Record{}

	if err := row.Scan(&record.Key, &record.Value, &timehelper); err != nil {
		if err == sql.ErrNoRows {
			return records, store.ErrRecordNotFound
		}
		return nil, err
	}

	// if the expiry is valid, we'll check if it has expired
	// otherwise, we default to appending
	if timehelper.Valid {
		// if the record has expired, then delete it instead
		// otherwise, store the expiry on the record and append
		if timehelper.Time.Before(time.Now()) {
			go s.Delete(record.Key)
			return records, store.ErrRecordNotFound
		}
		record.Expiry = time.Until(timehelper.Time)
		records = append(records, record)
	} else {
		records = append(records, record)
	}

	return records, nil
}

func (s *cockroachStore) read(key string, options store.ReadOptions) ([]*store.Record, error) {
	pattern := "%"

	if options.Prefix {
		pattern = key + pattern
	}

	if options.Suffix {
		pattern = pattern + key
	}

	records := []*store.Record{}

	var rows *sql.Rows

	var err error

	if options.Limit != 0 {
		rows, err = s.readOffset.Query(pattern, options.Limit, options.Offset)
	} else {
		rows, err = s.readMany.Query(pattern)
	}

	if err != nil {
		return nil, err
	}

	var timehelper pq.NullTime

	for rows.Next() {
		record := &store.Record{}

		if err := rows.Scan(&record.Key, &record.Value, &timehelper); err != nil {
			return records, err
		}

		// if the expiry is valid, we'll check if it has expired
		// otherwise, we default to appending
		if timehelper.Valid {
			// if the record has expired, then delete it instead
			// otherwise, append
			if timehelper.Time.Before(time.Now()) {
				go s.Delete(record.Key)
			} else {
				record.Expiry = time.Until(timehelper.Time)
				records = append(records, record)
			}
		} else {
			records = append(records, record)
		}
	}

	// TODO: better cleanup needed?
	closeErr := rows.Close()
	if closeErr != nil {
		return records, closeErr
	}

	rowsErr := rows.Err()
	if rowsErr != nil {
		return records, rowsErr
	}

	return records, nil
}

func (s *cockroachStore) List(opts ...store.ListOption) ([]string, error) {
	keys := []string{}

	var timehelper pq.NullTime

	rows, err := s.list.Query()
	if err != nil {
		if err == sql.ErrNoRows {
			return keys, nil
		}
		return nil, err
	}

	for rows.Next() {
		record := &store.Record{}

		if err := rows.Scan(&record.Key, &record.Value, &timehelper); err != nil {
			return keys, err
		}

		// if the expiry is valid, we'll check if it has expired
		// otherwise, we default to appending
		if timehelper.Valid {
			// if the record has expired, then delete it instead
			// otherwise, append
			if timehelper.Time.Before(time.Now()) {
				go s.Delete(record.Key)
			} else {
				keys = append(keys, record.Key)
			}
		} else {
			keys = append(keys, record.Key)
		}
	}

	// TODO: better cleanup needed?
	closeErr := rows.Close()
	if closeErr != nil {
		return keys, closeErr
	}

	rowsErr := rows.Err()
	if rowsErr != nil {
		return keys, rowsErr
	}

	return keys, nil
}

func (s *cockroachStore) Delete(key string, opts ...store.DeleteOption) error {
	if _, err := s.delete.Exec(key); err != nil {
		return err
	}

	return nil
}

func (s *cockroachStore) String() string {
	return "cockroach"
}

func (s *cockroachStore) configure() error {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return errors.New("failed to compile regex for database and table names")
	}
	s.options.Database = reg.ReplaceAllString(s.options.Database, "_")
	s.options.Table = reg.ReplaceAllString(s.options.Table, "_")

	source := s.options.Nodes[0]
	if _, err := url.Parse(source); err != nil {
		return err
	}

	client, err := sql.Open("postgres", source)
	if err != nil {
		return err
	}

	if err := client.Ping(); err != nil {
		return err
	}

	s.client = client

	return s.initDB()
}

func (s *cockroachStore) initDB() error {
	if _, err := s.client.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", s.options.Database)); err != nil {
		return err
	}

	if _, err := s.client.Exec(fmt.Sprintf("SET DATABASE = %s ;", s.options.Database)); err != nil {
		return err
	}

	if _, err := s.client.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
	(
		key text NOT NULL,
		value bytea,
		expiry timestamp with time zone,
		CONSTRAINT %s_pkey PRIMARY KEY (key)
	);`, s.options.Table, s.options.Table)); err != nil {
		return err
	}

	if _, err := s.client.Exec(fmt.Sprintf(`CREATE INDEX IF NOT EXISTS "%s" ON %s.%s USING btree ("key");`, "key_index_"+s.options.Table, s.options.Database, s.options.Table)); err != nil {
		return err
	}

	write, err := s.client.Prepare(fmt.Sprintf(`INSERT INTO %s.%s(key, value, expiry)
		VALUES ($1, $2::bytea, $3)
		ON CONFLICT (key)
		DO UPDATE
		SET value = EXCLUDED.value, expiry = EXCLUDED.expiry;`, s.options.Database, s.options.Table))
	if err != nil {
		return err
	}
	s.write = write

	readOne, err := s.client.Prepare(fmt.Sprintf("SELECT key, value, expiry FROM %s.%s WHERE key = $1;", s.options.Database, s.options.Table))
	if err != nil {
		return err
	}
	s.readOne = readOne

	readMany, err := s.client.Prepare(fmt.Sprintf("SELECT key, value, expiry FROM %s.%s WHERE key LIKE $1;", s.options.Database, s.options.Table))
	if err != nil {
		return err
	}
	s.readMany = readMany

	readOffset, err := s.client.Prepare(fmt.Sprintf("SELECT key, value, expiry FROM %s.%s WHERE key LIKE $1 ORDER BY key DESC LIMIT $2 OFFSET $3;", s.options.Database, s.options.Table))
	if err != nil {
		return err
	}
	s.readOffset = readOffset

	list, err := s.client.Prepare(fmt.Sprintf("SELECT key, value, expiry FROM %s.%s;", s.options.Database, s.options.Table))
	if err != nil {
		return err
	}
	s.list = list

	delete, err := s.client.Prepare(fmt.Sprintf("DELETE FROM %s.%s WHERE key = $1;", s.options.Database, s.options.Table))
	if err != nil {
		return err
	}
	s.delete = delete

	return nil
}

func NewStore(opts ...store.StoreOption) store.Store {
	options := store.NewStoreOptions(opts...)

	s := &cockroachStore{
		options: options,
	}

	if err := s.configure(); err != nil {
		log.Fatal(err)
	}

	return s
}
