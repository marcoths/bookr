package data_test

import (
	"testing"

	"github.com/marcoths/bookr/internal/bucket"
	"github.com/marcoths/bookr/internal/data"
	"github.com/stretchr/testify/assert"
)

var (
	seat       = &data.Seat{ID: "A1", Row: "A", Number: 1, Status: data.Available}
	bucketName = bucket.Seats.Name()
)

func TestDB(t *testing.T) {
	db := data.SetContext(t, bucketName)

	err := data.Create(db, seat)
	assert.NoError(t, err)

	got, err := data.Get(db, seat.ID)
	assert.NoError(t, err)

	assert.Equal(t, seat, got)
}

func TestFindConsecutive(t *testing.T) {
	db := data.SetContext(t, bucketName)
	seats := []*data.Seat{
		{ID: "A0", Row: "A", Number: 0, Status: data.Booked},
		{ID: "A1", Row: "A", Number: 1, Status: data.Booked},
		{ID: "A2", Row: "A", Number: 2, Status: data.Booked},
		{ID: "A3", Row: "A", Number: 3, Status: data.Booked},
		{ID: "A4", Row: "A", Number: 4, Status: data.Booked},
		{ID: "A5", Row: "A", Number: 5, Status: data.Booked},
		{ID: "A6", Row: "A", Number: 6, Status: data.Available},
		{ID: "A7", Row: "A", Number: 7, Status: data.Available},
		{ID: "B0", Row: "B", Number: 0, Status: data.Available},
		{ID: "B1", Row: "B", Number: 1, Status: data.Available},
		{ID: "B2", Row: "B", Number: 2, Status: data.Available},
		{ID: "B3", Row: "B", Number: 3, Status: data.Available},
		{ID: "B3", Row: "B", Number: 4, Status: data.Available},
	}
	for _, s := range seats {
		_ = data.Create(db, s)
	}

	got, err := data.FindConsecutive(db, "A", 5, 3, data.Available)

	// it should return an error because there are no consecutive seats available
	assert.Error(t, err)

	// it should return 3 consecutive seats starting from B0
	want := []*data.Seat{
		{ID: "B0", Row: "B", Number: 0, Status: data.Available},
		{ID: "B1", Row: "B", Number: 1, Status: data.Available},
		{ID: "B2", Row: "B", Number: 2, Status: data.Available},
	}

	got, err = data.FindConsecutive(db, "B", 0, 3, data.Available)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
