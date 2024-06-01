package cmd

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/marcoths/bookr/internal/bucket"
	"github.com/marcoths/bookr/internal/data"
	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestCommands(t *testing.T) {
	db := SetContext(t)

	createBooking(t, db)

	cases := []struct {
		name       string
		args       []string
		assertFunc func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool
	}{
		{
			name:       "Book",
			args:       []string{"book", "A0", "1"},
			assertFunc: assert.NoError,
		},
		{
			name:       "Cancel A0",
			args:       []string{"cancel", "A0", "1"},
			assertFunc: assert.NoError,
		},
		{
			name:       "Book A0 again",
			args:       []string{"book", "A0", "1"},
			assertFunc: assert.NoError,
		},
		{
			name:       "Book A0 expect error",
			args:       []string{"book", "A0", "1"},
			assertFunc: assert.Error,
		},
		{
			name:       "Book A1",
			args:       []string{"book", "A1", "1"},
			assertFunc: assert.NoError,
		},
		{
			name:       "Book A2 4 seats",
			args:       []string{"book", "A2", "4"},
			assertFunc: assert.NoError,
		},
		{
			name:       "Book A5 1 seat expect error",
			args:       []string{"book", "A5", "1"},
			assertFunc: assert.Error,
		},
		{
			name:       "Book A6 3 seats expect error",
			args:       []string{"book", "A6", "3"},
			assertFunc: assert.Error,
		},
		{
			name:       "Book A8 1 seats expect error",
			args:       []string{"book", "A8", "1"},
			assertFunc: assert.Error,
		},
		{
			name:       "Book U1 1 seats expect error",
			args:       []string{"book", "U1", "1"},
			assertFunc: assert.Error,
		},
	}

	cmd := NewRootCmd(db)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd.SetArgs(tc.args)
			err := cmd.Execute()
			tc.assertFunc(t, err)
		})

	}
}

// SetContext sets up the testing environment.
//
// It uses t.Cleanup() to close the database connection after the test and
// all its subtests are completed.
func SetContext(t testing.TB) *bolt.DB {
	t.Helper()

	dbFile, err := os.CreateTemp("", "*")
	assert.NoError(t, err)

	db, err := bolt.Open(dbFile.Name(), 0o600, &bolt.Options{Timeout: 1 * time.Second})
	assert.NoError(t, err, "Failed connecting to the database")
	bucketName := bucket.Seats.Name()
	db.Update(func(tx *bolt.Tx) error {
		// Ignore errors on purpose
		tx.DeleteBucket(bucketName)
		tx.CreateBucketIfNotExists(bucketName)

		return nil
	})
	os.Stdout = os.NewFile(0, "") // Mute stdout
	os.Stderr = os.NewFile(0, "") // Mute stderr
	t.Cleanup(func() {
		assert.NoError(t, db.Close(), "Failed connecting to the database")
	})

	return db
}

func createBooking(t *testing.T, db *bolt.DB) {
	// Create a booking
	t.Helper()

	for i := 0; i <= 7; i++ {
		s := &data.Seat{
			ID:     "A" + strconv.Itoa(i),
			Row:    "A",
			Number: uint(i),
			Status: data.Available,
		}
		_ = data.Create(db, s)
	}

}
