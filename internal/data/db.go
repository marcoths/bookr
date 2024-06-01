package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/marcoths/bookr/internal/bucket"
	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

var ErrNotFound = errors.New("not found")

// Get returns a seat from the database
func Get(db *bolt.DB, name string) (*Seat, error) {

	var seat *Seat
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.Seats.Name())

		record := b.Get([]byte(name))
		if record == nil {
			return fmt.Errorf("%s: %w", name, ErrNotFound)
		}

		if err := json.Unmarshal(record, &seat); err != nil {
			return fmt.Errorf("unmarshal seat: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return seat, nil
}

// FindConsecutive finds the first n consecutive seats starting from seatNumber with the provided status
func FindConsecutive(db *bolt.DB, row string, seatNumber, n uint, status SeatStatus) ([]*Seat, error) {
	start := fmt.Sprintf("%s%d", row, seatNumber)
	end := fmt.Sprintf("%s%d", row, (seatNumber+n)-1)
	found, err := Range(db, start, end, status)
	if err != nil {
		return nil, err
	}
	if len(found) >= int(n) {
		return found[:n], nil
	}
	return nil, fmt.Errorf("couldn't find %d %s seats", n, status.String())
}

func Range(db *bolt.DB, start, end string, status SeatStatus) ([]*Seat, error) {
	var seats []*Seat
	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucket.Seats.Name()).Cursor()
		for k, v := c.Seek([]byte(start)); k != nil && bytes.Compare(k, []byte(end)) <= 0; k, v = c.Next() {
			seat := &Seat{}
			if err := json.Unmarshal(v, seat); err != nil {
				return fmt.Errorf("unmarshal seat: %w", err)
			}
			if seat.Status == status {
				seats = append(seats, seat)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return seats, nil
}

// Create adds a new seat to the database
func Create(db *bolt.DB, seat *Seat) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.Seats.Name())

		record, err := json.Marshal(seat)
		if err != nil {
			return fmt.Errorf("marshal seat: %w", err)
		}

		if err := b.Put([]byte(seat.ID), record); err != nil {
			return fmt.Errorf("put seat: %w", err)
		}

		return nil
	})
}

// UpdateStatusTo updates the status of a seat
func UpdateStatusTo(db *bolt.DB, id string, status SeatStatus) error {

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.Seats.Name())

		record := b.Get([]byte(id))
		if record == nil {
			return fmt.Errorf("%s: %w", id, ErrNotFound)
		}

		seat := &Seat{}
		if err := json.Unmarshal(record, seat); err != nil {
			return fmt.Errorf("unmarshal seat: %w", err)
		}

		seat.Status = status

		record, err := json.Marshal(seat)
		if err != nil {
			return fmt.Errorf("marshal seat: %w", err)
		}

		if err := b.Put([]byte(id), record); err != nil {
			return fmt.Errorf("put seat: %w", err)
		}

		return nil
	})
}

// SetContext creates a new database and registers the necessary buckets for testing
func SetContext(t testing.TB, bucketName []byte) *bolt.DB {
	dbFile, err := os.CreateTemp("", "*")
	assert.NoError(t, err)

	db, err := bolt.Open(dbFile.Name(), 0o600, &bolt.Options{Timeout: 1 * time.Second})
	assert.NoError(t, err, "Failed connecting to the database")

	err = db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(bucketName)
		if _, err := tx.CreateBucketIfNotExists(bucketName); err != nil {
			return fmt.Errorf("couldn't create %q bucket: %w", bucketName, err)
		}
		return nil
	})
	assert.NoError(t, err)

	t.Cleanup(func() {
		err := db.Close()
		assert.NoError(t, err, "Failed closing the database")
	})

	return db
}

// Seed adds the provided data to the database
func Seed(db *bolt.DB, data []byte) error {

	// Try to get one seat, if it exists, assume the data is already seeded and return
	if _, err := Get(db, "A1"); err == nil {
		return nil
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.Seats.Name())

		var seats []*Seat
		if err := json.Unmarshal(data, &seats); err != nil {
			return fmt.Errorf("unmarshal seats: %w", err)
		}
		for _, seat := range seats {
			record, err := json.Marshal(seat)
			if err != nil {
				return fmt.Errorf("marshal seat: %w", err)
			}
			if err := b.Put([]byte(seat.ID), record); err != nil {
				return fmt.Errorf("put seat: %w", err)
			}
		}
		return nil
	})
}

// Register creates the necessary buckets for the application
func Register(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		// Create all the buckets except auth, it will be created in setParameters()
		name := bucket.Seats.Name()
		if _, err := tx.CreateBucketIfNotExists(name); err != nil {
			return fmt.Errorf("create bucket %s: %w", name, err)
		}
		return nil
	})
}
